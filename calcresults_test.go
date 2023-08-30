package main

import (
	"math"
	"testing"

	"github.com/brainsik/bae/plane"
)

func TestCalcResultsMerge(t *testing.T) {
	dst := CalcResults{
		plane.ImagePoint{X: 0, Y: 0}: &CalcResult{},
		plane.ImagePoint{X: 1, Y: 1}: &CalcResult{complex(-1, 1), 1, true, false},
	}
	src := CalcResults{
		plane.ImagePoint{X: 0, Y: 0}: &CalcResult{complex(-2, 2), 2, false, true},
		plane.ImagePoint{X: 1, Y: 1}: &CalcResult{complex(-3, 3), 3, false, true},
	}

	dst.Merge(src)
	result := dst

	expect := CalcResults{
		plane.ImagePoint{X: 0, Y: 0}: &CalcResult{0, 2, false, true},
		plane.ImagePoint{X: 1, Y: 1}: &CalcResult{complex(-1, 1), 4, true, true},
	}

	for k := range result {
		if *result[k] != *expect[k] {
			t.Error(*result[k], *expect[k])
		}
	}
}

func TestCalcResultsAdd(t *testing.T) {
	crs := make(CalcResults)
	pt := plane.ImagePoint{X: 4, Y: 2}
	z := complex(7, 7)
	val := uint(42)

	// Add new result
	crs.Add(pt, z, val)

	result, ok := crs[pt]
	if !ok {
		t.Fatalf("Expected %v to exist in CalcResults: %v", pt, crs)
	}
	if result.Val != val {
		t.Errorf("Expected result to have value %d, got %d", val, result.Val)
	}
	if result.Z != z {
		t.Errorf("Expected result to have value %v, got %v", z, result.Z)
	}

	// Add to existing result
	crs.Add(pt, -z, val)

	if result.Val != val*2 {
		t.Errorf("Expected result to have value %d, got %d", val*2, result.Val)
	}
	if result.Z != z {
		t.Errorf("Expected result to have value %v, got %v", z, result.Z)
	}
}

func TestCalcResultsMaxEscaped(t *testing.T) {
	crs := CalcResults{
		plane.ImagePoint{X: 0, Y: 0}: &CalcResult{complex(0, 0), 10, true, false},
		plane.ImagePoint{X: 1, Y: 1}: &CalcResult{complex(1, 1), 20, false, false},
	}
	result := crs.MaxEscaped()
	expect := 10.0
	if result != expect {
		t.Error(expect, result)
	}
}

func TestCalcResultsEmptyZero(t *testing.T) {
	crs := make(CalcResults)
	result := crs.Sum()
	expect := uint(0)
	if result != expect {
		t.Errorf("Expected result to be %v, got %v", expect, result)
	}
}

func TestCalcResultsEmptyNaN(t *testing.T) {
	crs := make(CalcResults)
	testCases := []struct {
		name string
		f    func() float64
	}{
		{"Max", crs.Max},
		{"MaxEscaped", crs.MaxEscaped},
		{"Min", crs.Min},
		{"Avg", crs.Avg},
		{"AvgEscaped", crs.AvgEscaped},
		{"Median", crs.Median},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.f()
			if !math.IsNaN(tc.f()) {
				t.Errorf("Expected result to be %v, got %v", math.NaN(), result)
			}
		})
	}
}
