package vptree

import (
	"errors"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/ef-ds/deque"
)

type kvp struct {
	key   float64
	value interface{}
}

type heap struct {
	Data []kvp
}

func (h *heap) Root() float64 { return h.Data[0].key }

func (h *heap) Len() int { return len(h.Data) }

func (h *heap) leftChildIdx(k int) int { return 2*k + 1 }

func (h *heap) rightChildIdx(k int) int { return 2*k + 2 }

func (h *heap) parentIdx(k int) int { return (k - 1) / 2 }

func (h *heap) fixUp(k int) {
	for k > 0 {
		parentIdx := h.parentIdx(k)
		parentItem := h.Data[parentIdx]
		currentItem := h.Data[k]
		if parentItem.key < currentItem.key {
			h.swap(parentIdx, k)
			k = parentIdx
		} else {
			return
		}
	}
}

func (h *heap) fixDown(k int) {
	n := h.Len()
	for k < n {
		leftIdx := h.leftChildIdx(k)
		rightIdx := h.rightChildIdx(k)
		if leftIdx < n {
			maxChildIdx := leftIdx
			if rightIdx < n && h.Data[leftIdx].key < h.Data[rightIdx].key {
				maxChildIdx = rightIdx
			}

			if h.Data[k].key < h.Data[maxChildIdx].key {
				h.swap(k, maxChildIdx)
				k = maxChildIdx
			} else {
				return
			}
		} else {
			return
		}
	}
}

func (h *heap) swap(i, j int) {
	h.Data[i], h.Data[j] = h.Data[j], h.Data[i]
}

func (h *heap) Pop() (*kvp, error) {
	if h.Len() == 0 {
		return nil, errors.New("Heap is empty")
	}
	h.swap(0, h.Len()-1)
	target := h.Data[h.Len()-1]
	h.Data = h.Data[:h.Len()-1]
	h.fixDown(0)
	return &target, nil
}

func (h *heap) Push(item kvp) {
	h.Data = append(h.Data, item)
	h.fixUp(h.Len() - 1)
}

// VPNode is a node in a VPTree
type VPNode struct {
	Left         *VPNode
	Right        *VPNode
	LeftMin      float64
	LeftMax      float64
	RightMin     float64
	RightMax     float64
	VantagePoint interface{}
}

func makeNode(point interface{}) *VPNode {
	n := new(VPNode)
	n.VantagePoint = point
	n.Left = nil
	n.Right = nil
	n.LeftMin = math.MaxFloat64
	n.LeftMax = 0
	n.RightMin = math.MaxFloat64
	n.RightMax = 0
	return n
}

// DistanceFnc is a function that computes the distance between two points
// Requires: DistanceFnc(point1, point2) >= 0
type DistanceFnc func(point1, point2 interface{}) float64

// VPTree implements the Vantage Point Tree
type VPTree struct {
	Root        *VPNode
	distanceFnc DistanceFnc
}

func kthElement(distances []float64, k int) float64 {
	// if the slice is small, simply sort it
	n := len(distances)
	if n <= 30 {
		sort.Float64s(distances)
		return distances[k]
	}

	if n == 1 && k == 0 {
		return distances[k]
	}

	// use the randomized quick select algorithm
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(n)
	distances[idx], distances[n-1] = distances[n-1], distances[idx]

	i := 0
	j := n - 2
	pivotIdx := n - 1
	pivot := distances[pivotIdx]

	for true {
		for i < n-1 && distances[i] <= pivot {
			i++
		}
		for j > 0 && distances[j] > pivot {
			j--
		}
		if i < j {
			distances[i], distances[j] = distances[j], distances[i]
		} else {
			distances[i], distances[pivotIdx] = distances[pivotIdx], distances[i]
			pivotIdx = i
			break
		}
	}

	if pivotIdx == k {
		return distances[pivotIdx]
	} else if pivotIdx > k {
		return kthElement(distances[:pivotIdx], k)
	} else {
		return kthElement(distances[pivotIdx+1:], k-pivotIdx-1)
	}
}

func findMedian(distances []float64) float64 {
	copySlice := make([]float64, len(distances))
	copy(copySlice, distances)
	return kthElement(copySlice, len(distances)/2)
}

func buildTree(points []interface{}, distanceFnc DistanceFnc) *VPNode {
	if len(points) == 0 {
		return nil
	}
	currentNode := makeNode(points[len(points)-1])
	points = points[:len(points)-1]

	if len(points) == 0 {
		return currentNode
	}

	distances := make([]float64, len(points))

	for i, point := range points {
		distances[i] = distanceFnc(currentNode.VantagePoint, point)
	}

	median := findMedian(distances)

	var leftPoints []interface{}
	var rightPoints []interface{}

	swapLastTwo := func(slice []interface{}) {
		lenSlice := len(slice)
		slice[lenSlice-1], slice[lenSlice-2] = slice[lenSlice-2], slice[lenSlice-1]
	}

	for i, dist := range distances {
		if dist >= median {
			currentNode.RightMin = math.Min(dist, currentNode.RightMin)
			currentNode.RightMax = math.Max(dist, currentNode.RightMax)
			rightPoints = append(rightPoints, points[i])

			if currentNode.RightMax != dist {
				// make sure that the last element of the slice is
				// the furthest point from the current vp point
				swapLastTwo(rightPoints)
			}
		} else {
			currentNode.LeftMin = math.Min(dist, currentNode.LeftMin)
			currentNode.LeftMax = math.Max(dist, currentNode.LeftMax)
			leftPoints = append(leftPoints, points[i])

			if currentNode.LeftMax != dist {
				swapLastTwo(leftPoints)
			}
		}
	}

	currentNode.Left = buildTree(leftPoints, distanceFnc)
	currentNode.Right = buildTree(rightPoints, distanceFnc)
	return currentNode
}

// BuildTree will return the root of the VP-Tree built from
// `points` using the `distanceFnc` for computing the distances
func BuildTree(points []interface{}, distanceFnc DistanceFnc) *VPTree {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(points), func(i, j int) { points[i], points[j] = points[j], points[i] })
	root := buildTree(points, distanceFnc)
	return &VPTree{
		Root:        root,
		distanceFnc: distanceFnc,
	}
}

// KNNSearch will return the k nearest neighbours of the given `point`
// in the VP-Tree
func (tree *VPTree) KNNSearch(point interface{}, k uint) (map[interface{}]float64, error) {
	if k < 1 {
		return nil, errors.New("Invalid k")
	}
	root := tree.Root
	nodesToVisit := deque.New()
	nodesToVisit.PushFront(kvp{0, root})

	tau := math.MaxFloat64

	results := new(heap)

	resultsLen := func() uint { return uint(results.Len()) }

	for nodesToVisit.Len() > 0 {
		pair, _ := nodesToVisit.PopFront()
		kvpObj := pair.(kvp)
		d0, currentNode := kvpObj.key, kvpObj.value.(*VPNode)
		if currentNode == nil || d0 > tau {
			continue
		}
		dist := tree.distanceFnc(point, currentNode.VantagePoint)

		if dist < tau {
			if resultsLen() == k {
				results.Pop()
			}
			results.Push(kvp{dist, currentNode.VantagePoint})
			if resultsLen() == k {
				tau = results.Root()
			}
		}

		if currentNode.Left == nil && currentNode.Right == nil {
			continue
		}

		if currentNode.LeftMin <= dist && dist <= currentNode.LeftMax {
			nodesToVisit.PushFront(kvp{0, currentNode.Left})
		} else if currentNode.LeftMin-tau <= dist && dist <= currentNode.LeftMax+tau {
			if dist < currentNode.LeftMin {
				nodesToVisit.PushBack(kvp{currentNode.LeftMin - dist, currentNode.Left})
			} else {
				nodesToVisit.PushBack(kvp{dist - currentNode.LeftMax, currentNode.Left})
			}
		}

		if currentNode.RightMin <= dist && dist <= currentNode.RightMax {
			nodesToVisit.PushFront(kvp{0, currentNode.Right})
		} else if currentNode.RightMin-tau <= dist && dist <= currentNode.RightMax+tau {
			if dist < currentNode.RightMin {
				nodesToVisit.PushBack(kvp{currentNode.RightMin - dist, currentNode.Right})
			} else {
				nodesToVisit.PushBack(kvp{dist - currentNode.RightMax, currentNode.Right})
			}
		}
	}

	knnMap := make(map[interface{}]float64)

	for _, pair := range results.Data {
		knnMap[pair.value] = pair.key
	}

	return knnMap, nil
}

// RangeSearch will return all the points within a `threshold`
// distance from the given `point`
func (tree *VPTree) RangeSearch(point interface{}, threshold float64) (map[interface{}]float64, error) {
	if threshold < 0 {
		return nil, errors.New("Threshold must be positive")
	}

	root := tree.Root
	nodesToVisit := deque.New()
	nodesToVisit.PushFront(kvp{0, root})

	rangeMap := make(map[interface{}]float64)

	for nodesToVisit.Len() > 0 {
		pair, _ := nodesToVisit.PopFront()
		kvpObj := pair.(kvp)
		d0, currentNode := kvpObj.key, kvpObj.value.(*VPNode)
		if currentNode == nil || d0 > threshold {
			continue
		}

		dist := tree.distanceFnc(point, currentNode.VantagePoint)
		if dist <= threshold {
			rangeMap[currentNode.VantagePoint] = dist
		}

		if currentNode.Left == nil && currentNode.Right == nil {
			continue
		}

		if currentNode.LeftMin <= dist && dist <= currentNode.LeftMax {
			nodesToVisit.PushFront(kvp{0, currentNode.Left})
		} else if currentNode.LeftMin-threshold <= dist && dist <= currentNode.LeftMax+threshold {
			if dist < currentNode.LeftMin {
				nodesToVisit.PushBack(kvp{currentNode.LeftMin - dist, currentNode.Left})
			} else {
				nodesToVisit.PushBack(kvp{dist - currentNode.LeftMax, currentNode.Left})
			}
		}

		if currentNode.RightMin <= dist && dist <= currentNode.RightMax {
			nodesToVisit.PushFront(kvp{0, currentNode.Right})
		} else if currentNode.RightMin-threshold <= dist && dist <= currentNode.RightMax+threshold {
			if dist < currentNode.RightMin {
				nodesToVisit.PushBack(kvp{currentNode.RightMin - dist, currentNode.Right})
			} else {
				nodesToVisit.PushBack(kvp{dist - currentNode.RightMax, currentNode.Right})
			}
		}

	}

	return rangeMap, nil
}
