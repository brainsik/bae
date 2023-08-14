package main

import (
	"math"
	"time"
)

const (
	WIDTH  = 3840
	ASPECT = 0.625
)

func main() {
	params := &coldwave2
	params.ColorImage(3)
	params.Plane.WritePNG("image.png")
}

func TimestampMilli() string {
	return time.Now().Format(time.StampMilli)
}

var klein = CalcParams{ //nolint:unused
	Plane: NewPlane(complex(-0.1, -0.54), complex(2.6, 2.6*ASPECT), WIDTH),

	Style: Attractor,
	ZF:    zf_klein,
	C:     complex(-0.172, -1.136667),

	CF:  cf_luma_clip_percent,
	CFP: ColorFuncParams{Clip: 10},

	CalcArea: PlaneView{complex(0, 0), complex(0, 0)},
	RPoints:  1,
	IPoints:  1,

	Iterations: int(math.Pow(2, 23)),
	Limit:      4,
}

var klein2_allpts = CalcParams{ //nolint:unused
	Plane: NewPlane(complex(-0.2, -1), complex(2.2, 2.2*ASPECT), WIDTH),

	Style: Attractor,
	ZF:    zf_klein2,
	C:     complex(-0.172, -1.136667),

	Limit: 4,
}.NewAllPoints(128, cf_luma_clip_percent, ColorFuncParams{Clip: 8, Gamma: 2.0})

var coldwave1 = CalcParams{ //nolint:unused
	Plane: NewPlane(complex(-0.19, 0.19), complex(0.8, 0.8*ASPECT), WIDTH),

	Style: Attractor,
	ZF:    zf_klein,
	C:     complex(0, 0),

	CF:  cf_luma_clip_value,
	CFP: ColorFuncParams{Clip: math.Pow(2, 6)},

	CalcArea: PlaneView{complex(-0.53, -0.001), complex(-0.46, 0.499)},
	RPoints:  8,
	IPoints:  25000,

	Iterations: 1024,
	Limit:      4,
}

var coldwave1_allpts = coldwave1.NewAllPoints( //nolint:unused
	16, cf_luma_clip_value, ColorFuncParams{Clip: math.Pow(2, 8)})

var coldwave2 = CalcParams{ //nolint:unused
	Plane: NewPlane(complex(-0.22, -0.175), complex(3.75, 3.75*ASPECT), WIDTH),

	Style: Attractor,
	ZF:    zf_klein,
	C:     complex(-0.1278, 0.0),

	CF:  cf_luma_clip_percent,
	CFP: ColorFuncParams{Clip: 10},

	CalcArea: PlaneView{complex(-0.5, -0.255), complex(-0.5, 0.505)},
	RPoints:  1,
	IPoints:  3080,

	Iterations: 4096,
	Limit:      4,
}

var burning_ship = CalcParams{ //nolint:unused
	Plane: NewPlane(complex(1.75, 0.038), complex(0.145, 0.145*ASPECT), WIDTH),

	Style: Mandelbrot,
	ZF:    zf_burning_ship,

	CF:  cf_escaped_clip_percent,
	CFP: ColorFuncParams{Clip: 20},

	Iterations: 1024,
	Limit:      4,
}

var julia_classic = CalcParams{ //nolint:unused
	Plane: NewPlane(complex(0, 0), complex(4, 4*ASPECT), WIDTH),

	Style: Julia,
	ZF:    zf_mandelbrot,
	C:     complex(0.285, 0.01),

	CF:  cf_escaped_clip_percent,
	CFP: ColorFuncParams{Clip: 50, Gamma: 2.8},

	Iterations: 493,
	Limit:      4,
}
