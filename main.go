package main

import (
	"math"
	"time"
)

const WIDTH = 1920
const ASPECT = 0.625

func main() {
	params := &coldwave1
	params.ColorImage(6, Julia)
	params.WritePNG("image.png")
}

func TimestampMilli() string {
	return time.Now().Format(time.StampMilli)
}

var mandelbrot = AttractorParams{
	plane: NewPlane(complex(-0.75, 0), complex(3, 3*ASPECT), WIDTH),

	zf: mandelbrot_f,
	cf: escaped_blue,
	c:  complex(0, 0), // unused in mandelbrot

	calc_area: PlaneView{complex(-2, -1.25), complex(2, 1.25)}, // unused in mandelbrot (for now)
	r_points:  WIDTH,
	i_points:  WIDTH * ASPECT,

	iterations: 64,
	limit:      4,
}

var julia = AttractorParams{
	plane: NewPlane(complex(0, 0), complex(4, 4*ASPECT), WIDTH),

	zf: mandelbrot_f,
	cf: escaped_blue,
	c:  complex(0.285, 0.01),

	calc_area: PlaneView{complex(-2, -1.25), complex(2, 1.25)}, // unused in mandelbrot/julia (for now)
	r_points:  WIDTH,
	i_points:  WIDTH * ASPECT,

	iterations: 280,
	limit:      4,
}

var klein = AttractorParams{
	plane: NewPlane(complex(0, -1.5), complex(2.6, 2.6*ASPECT), WIDTH),

	zf: burning_ship,
	cf: luma_ceil_10pct,
	c:  complex(-0.172, -1.136667),

	calc_area: PlaneView{complex(0, 0), complex(0, 0)},
	r_points:  1,
	i_points:  1,

	iterations: int(math.Pow(2, 23)),
	limit:      4,
}

var coldwave1 = AttractorParams{
	plane: NewPlane(complex(-0.19, -0.1), complex(0.8, 0.8*ASPECT), WIDTH),

	zf: burning_ship,
	cf: luma_ceil_6bit,
	c:  complex(0, 0),

	calc_area: PlaneView{complex(-0.53, -0.001), complex(-0.46, 0.499)},
	r_points:  8,
	i_points:  25000,

	iterations: 1024,
	limit:      4,
}

var coldwave1_allpts = coldwave1.NewAllPoints(16, luma_ceil_8bit)

var coldwave2 = AttractorParams{
	plane: NewPlane(complex(-0.22, -1.65), complex(4, 4*ASPECT), WIDTH),

	zf: burning_ship,
	cf: luma_ceil_20pct,
	c:  complex(-0.1278, 0.0),

	calc_area: PlaneView{complex(-0.5, -0.255), complex(-0.5, 0.505)},
	r_points:  1,
	i_points:  3080,

	iterations: 2048,
	limit:      4,
}
