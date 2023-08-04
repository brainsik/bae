package main

import (
	"fmt"
	"image/color"
	"math"
	"math/cmplx"
)

type ZFunc struct {
	desc string
	f    func(z, c complex128) complex128
}

type ColorFunc struct {
	desc string
	f    func(CalcResults) ColorResults
}

type Calculator func(z, c complex128) complex128

type CalcPoint struct {
	z  complex128
	xy ImagePoint
}

type CalcResult struct {
	z   complex128
	val uint
}

type CalcResults map[ImagePoint]*CalcResult

type Colorizer func(CalcResults) ColorResults

type ColorResults map[ImagePoint]color.NRGBA

func (cr *CalcResult) Add(n uint) {
	prev := cr.val
	cr.val += n

	if cr.val <= prev && n > 0 {
		fmt.Printf("CalcResult %v Overflowed: setting value to max\n", cr)
		cr.val = ^uint(0)
	}
}

func (cr CalcResults) Max() (max uint) {
	for _, v := range cr {
		if v.val > max {
			max = v.val
		}
	}
	fmt.Printf("CalcResults max value: %v\n", max)
	return
}

var burning_ship = ZFunc{
	desc: `Burning Ship: z^2 - y + i|x| + c`,
	f: func(z, c complex128) complex128 {
		return cmplx.Pow(z, 2.0) + complex(-imag(z), math.Abs(real(z))) + c
	},
}
