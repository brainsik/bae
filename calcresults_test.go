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
