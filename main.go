package main

import (
	"log"
	"math"

	"github.com/brainsik/bae/plane"
)

const (
	HEIGHT = 800
	ASPECT = 1.6
)

func main() {
	log.SetFlags(log.Lshortfile)

	params := coldwave2
	params.ColorImage()
	params.Plane.WritePNG("image.png")
}

// Single orbit attractor.
var klein = NewCalcParams(CalcParams{ //nolint:unused
	Plane: plane.NewPlane(complex(-0.1, -0.54), complex(1.6*ASPECT, 1.6), HEIGHT).WithInverted(),

	Style:      Attractor,
	ZF:         zf_klein,
	C:          complex(-0.172, -1.136667),
	Iterations: int(math.Pow(2, 20)),

	CalcArea: plane.PlaneView{Min: complex(-0.4, -0.1), Max: complex(0.4, 0.1)},
	RPoints:  3,
	IPoints:  3,

	CF:  cf_luma_clip_percent_max,
	CFP: ColorFuncParams{Clip: 10},
})

// Multi-orbit map.
var klein2_allpts = NewCalcParams(CalcParams{ //nolint:unused
	Plane: plane.NewPlane(complex(-0.2, -1), complex(2.2*ASPECT, 2.2), HEIGHT),

	Style: Attractor,
	ZF:    zf_klein2,
	C:     complex(-0.172, -1.136667),
}.NewAllPoints(128, cf_luma_clip_percent_max, ColorFuncParams{Clip: 8, Gamma: 2.0}))

// Multi-orbit map.
var coldwave1 = NewCalcParams(CalcParams{ //nolint:unused
	Plane: plane.NewPlane(complex(-0.19, 0.19), complex(0.8*ASPECT, 0.8), HEIGHT),

	Style:      Attractor,
	ZF:         zf_klein,
	C:          complex(0, 0),
	Iterations: 1024,

	CalcArea: plane.PlaneView{Min: complex(-0.53, -0.001), Max: complex(-0.46, 0.499)},
	RPoints:  8,
	IPoints:  25000,

	CF:  cf_luma_clip_value,
	CFP: ColorFuncParams{Clip: math.Pow(2, 6)},
})

// Multi-orbit map.
var coldwave1_allpts = NewCalcParams(coldwave1.NewAllPoints( //nolint:unused
	16, cf_luma_clip_value, ColorFuncParams{Clip: math.Pow(2, 8)}))

// Single orbit attractor.
var coldwave2 = NewCalcParams(CalcParams{ //nolint:unused
	Plane: plane.NewPlane(complex(-0.22, -0.175), complex(3.75*ASPECT, 3.75), HEIGHT),

	Style:      Attractor,
	ZF:         zf_klein,
	C:          complex(-0.1278, 0.0),
	Iterations: 4096,

	CalcArea: plane.PlaneView{Min: complex(-0.5, -0.255), Max: complex(-0.5, 0.505)},
	RPoints:  1,
	IPoints:  3080,

	CF:  cf_luma_clip_percent_max,
	CFP: ColorFuncParams{Clip: 10},
})

var coldwave2_julia = NewCalcParams(CalcParams{ //nolint:unused
	Plane: plane.NewPlane(complex(-0.22, -0.175), complex(3.75*ASPECT, 3.75), HEIGHT),

	Style:      Julia,
	ZF:         zf_klein,
	C:          complex(-0.1278, 0.0),
	Iterations: 96,

	CF:  cf_escaped_clip_percent_max,
	CFP: ColorFuncParams{Clip: 50},
})

var julia_classic = NewCalcParams(CalcParams{ //nolint:unused
	Plane: plane.NewPlane(complex(0, 0), complex(4*ASPECT, 4), HEIGHT),

	Style:      Julia,
	ZF:         zf_mandelbrot,
	C:          complex(0.285, 0.01),
	Iterations: 493,

	CF:  cf_escaped_clip_percent_max,
	CFP: ColorFuncParams{Clip: 50},
})

var burning_ship = NewCalcParams(CalcParams{ //nolint:unused
	// Plane: plane.NewPlane(complex(1.75, 0.038), complex(0.145, 0.145*ASPECT_INV), HEIGHT),
	Plane: plane.NewPlane(complex(-1.765, -0.035), complex(0.15*ASPECT, 0.15), HEIGHT).WithInverted(),

	Style:      Mandelbrot,
	ZF:         zf_burning_ship,
	Iterations: 256,

	CF:  cf_escaped_clip_percent_avg,
	CFP: ColorFuncParams{Clip: 400},
})

var mandelbrot = NewCalcParams(CalcParams{ //nolint:unused
	Plane: plane.NewPlane(complex(-0.5, 0), complex(4*ASPECT, 4), HEIGHT),

	Style:      Mandelbrot,
	ZF:         zf_mandelbrot,
	Iterations: 256,

	CF:  cf_escaped_clip_percent_avg,
	CFP: ColorFuncParams{Clip: 400},
})
