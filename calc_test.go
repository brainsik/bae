package main

import (
	"sort"
	"testing"

	"github.com/brainsik/bae/plane"
)

func TestMakePlaneProblemSetSinglePointPanic(t *testing.T) {
	params := CalcParams{
		Plane:    plane.NewPlane(complex(0, 0), complex(2, 2), 10),
		CalcArea: plane.PlaneView{Min: complex(-1, -1), Max: complex(1, 1)},
		RPoints:  1,
		IPoints:  1,
	}

	defer func() { _ = recover() }()
	params.MakePlaneProblemSet() // should panic
	t.Errorf("Expected MakeAttractorProblemSet to panic.")
}

func NewPlane(c1, c2 complex128, i int) {
	panic("unimplemented")
}

func TestMakePlaneProblemSetMultiplePoints(t *testing.T) {
	params := CalcParams{
		Plane:    plane.NewPlane(complex(0, 0), complex(2, 2), 10),
		CalcArea: plane.PlaneView{Min: complex(-1, -1), Max: complex(1, 1)},
		RPoints:  3,
		IPoints:  3,
	}

	// Get.
	expect := []CalcPoint{
		{(-1 + 1i), plane.ImagePoint{X: 0, Y: 0}},
		{(-1 + 0i), plane.ImagePoint{X: 0, Y: 5}},
		{(-1 - 1i), plane.ImagePoint{X: 0, Y: 10}},

		{(0 + 1i), plane.ImagePoint{X: 5, Y: 0}},
		{(0 + 0i), plane.ImagePoint{X: 5, Y: 5}},
		{(0 - 1i), plane.ImagePoint{X: 5, Y: 10}},

		{(1 + 1i), plane.ImagePoint{X: 10, Y: 0}},
		{(1 + 0i), plane.ImagePoint{X: 10, Y: 5}},
		{(1 - 1i), plane.ImagePoint{X: 10, Y: 10}},
	}
	result := params.MakePlaneProblemSet()
	if len(result) != len(expect) {
		t.Fatalf("len(result) = %d; want %d", len(result), len(expect))
	}

	// Sort.
	less := func(cp []CalcPoint, i, j int) bool {
		ix := cp[i].XY.X
		jx := cp[j].XY.X
		switch {
		case ix < jx:
			return true
		case ix == jx:
			return cp[i].XY.Y < cp[j].XY.Y
		default: // ix > jx
			return false
		}
	}
	sort.Slice(result, func(i, j int) bool { return less(result, i, j) })
	sort.Slice(expect, func(i, j int) bool { return less(expect, i, j) })

	// Compare.
	for i := 0; i < len(expect); i++ {
		if result[i] != expect[i] {
			t.Errorf("Result and expect differed at index %d: %v != %v", i, result[i], expect[i])
		}
	}
}

func TestMakePlaneProblemSetSinglePoint(t *testing.T) {
	params := CalcParams{
		Plane:    plane.NewPlane(complex(0, 0), complex(2, 2), 10),
		CalcArea: plane.PlaneView{Min: complex(-1, -1), Max: complex(-1, -1)},
		RPoints:  1,
		IPoints:  1,
	}

	// Get.
	expect := []CalcPoint{{(-1 - 1i), plane.ImagePoint{X: 0, Y: 10}}}
	result := params.MakePlaneProblemSet()
	if len(result) != 1 {
		t.Fatalf("len(result) = %d; want 1", len(result))
	}

	// Compare.
	if result[0] != expect[0] {
		t.Error(result[0], expect[0])
	}
}
