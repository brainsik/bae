package main

import (
	"math"
	"time"
)

const WIDTH = 2560
const ASPECT = 0.625

func main() {
	params := &coldwave1_allpts
	params.ColorImage(0)
	params.WritePNG("image.png")
}

func TimestampMilli() string {
	return time.Now().Format(time.StampMilli)
}

var klein = AttractorParams{
	plane: NewPlane(complex(0, -0.53), complex(2.6, 2.6*ASPECT), WIDTH),

	zf: burning_ship,
	cf: luma_ceil_20pct,
	c:  complex(-0.172, -1.136667),

	r_start: 0.0,
	r_end:   0.1,
	i_start: 0.0,
	i_end:   0.1,

	r_points: 1,
	i_points: 1,

	iterations: int(math.Pow(2, 24)),
	limit:      4,
}

var klein_allpts = AttractorParams{
	plane: NewPlane(complex(0, -0.53), complex(2.6, 2.6*ASPECT), WIDTH),

	zf: burning_ship,
	cf: luma_ceil_10pct,
	c:  complex(-0.172, -1.136667),

	r_start: -2.6,
	r_end:   2.6,
	i_start: -1.625,
	i_end:   1.625,

	r_points: WIDTH,
	i_points: int(WIDTH * ASPECT),

	iterations: 255,
	limit:      4,
}

var coldwave1 = AttractorParams{
	plane: NewPlane(complex(-0.2, 0.18), complex(1, 1*ASPECT), WIDTH),

	zf: burning_ship,
	cf: luma_ceil_6bit,
	c:  complex(0, 0),

	r_start: -0.53,
	r_end:   -0.46,

	i_start: -0.001,
	i_end:   0.499,

	r_points: 4,
	i_points: 10000,

	iterations: 1024,
	limit:      4,
}

var coldwave1_allpts = coldwave1.NewAllPoints(64, luma_ceil_8bit)

var coldwave2 = AttractorParams{
	plane: NewPlane(complex(-0.22, 0.15), complex(5, 5*ASPECT), WIDTH),

	zf: burning_ship,
	cf: luma_ceil_20pct,
	c:  complex(-0.1278, 0.0),

	r_start: -0.53,
	r_end:   -0.46,

	i_start: -0.255,
	i_end:   0.505,

	r_points: 1,
	i_points: 3040,

	iterations: 2048,
	limit:      4,
}

var coldwave2_allptsorig = AttractorParams{
	plane:   NewPlane(complex(-0.22, 0.15), complex(5, 5*ASPECT), WIDTH),
	zf:      burning_ship,
	cf:      luma_ceil_20pct,
	c:       complex(-0.1278, 0.0),
	r_start: -0.53,
	r_end:   -0.46,

	i_start: -0.255,
	i_end:   0.505,

	r_points: 1,
	i_points: 3040,

	iterations: 2048,
	limit:      4,
}
