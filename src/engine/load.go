package engine

import (
	"encoding/csv"
	"errors"
	"image"
	"io"
	"log"
	"os"
	"strconv"

	// for decoding
	_ "image/jpeg"
	_ "image/png"

	phash "github.com/jx3yang/imgsearchengine/src/phash"
	vptree "github.com/jx3yang/imgsearchengine/src/vptree"
)

const (
	pathCol  string = "path"
	phashCol string = "phash"
)

func distanceFnc(img1, img2 interface{}) float64 {
	return phash.NormHammingDist(img1.(*ImageInfo).GetPHash(), img2.(*ImageInfo).GetPHash())
}

func processCSV(rc io.Reader, sep rune) (<-chan []string, <-chan []string) {
	ch := make(chan []string)
	headch := make(chan []string)
	go func() {
		r := csv.NewReader(rc)
		r.Comma = sep
		defer close(ch)
		defer close(headch)

		headers, err := r.Read()
		if err != nil {
			log.Fatal(err)
		}
		headch <- headers

		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			ch <- rec
		}

	}()
	return headch, ch
}

func processEntries(ch <-chan []string, pathIdx int, phashIdx int) <-chan *ImageInfo {
	imgCh := make(chan *ImageInfo)
	computePHash := phashIdx < 0

	go func() {
		defer close(imgCh)
		for elem := range ch {
			path := elem[pathIdx]
			var hash phash.PHash
			if computePHash {
				file, err := os.Open(path)
				if err != nil {
					log.Fatal("Unabled to read ", path)
				}
				img, _, _ := image.Decode(file)
				hash = phash.GetPHash(img)
				file.Close()
			} else {
				n, err := strconv.ParseUint(elem[phashIdx], 10, 64)
				if err != nil {
					log.Fatal("Image with path ", path, " has invalid PHash")
				}
				hash = phash.PHash(n)
			}

			imgCh <- NewImageInfo(hash, path)
		}
	}()

	return imgCh
}

func parseColumns(csvFile *os.File, sep rune, withPhashCol bool) (<-chan *ImageInfo, error) {
	headch, ch := processCSV(csvFile, sep)
	headers := <-headch

	// find the column containing the paths and phashes
	pathIdx := 0
	phashIdx := -1
	if withPhashCol {
		phashIdx = 0
	}

	foundPathColumn := false
	foundPhashColumn := !withPhashCol

	for i, col := range headers {
		if !foundPathColumn && col == pathCol {
			pathIdx = i
			foundPathColumn = true
		}
		if !foundPhashColumn && col == phashCol {
			phashIdx = i
			foundPhashColumn = true
		}
		if foundPathColumn && foundPhashColumn {
			break
		}
	}

	if !foundPathColumn {
		return nil, errors.New("Did not find the path column")
	}

	if !foundPhashColumn {
		return nil, errors.New("Did not find the phash column")
	}

	return processEntries(ch, pathIdx, phashIdx), nil
}

func load(csvPath string, sep rune, withPhashCol bool) (*vptree.VPTree, error) {
	csvFile, err := os.Open(csvPath)
	defer csvFile.Close()
	if err != nil {
		return nil, err
	}

	ch, err := parseColumns(csvFile, sep, withPhashCol)
	if err != nil {
		return nil, err
	}

	points := make([]interface{}, 0)

	for elem := range ch {
		points = append(points, elem)
	}

	tree := vptree.BuildTree(points, distanceFnc)
	return tree, nil
}

// LoadFromCSV loads the given CSV file containing the paths
// of the images, computes the PHashes of the images, and returns the VP-Tree
// containing the PHash and path of each image
func LoadFromCSV(csvPath string, sep rune) (*vptree.VPTree, error) {
	return load(csvPath, sep, false)
}

// LoadFromCSVPHash loads the given CSV file containg the paths
// of the images and the corresponding PHashes, and returns the
// VP-Tree containing the PHash and path of each image
// Must contain the headers "phash" and "path"
func LoadFromCSVPHash(csvPath string, sep rune) (*vptree.VPTree, error) {
	return load(csvPath, sep, true)
}
