package main

import (
	"sort"
	"testing"
)

func TestMakePlaneProblemSetSinglePointPanic(t *testing.T) {
	ap := AttractorParams{
		plane:     NewPlane(complex(0, 0), complex(2, 2), 10),
		calc_area: PlaneView{complex(-1, -1), complex(1, 1)},
		r_points:  1,
		i_points:  1,
	}

	defer func() { _ = recover() }()
	ap.MakePlaneProblemSet() // should panic
	t.Errorf("Expected MakeAttractorProblemSet to panic.")
}

func TestMakePlaneProblemSetMultiplePoints(t *testing.T) {
	ap := AttractorParams{
		plane:     NewPlane(complex(0, 0), complex(2, 2), 10),
		calc_area: PlaneView{complex(-1, -1), complex(1, 1)},
		r_points:  3,
		i_points:  3,
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
	result := ap.MakePlaneProblemSet()
	if len(result) != len(expect) {
		t.Fatalf("len(result) = %d; want %d", len(result), len(expect))
	}

	// Sort.
	less := func(cp []CalcPoint, i, j int) bool {
		ix := cp[i].xy.x
		jx := cp[j].xy.x
		switch {
		case ix < jx:
			return true
		case ix == jx:
			return cp[i].xy.y < cp[j].xy.y
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
	ap := AttractorParams{
		plane:     NewPlane(complex(0, 0), complex(2, 2), 10),
		calc_area: PlaneView{complex(-1, -1), complex(-1, -1)},
		r_points:  1,
		i_points:  1,
	}

	// Get.
	expect := []CalcPoint{{(-1 - 1i), ImagePoint{0, 10}}}
	result := ap.MakePlaneProblemSet()
	if len(result) != 1 {
		t.Fatalf("len(result) = %d; want 1", len(result))
	}

	// Compare.
	if result[0] != expect[0] {
		t.Error(result[0], expect[0])
	}
}
