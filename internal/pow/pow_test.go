package pow

import (
	"crypto/rand"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCheckSolution(t *testing.T) {
	tt := []struct {
		token      []byte
		nonce      uint64
		difficulty uint8
		result     bool
	}{
		{
			token:      []byte{85, 143, 166, 105, 104, 131, 154, 125, 140, 160},
			nonce:      9223372036859668424,
			difficulty: 24,
			result:     true,
		},
		{
			token:      []byte{85, 143, 166, 105, 104, 131, 154, 125, 140, 160},
			nonce:      8223372036859668424,
			difficulty: 24,
			result:     false,
		},
	}

	for i, tc := range tt {
		result := CheckSolution(tc.token, tc.nonce, tc.difficulty)
		require.Equal(t, tc.result, result, fmt.Sprintf("case #%d", i+1))
	}
}

func TestFindNonce(t *testing.T) {
	tt := []struct {
		token       []byte
		nonce       uint64
		difficulty  uint8
		concurrency int
		result      FindNonceResult
	}{
		{
			token:       []byte{85, 143, 166, 105, 104, 131, 154, 125, 140, 160},
			difficulty:  18,
			concurrency: 1,
			result: FindNonceResult{
				Hash:         []byte{0, 0, 21, 116, 138, 195, 212, 228, 58, 159, 63, 135, 106, 158, 162, 84, 216, 30, 105, 77, 58, 33, 90, 190, 113, 254, 30, 131, 92, 27, 82, 11},
				Nonce:        53302,
				LeadingZeros: 19,
			},
		},
		{
			token:       []byte{85, 143, 166, 105, 104, 131, 154, 125, 140, 160},
			difficulty:  18,
			concurrency: 2,
			result: FindNonceResult{
				Hash:         []byte{0, 0, 21, 116, 138, 195, 212, 228, 58, 159, 63, 135, 106, 158, 162, 84, 216, 30, 105, 77, 58, 33, 90, 190, 113, 254, 30, 131, 92, 27, 82, 11},
				Nonce:        53302,
				LeadingZeros: 19,
			},
		},
		{
			token:       []byte{85, 143, 166, 105, 104, 131, 154, 125, 140, 160},
			difficulty:  18,
			concurrency: 4,
			result: FindNonceResult{
				Hash:         []byte{0, 0, 28, 126, 191, 220, 213, 31, 228, 173, 172, 3, 146, 213, 165, 111, 9, 10, 177, 61, 141, 82, 210, 74, 178, 39, 44, 128, 48, 168, 4, 159},
				Nonce:        13835058055282197610,
				LeadingZeros: 19,
			},
		},
		{
			token:       []byte{85, 143, 166, 105, 104, 131, 154, 125, 140, 160},
			difficulty:  18,
			concurrency: 8,
			result: FindNonceResult{
				Hash:         []byte{0, 0, 21, 116, 138, 195, 212, 228, 58, 159, 63, 135, 106, 158, 162, 84, 216, 30, 105, 77, 58, 33, 90, 190, 113, 254, 30, 131, 92, 27, 82, 11},
				Nonce:        53302,
				LeadingZeros: 19,
			},
		},
	}

	for i, tc := range tt {
		require.Eventually(t, func() bool {
			result := FindNonce(tc.difficulty, tc.token, tc.concurrency)
			return reflect.DeepEqual(tc.result.Hash, result.Hash) &&
				tc.result.Nonce == result.Nonce &&
				tc.result.LeadingZeros == result.LeadingZeros
		},
			4*time.Second,
			100*time.Millisecond,
			fmt.Sprintf("case #%d", i+1),
		)
	}
}

func BenchmarkFindNonceConcurrency(b *testing.B) {
	// token := []byte{85, 143, 166, 105, 104, 131, 154, 125, 140, 160}
	tokenSize := 64
	difficulty := uint8(20)
	concurrency := []int{1, 2, 4, 8, 16}

	for _, c := range concurrency {
		b.Run(fmt.Sprintf("concurrency_%d", c), func(b *testing.B) {
			token := make([]byte, tokenSize)
			rand.Read(token)

			for i := 0; i < b.N; i++ {
				FindNonce(difficulty, token, c)
			}
		})
	}
}
