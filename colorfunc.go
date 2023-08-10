package main

import (
	"fmt"
	"image/color"
	"math"
	"math/cmplx"
)

type ColorFunc struct {
	desc string
	f    func(CalcResults, ColorFuncParams) ColorResults
}

type ColorFuncParams struct {
	clip, gamma float64
	showclip    bool
}

func (cf ColorFunc) String() string {
	return fmt.Sprintf("ColorFunc: %s", cf.desc)
}

func (cfp ColorFuncParams) String() string {
	return fmt.Sprintf("ColorFuncParams{gamma:%f, clip:%f}", cfp.gamma, cfp.clip)
}

func GammaScale(val, max, gamma float64) float64 {
	if gamma <= 0 {
		gamma = 2.2 // default
	}

	scaled := complex(math.Min(val, max)/max, 0)
	gamma_correction := complex(1.0/gamma, 0)
	return real(cmplx.Pow(scaled, gamma_correction))
}

var cf_luma_clip_percent = ColorFunc{
	desc: `Brightness clips at given percent of max`,
	f: func(histogram CalcResults, params ColorFuncParams) ColorResults {
		coloring := make(ColorResults)
		max := (params.clip / 100) * float64(histogram.Max())
		for xy, v := range histogram {
			val := float64(v.val)
			if params.showclip && val >= max+1 {
				coloring[xy] = color.NRGBA{0xff, 0xd4, 0x79, 0xff}
			} else {
				brightness := uint8(255 * GammaScale(val, max, params.gamma))
				coloring[xy] = color.NRGBA{brightness, brightness, brightness, 0xff}
			}
		}
		return coloring
	},
}

var cf_luma_clip_value = ColorFunc{ //nolint:unused
	desc: `Brightness clips at given value`,
	f: func(histogram CalcResults, params ColorFuncParams) ColorResults {
		coloring := make(ColorResults)
		max := params.clip
		for xy, v := range histogram {
			val := float64(v.val)
			if params.showclip && val >= max+1 {
				coloring[xy] = color.NRGBA{0xff, 0xd4, 0x79, 0xff}
			} else {
				brightness := uint8(255 * GammaScale(val, max, params.gamma))
				coloring[xy] = color.NRGBA{brightness, brightness, brightness, 0xff}
			}
		}
		return coloring
	},
}

var cf_escaped_1bit = ColorFunc{ //nolint:unused
	desc: `Escaped points are white (1bit color)`,
	f: func(histogram CalcResults, params ColorFuncParams) ColorResults {
		coloring := make(ColorResults)
		for xy, v := range histogram {
			if v.escaped {
				coloring[xy] = color.NRGBA{0xff, 0xff, 0xff, 0xff}
			} else {
				coloring[xy] = color.NRGBA{0, 0, 0, 0xff}
			}
		}
		return coloring
	},
}

var cf_escaped_clip_percent = ColorFunc{ //nolint:unused
	desc: `Blue brightness depends on number of iterations to escape`,
	f: func(histogram CalcResults, params ColorFuncParams) ColorResults {
		coloring := make(ColorResults)
		max := (params.clip / 100) * float64(histogram.Max())
		for xy, v := range histogram {
			val := float64(v.val)
			if v.escaped {
				brightness := uint8(255 * GammaScale(val, max, params.gamma))
				coloring[xy] = color.NRGBA{brightness / 4, brightness / 4, brightness, 255}
			} else {
				coloring[xy] = color.NRGBA{0, 0, 0, 0xff}
			}
		}
		return coloring
	},
}
