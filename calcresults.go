package main

import (
	"fmt"
	"math"
	"sort"
)

type CalcResult struct {
	z   complex128
	val uint

	escaped, periodic bool
}

type CalcResults map[ImagePoint]*CalcResult

func (cr *CalcResult) Add(n uint) {
	prev := cr.val
	cr.val += n

	if cr.val <= prev && n > 0 {
		fmt.Printf("CalcResult %v Overflowed: setting value to max\n", cr)
		cr.val = math.MaxUint
	}
}

func (cr CalcResults) Add(xy ImagePoint, z complex128, val uint) *CalcResult {
	cr_xy, ok := cr[xy]
	if !ok {
		cr[xy] = &CalcResult{z, val, false, false}
	} else {
		cr_xy.Add(1)
	}
	return cr[xy]
}

func (cr CalcResults) Merge(src CalcResults) {
	for k, v := range src {
		dst, ok := cr[k]
		if !ok {
			cr[k] = v
		} else {
			dst.val += v.val
			dst.escaped = dst.escaped || v.escaped
			dst.periodic = dst.periodic || v.periodic
		}
	}
}

func (cr CalcResults) Max() (max uint) {
	for _, v := range cr {
		if v.val > max {
			max = v.val
		}
	}
	return
}

func (cr CalcResults) Min() (min uint) {
	min = math.MaxUint
	for _, v := range cr {
		if v.val < min {
			min = v.val
		}
	}
	return
}

func (cr CalcResults) Sum() (sum uint) {
	for _, v := range cr {
		sum += v.val
	}
	return
}

func (cr CalcResults) Avg() float64 {
	return float64(cr.Sum()) / float64(len(cr))
}

func (cr CalcResults) Median() uint {
	var vals []uint
	for _, v := range cr {
		vals = append(vals, v.val)
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	return vals[len(vals)/2]
}

func (cr CalcResults) PrintStats() {
	var periodic, escaped uint
	for _, v := range cr {
		if v.escaped {
			escaped++
		} else if v.periodic {
			periodic++
		}
	}
	escaped_pct := 100 * float64(escaped) / float64(len(cr))
	periodic_pct := 100 * float64(periodic) / float64(len(cr))

	fmt.Printf("[CalcResults] "+
		"total: %d, max: %d, avg: %.1f, median: %d, min: %d, "+
		"escaped: %d (%.1f%%), periodic: %d (%.1f%%)\n",
		len(cr), cr.Max(), cr.Avg(), cr.Median(), cr.Min(),
		escaped, escaped_pct, periodic, periodic_pct)
	return
}
