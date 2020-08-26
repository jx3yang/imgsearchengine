package phash

import (
	"strconv"
	"testing"

	_ "image/jpeg"
	_ "image/png"
)

func TestHammingDistance(t *testing.T) {
	// arrange
	hashstr1 := "1001011010110"
	hashstr2 := "1100100111000"

	// number of bits that are different
	want := 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1

	hash1, _ := strconv.ParseUint(hashstr1, 2, 64)
	hash2, _ := strconv.ParseUint(hashstr2, 2, 64)

	// act
	got := hammingDist(PHash(hash1), PHash(hash2))

	// assert
	if got != want {
		t.Errorf("findMedian() = %d, want %d", got, want)
	}
}
