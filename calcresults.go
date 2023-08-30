package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/brainsik/bae/plane"
)

// CalcResult are the calculation results for a point z.
type CalcResult struct {
	Z   complex128
	Val uint

	Escaped, Periodic bool
}

// CalcResults maps an ImagePoint to the corresponding CalcResult for that point.
type CalcResults map[plane.ImagePoint]*CalcResult

// Add adds n, printing a warning if the value overflows.
func (cr *CalcResult) Add(n uint) {
	prev := cr.Val
	cr.Val += n

	if cr.Val <= prev && n > 0 {
		fmt.Printf("Warning: CalcResult %v Overflowed: setting value to max\n", cr)
		cr.Val = math.MaxUint
	}
}

// Add adds val to the existing val or creates a new CalcResult with val.
func (cr CalcResults) Add(xy plane.ImagePoint, z complex128, val uint) *CalcResult {
	cr_xy, ok := cr[xy]
	if !ok {
		cr[xy] = &CalcResult{z, val, false, false}
	} else {
		cr_xy.Add(val)
	}
	return cr[xy]
}

// Merge makes a union with src by adding vals and setting booleans to the OR'd value.
func (cr CalcResults) Merge(src CalcResults) {
	for k, v := range src {
		dst, ok := cr[k]
		if !ok {
			cr[k] = v
		} else {
			dst.Val += v.Val
			dst.Escaped = dst.Escaped || v.Escaped
			dst.Periodic = dst.Periodic || v.Periodic
		}
	}
}

/* Statistics */

// Max returns the highest val.
func (cr CalcResults) Max() float64 {
	if len(cr) <= 0 {
		return math.NaN()
	}
	var max uint
	for _, v := range cr {
		if v.Val > max {
			max = v.Val
		}
	}
	return float64(max)
}

// MaxEscaped returns the highest escaped val.
func (cr CalcResults) MaxEscaped() float64 {
	if len(cr) <= 0 {
		return math.NaN()
	}
	var max uint
	for _, v := range cr {
		if !v.Escaped {
			continue
		}
		if v.Val > max {
			max = v.Val
		}
	}
	return float64(max)
}

// Min returns the lowest val.
func (cr CalcResults) Min() float64 {
	if len(cr) <= 0 {
		return math.NaN()
	}
	var min uint = math.MaxUint
	for _, v := range cr {
		if v.Val < min {
			min = v.Val
		}
	}
	return float64(min)
}

// Sum returns the sum of all vals.
func (cr CalcResults) Sum() (sum uint) {
	for _, v := range cr {
		sum += v.Val
	}
	return
}

// Avg returns the average of all vals.
func (cr CalcResults) Avg() float64 {
	return float64(cr.Sum()) / float64(len(cr))
}

// AvgEscaped returns the average of all escaped vals.
func (cr CalcResults) AvgEscaped() float64 {
	var num, sum float64
	for _, v := range cr {
		if v.Escaped {
			sum += float64(v.Val)
			num++
		}
	}

	return sum / num
}

// Median returns the median of the sorted vals.
func (cr CalcResults) Median() float64 {
	if len(cr) <= 0 {
		return math.NaN()
	}

	var vals []uint
	for _, v := range cr {
		vals = append(vals, v.Val)
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	return float64(vals[len(vals)/2])
}

// PrintStats outputs a variety of stats about what's in the CalcResults.
func (cr CalcResults) PrintStats() {
	var periodic, escaped uint
	for _, v := range cr {
		if v.Escaped {
			escaped++
		} else if v.Periodic {
			periodic++
		}
	}
	escaped_pct := 100 * float64(escaped) / float64(len(cr))
	periodic_pct := 100 * float64(periodic) / float64(len(cr))

	fmt.Printf("[CalcResults] "+
		"total: %d, max(esc): %.0f(%.0f), avg(esc): %.1f(%.1f), median: %.0f, min: %.0f, "+
		"escaped: %d (%.1f%%), periodic: %d (%.1f%%)\n",
		len(cr), cr.Max(), cr.MaxEscaped(), cr.Avg(), cr.AvgEscaped(), cr.Median(), cr.Min(),
		escaped, escaped_pct, periodic, periodic_pct)
}
