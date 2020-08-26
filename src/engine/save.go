package engine

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	vptree "github.com/jx3yang/imgsearchengine/src/vptree"
)

func treeTraversal(tree *vptree.VPTree) <-chan *ImageInfo {
	root := tree.Root
	ch := make(chan *ImageInfo)

	var walk func(*vptree.VPNode)
	walk = func(node *vptree.VPNode) {
		ch <- node.VantagePoint.(*ImageInfo)
		if node.Left != nil {
			walk(node.Left)
		}
		if node.Right != nil {
			walk(node.Right)
		}
	}

	go func() {
		defer close(ch)
		walk(root)
	}()

	return ch
}

// SaveTreeInfo will save the content of a VP-Tree holding
// *ImageInfo structs as nodes in a CSV file given its path
func SaveTreeInfo(tree *vptree.VPTree, csvPath string, sep rune) {
	file, err := os.Create(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = sep
	defer writer.Flush()

	writer.Write([]string{pathCol, phashCol})

	ch := treeTraversal(tree)

	for elem := range ch {
		row := []string{elem.GetPath(), strconv.FormatUint(uint64(elem.GetPHash()), 10)}
		writer.Write(row)
	}
}
