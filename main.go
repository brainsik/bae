package main

import (
	"math"
	"time"
)

const WIDTH = 1920
const ASPECT = 0.625

func main() {
	params := &klein
	params.ColorImage(6)
	params.WritePNG("image.png")
}

func TimestampMilli() string {
	return time.Now().Format(time.StampMilli)
}

var klein = CalcParams{
	plane: NewPlane(complex(0, -1.5), complex(2.6, 2.6*ASPECT), WIDTH),

	style: Attractor,
	zf:    zf_klein,
	c:     complex(-0.172, -1.136667),

	cf:  cf_luma_clip_percent,
	cfp: ColorFuncParams{clip: 10},

	calc_area: PlaneView{complex(0, 0), complex(0, 0)},
	r_points:  1,
	i_points:  1,

	iterations: int(math.Pow(2, 23)),
	limit:      4,
}

var coldwave1 = CalcParams{
	plane: NewPlane(complex(-0.19, -0.1), complex(0.8, 0.8*ASPECT), WIDTH),

	style: Attractor,
	zf:    zf_klein,
	c:     complex(0, 0),

	cf:  cf_luma_clip_value,
	cfp: ColorFuncParams{clip: math.Pow(2, 6)},

	calc_area: PlaneView{complex(-0.53, -0.001), complex(-0.46, 0.499)},
	r_points:  8,
	i_points:  25000,

	iterations: 1024,
	limit:      4,
}

var coldwave1_allpts = coldwave1.NewAllPoints(
	16, cf_luma_clip_value, ColorFuncParams{clip: math.Pow(2, 8)})

var coldwave2 = CalcParams{
	plane: NewPlane(complex(-0.22, -1.65), complex(4, 4*ASPECT), WIDTH),

	style: Attractor,
	zf:    zf_klein,
	c:     complex(-0.1278, 0.0),

	cf:  cf_luma_clip_percent,
	cfp: ColorFuncParams{clip: 10},

	calc_area: PlaneView{complex(-0.5, -0.255), complex(-0.5, 0.505)},
	r_points:  1,
	i_points:  3080,

	iterations: 2048,
	limit:      4,
}

var burning_ship = CalcParams{
	plane: NewPlane(complex(1.75, 0.038), complex(0.145, 0.145*ASPECT), WIDTH),

	style: Mandelbrot,
	zf:    zf_burning_ship,

	cf:  cf_escaped_clip_percent,
	cfp: ColorFuncParams{clip: 20},

	iterations: 1024,
	limit:      4,
}

var julia_classic = CalcParams{
	plane: NewPlane(complex(0, 0), complex(4, 4*ASPECT), WIDTH),

	style: Julia,
	zf:    zf_mandelbrot,
	c:     complex(0.285, 0.01),

	cf:  cf_escaped_clip_percent,
	cfp: ColorFuncParams{clip: 45},

	iterations: 493,
	limit:      4,
}
