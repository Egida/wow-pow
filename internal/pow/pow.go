package pow

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"
	"time"
)

const nonceSizeInBytes = 8

type FindNonceResult struct {
	Nonce        uint64
	Hash         []byte
	Duration     time.Duration
	LeadingZeros int
	Error        error
}

type task struct {
	LowerBound uint64
	UpperBound uint64
}

func FindNonce(difficulty uint8, token []byte, concurrency int) FindNonceResult {
	if int(difficulty) > 8*len(token) {
		return FindNonceResult{
			Error: errors.New("difficulty more then bits in the given token"),
		}
	}

	var (
		begin       = time.Now()
		taskCh      = make(chan task)
		resultCh    = make(chan FindNonceResult)
		doneCh      = make(chan struct{})
		ctx, cancel = context.WithCancel(context.Background())
	)

	for i := 0; i < concurrency; i++ {
		t := make([]byte, len(token))
		copy(t, token)
		go findNonceWorker(i, ctx, taskCh, resultCh, doneCh, difficulty, t)
	}

	go func() {
		defer close(taskCh)

		step := uint64(math.MaxUint64 / uint64(concurrency))
		i, next := uint64(0), uint64(0)

		for next != math.MaxUint64 {
			i = next
			next = calcUpperBoundUint64(i, step)

			select {
			case <-ctx.Done():
				return
			case taskCh <- task{
				LowerBound: i,
				UpperBound: next,
			}:
			}
		}
	}()

	cDone := 0

	for {
		select {
		case result := <-resultCh:
			result.Duration = time.Since(begin)
			cancel()

			return result
		case <-doneCh:
			cDone++
			if cDone == concurrency {
				return FindNonceResult{
					Error: errors.New("failure to find nonce for the given difficulty"),
				}
			}
		}
	}
}

func findNonceWorker(
	i int,
	ctx context.Context, taskCh chan task, resultCh chan FindNonceResult,
	doneCh chan struct{}, difficulty uint8, token []byte,
) {
	defer func() {
		select {
		case <-ctx.Done():
		case doneCh <- struct{}{}:
		}
	}()

	var (
		hash          [32]byte
		cLeadingZeros int
		difficultyInt int = int(difficulty)
		buf               = append(token, make([]byte, nonceSizeInBytes)...)
	)

	for task := range taskCh {
		for nonce := task.LowerBound; nonce < task.UpperBound; nonce++ {
			select {
			case <-ctx.Done():
				return
			default:
			}

			binary.BigEndian.PutUint64(buf[len(token):], nonce)

			hash = sha256.Sum256(buf)
			cLeadingZeros = countLeadingZeros(hash[:])

			if cLeadingZeros >= difficultyInt {
				select {
				case <-ctx.Done():
				case resultCh <- FindNonceResult{
					Nonce:        nonce,
					Hash:         hash[:],
					LeadingZeros: cLeadingZeros,
				}:
				}

				return
			}
		}
	}

}

func CheckSolution(token []byte, nonce uint64, difficulty uint8) bool {
	buf := append(token, make([]byte, nonceSizeInBytes)...)

	binary.BigEndian.PutUint64(buf[len(token):], nonce)

	hash := sha256.Sum256(buf)

	if countLeadingZeros(hash[:]) >= int(difficulty) {
		return true
	}

	return false
}

func GenerateToken(buf []byte) error {
	if _, err := rand.Read(buf); err != nil {
		return err
	}

	return nil
}
