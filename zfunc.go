package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

type ZFunc struct {
	desc string
	f    func(z, c complex128) complex128
}

func (zf ZFunc) String() string {
	return fmt.Sprintf("ZFunc: %s", zf.desc)
}

var zf_burning_ship = ZFunc{
	desc: `Burning Ship :: (|x| + i|y|)^2 + c`,
	f: func(z, c complex128) complex128 {
		return cmplx.Pow(-complex(math.Abs(real(z)), math.Abs(imag(z))), 2.0) - c
	},
}

var zf_klein = ZFunc{
	desc: `Klein :: z^2 - y + i|x| + c`,
	f: func(z, c complex128) complex128 {
		return cmplx.Pow(z, 2.0) + complex(-imag(z), math.Abs(real(z))) + c
	},
}

var zf_mandelbrot = ZFunc{
	desc: `Mandelbrot :: z^2 + c`,
	f: func(z, c complex128) complex128 {
		return cmplx.Pow(z, 2.0) + c
	},
}
