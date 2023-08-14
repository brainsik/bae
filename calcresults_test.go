package main

import (
	"testing"
)

func TestCalcResultsMerge(t *testing.T) {
	dst := CalcResults{
		ImagePoint{0, 0}: &CalcResult{},
		ImagePoint{1, 1}: &CalcResult{complex(-1, 1), 1, true, false},
	}
	src := CalcResults{
		ImagePoint{0, 0}: &CalcResult{complex(-2, 2), 2, false, true},
		ImagePoint{1, 1}: &CalcResult{complex(-3, 3), 3, false, true},
	}

	dst.Merge(src)
	result := dst

	expected := CalcResults{
		ImagePoint{0, 0}: &CalcResult{0, 2, false, true},
		ImagePoint{1, 1}: &CalcResult{complex(-1, 1), 4, true, true},
	}

	for k := range result {
		if *result[k] != *expected[k] {
			t.Fatal(*result[k], *expected[k])
		}
	}
}

func TestCalcResultsAdd(t *testing.T) {
	crs := make(CalcResults)
	pt := ImagePoint{4, 2}
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
