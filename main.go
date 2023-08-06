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
	plane: NewPlane(complex(-0.2, -0.15), complex(0.9, 0.9*ASPECT), WIDTH),

	zf: burning_ship,
	cf: luma_ceil_6bit,
	c:  complex(0, 0),

	calc_area: PlaneView{complex(-0.53, -0.001), complex(-0.46, 0.499)},
	r_points:  8,
	i_points:  25000,

	iterations: 1024,
	limit:      4,
}

var coldwave1_allpts = coldwave1.NewAllPoints(64, luma_ceil_8bit)

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
