package engine

import phash "github.com/jx3yang/imgsearchengine/src/phash"

// ImageInfo contains the PHash as well as the path of an image
type ImageInfo struct {
	hash phash.PHash
	path string
}

// NewImageInfo returns a struct containing the hash and path of the image
func NewImageInfo(hash phash.PHash, path string) *ImageInfo {
	return &ImageInfo{
		hash: hash,
		path: path,
	}
}

// GetPHash returns the PHash of the associated image
func (imgInfo *ImageInfo) GetPHash() phash.PHash { return imgInfo.hash }

// GetPath returns the path of the associated image
func (imgInfo *ImageInfo) GetPath() string { return imgInfo.path }
