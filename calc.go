package main

import (
	"fmt"
	"image/color"
	"math"
	"sort"
)

type CalcStyle int

const (
	Mandelbrot CalcStyle = iota
	Julia
	Attractor
)

type Calculator func(z, c complex128) complex128

type CalcPoint struct {
	z  complex128
	xy ImagePoint
}

func (cp CalcPoint) String() string {
	return fmt.Sprintf("{%v, %v}", cp.z, cp.xy)
}

type CalcResult struct {
	z   complex128
	val uint

	escaped, periodic bool
}

type CalcResults map[ImagePoint]*CalcResult

type Colorizer func(CalcResults) ColorResults

type ColorResults map[ImagePoint]color.NRGBA

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

func (cr CalcResults) PrintStats() {
	var min, max uint
	var periodic, escaped uint
	var sum, old_sum uint

	min = math.MaxUint
	values := []uint{}
	for _, v := range cr {
		old_sum = sum
		sum += v.val
		if sum < old_sum {
			fmt.Printf("CalcResults.PrintStats(): Sum overflowed!\n")
			return
		}

		if v.val < min {
			min = v.val
		}

		if v.val > max {
			max = v.val
		}

		if v.escaped {
			escaped++
		} else if v.periodic {
			periodic++
		}

		values = append(values, v.val)
	}
	avg := float64(sum) / float64(len(values))

	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	median := values[len(values)/2]

	escaped_pct := 100 * float64(escaped) / float64(len(cr))
	periodic_pct := 100 * float64(periodic) / float64(len(cr))

	fmt.Printf("[CalcResults] total: %d, max: %d, avg: %.1f, median: %d, min: %d, escaped: %d (%.1f%%), periodic: %d (%.1f%%)\n",
		len(cr), max, avg, median, min, escaped, escaped_pct, periodic, periodic_pct)
	return
}

func (cr CalcResults) Max() (max uint) {
	for _, v := range cr {
		if v.val > max {
			max = v.val
		}
	}
	return
}
