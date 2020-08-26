package phash

import (
	"image"
	"math/bits"

	"github.com/corona10/goimagehash"
)

// PHash is the type representing the signed perception hash of an image
type PHash uint64

// GetPHash returns the signed PHash of a given image
func GetPHash(img image.Image) PHash {
	pHash, _ := goimagehash.PerceptionHash(img)
	return PHash(pHash.GetHash())
}

func hammingDist(hash1, hash2 PHash) int {
	xorResult := uint64(hash1) ^ uint64(hash2)
	return bits.OnesCount64(xorResult)
}

// NormHammingDist returns the normalized hamming distance between two PHashes
func NormHammingDist(hash1, hash2 PHash) float64 {
	return float64(hammingDist(hash1, hash2)) / 64
}
