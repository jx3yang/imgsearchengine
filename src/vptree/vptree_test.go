package vptree

import (
	"math"
	"reflect"
	"testing"
)

func TestMedian(t *testing.T) {
	// arrange
	s := []float64{
		89, 86, 91, 8, 46, 3, 23, 52, 24, 43, 75, 10, 81, 54,
		64, 50, 15, 92, 6, 63, 38, 16, 17, 45, 85, 2, 90, 59,
		60, 73, 22, 70, 12, 31, 14, 57, 79, 30, 78, 61, 29,
		65, 62, 84, 39, 83, 42, 32, 56, 20, 97, 40, 77, 67,
		87, 11, 27, 93, 47, 88, 9, 28, 96, 55, 19, 25, 26, 69,
		72, 71, 1, 51, 100, 5, 4, 7, 99, 35, 82, 44, 74, 36,
		66, 34, 33, 58, 68, 80, 76, 53, 18, 49, 98, 37, 21, 95, 48, 41, 94, 13,
	}

	want := float64(51)

	// act
	got := findMedian(s)

	// assert
	if got != want {
		t.Errorf("findMedian() = %f, want %f", got, want)
	}
}

func TestKNNSearch(t *testing.T) {
	// arrange
	points := make([]interface{}, 0)
	points = append(points, 2.3, 4.2, 1.3, 9.3, 0.1, 1.1, 2.4)
	distanceFnc := func(point1, point2 interface{}) float64 { return math.Abs(point1.(float64) - point2.(float64)) }

	point := 3.
	k := 3

	want := make(map[interface{}]float64)

	want[2.3] = distanceFnc(point, 2.3)
	want[4.2] = distanceFnc(point, 4.2)
	want[2.4] = distanceFnc(point, 2.4)

	// act
	node := BuildTree(points, distanceFnc)
	got, _ := node.KNNSearch(point, uint(k))

	// assert
	eq := reflect.DeepEqual(want, got)

	if !eq {
		t.Errorf("want and got not equal")
	}
}

func TestRangeSearch(t *testing.T) {
	// arrange
	points := make([]interface{}, 0)
	points = append(points, 2.3, 4.2, 1.3, 9.3, 0.1, 1.1, 2.4)
	distanceFnc := func(point1, point2 interface{}) float64 { return math.Abs(point1.(float64) - point2.(float64)) }

	point := 3.
	threshold := 3.

	want := make(map[interface{}]float64)

	want[2.3] = distanceFnc(point, 2.3)
	want[4.2] = distanceFnc(point, 4.2)
	want[1.3] = distanceFnc(point, 1.3)
	want[0.1] = distanceFnc(point, 0.1)
	want[1.1] = distanceFnc(point, 1.1)
	want[2.4] = distanceFnc(point, 2.4)

	node := BuildTree(points, distanceFnc)
	got, _ := node.RangeSearch(point, threshold)

	// assert
	eq := reflect.DeepEqual(want, got)

	if !eq {
		t.Errorf("want and got not equal")
	}
}
