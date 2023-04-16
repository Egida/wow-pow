package pow

import (
	"math"
	"math/bits"
)

func calcUpperBoundUint64(i, step uint64) uint64 {
	if math.MaxUint64-step <= i {
		return math.MaxUint64
	}

	return i + step
}

func countLeadingZeros(data []byte) int {
	result := 0

	for _, v := range data {
		if v == 0 {
			result += 8
		} else {
			result += bits.LeadingZeros8(v)
			break
		}
	}

	return result

}
