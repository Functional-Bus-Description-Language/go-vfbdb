package utils

import (
	"fmt"
)

// Uint64Mask returns 64 bits mask.
// If mask can't be represented on 64 bits it panics.
// The returned mask is always shifted to the right. For example, the result for
// Mask(2, 1) is 3 (0b11), not 6 (0b110).
func Uint64Mask(startBit, endBit int64) uint64 {
	width := endBit - startBit + 1
	if width > 64 {
		panic(fmt.Sprintf("cannot convert mask of width %d to uint64", width))
	}
	return (1 << (endBit - startBit + 1)) - 1
}
