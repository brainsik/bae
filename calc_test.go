package main

import (
	"sort"
	"testing"
)

func TestMakePlaneProblemSetSinglePointPanic(t *testing.T) {
	params := CalcParams{
		Plane:    NewPlane(complex(0, 0), complex(2, 2), 10),
		CalcArea: PlaneView{complex(-1, -1), complex(1, 1)},
		RPoints:  1,
		IPoints:  1,
	}

	defer func() { _ = recover() }()
	params.MakePlaneProblemSet() // should panic
	t.Errorf("Expected MakeAttractorProblemSet to panic.")
}

func TestMakePlaneProblemSetMultiplePoints(t *testing.T) {
	params := CalcParams{
		Plane:    NewPlane(complex(0, 0), complex(2, 2), 10),
		CalcArea: PlaneView{complex(-1, -1), complex(1, 1)},
		RPoints:  3,
		IPoints:  3,
	}

	// Get.
	expect := []CalcPoint{
		{(-1 + 1i), ImagePoint{0, 0}},
		{(-1 + 0i), ImagePoint{0, 5}},
		{(-1 - 1i), ImagePoint{0, 10}},

		{(0 + 1i), ImagePoint{5, 0}},
		{(0 + 0i), ImagePoint{5, 5}},
		{(0 - 1i), ImagePoint{5, 10}},

		{(1 + 1i), ImagePoint{10, 0}},
		{(1 + 0i), ImagePoint{10, 5}},
		{(1 - 1i), ImagePoint{10, 10}},
	}
	result := params.MakePlaneProblemSet()
	if len(result) != len(expect) {
		t.Fatalf("len(result) = %d; want %d", len(result), len(expect))
	}

	// Sort.
	less := func(cp []CalcPoint, i, j int) bool {
		ix := cp[i].XY.x
		jx := cp[j].XY.x
		switch {
		case ix < jx:
			return true
		case ix == jx:
			return cp[i].XY.y < cp[j].XY.y
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
		Plane:    NewPlane(complex(0, 0), complex(2, 2), 10),
		CalcArea: PlaneView{complex(-1, -1), complex(-1, -1)},
		RPoints:  1,
		IPoints:  1,
	}

	// Get.
	expect := []CalcPoint{{(-1 - 1i), ImagePoint{0, 10}}}
	result := params.MakePlaneProblemSet()
	if len(result) != 1 {
		t.Fatalf("len(result) = %d; want 1", len(result))
	}

	// Compare.
	if result[0] != expect[0] {
		t.Error(result[0], expect[0])
	}
}
