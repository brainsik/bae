package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

// ZFunc represents the math function f(z, c).
type ZFunc struct {
	Desc string
	F    func(z, c complex128) complex128
}

func (zf ZFunc) String() string {
	return fmt.Sprintf("ZFunc: %s", zf.Desc)
}

var zf_burning_ship = ZFunc{ //nolint:unused
	Desc: `Burning Ship :: (|x| + i|y|)^2 + c`,
	F: func(z, c complex128) complex128 {
		return cmplx.Pow(-complex(math.Abs(real(z)), math.Abs(imag(z))), 2.0) - c
	},
}

var zf_klein = ZFunc{ //nolint:unused
	Desc: `Klein :: z^2 - y + i|x| + c`,
	F: func(z, c complex128) complex128 {
		return cmplx.Pow(z, 2.0) + complex(-imag(z), math.Abs(real(z))) + c
	},
}

var zf_klein2 = ZFunc{ //nolint:unused
	Desc: `Klein :: z^2 + |y| + ix + c`,
	F: func(z, c complex128) complex128 {
		return cmplx.Pow(z, 2.0) + complex(math.Abs(imag(z)), real(z)) + c
	},
}

var zf_mandelbrot = ZFunc{ //nolint:unused
	Desc: `Mandelbrot :: z^2 + c`,
	F: func(z, c complex128) complex128 {
		return cmplx.Pow(z, 2.0) + c
	},
}
