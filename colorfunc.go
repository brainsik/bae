package main

import (
	"fmt"
	"image/color"
	"math"
	"math/cmplx"
)

// ColorFunc represents the alorithm used to determine the color of pixel in the image.
type ColorFunc struct {
	Desc string
	F    func(CalcResults, ColorFuncParams) ColorResults
}

// ColorFuncParams contains paramenters needed by a ColorFunc algorithm.
type ColorFuncParams struct {
	Clip, Gamma float64
	Showclip    bool
}

func (cf ColorFunc) String() string {
	return fmt.Sprintf("ColorFunc: %s", cf.Desc)
}

func (cfp ColorFuncParams) String() string {
	return fmt.Sprintf("ColorFuncParams{gamma:%f, clip:%f}", cfp.Gamma, cfp.Clip)
}

// GammaScale returns a scaled and gamma corrected brightness.
func GammaScale(val, max, gamma float64) float64 {
	if gamma <= 0 {
		gamma = 2.2 // default
	}

	scaled := complex(math.Min(val, max)/max, 0)
	gamma_correction := complex(1.0/gamma, 0)
	return real(cmplx.Pow(scaled, gamma_correction))
}

var cf_luma_clip_value = ColorFunc{ //nolint:unused
	Desc: `Brightness clips at given value`,
	F: func(histogram CalcResults, params ColorFuncParams) ColorResults {
		coloring := make(ColorResults)
		max := params.Clip
		for xy, v := range histogram {
			val := float64(v.Val)
			if params.Showclip && val >= max+1 {
				coloring[xy] = color.NRGBA{0xff, 0xd4, 0x79, 0xff}
			} else {
				brightness := uint8(255 * GammaScale(val, max, params.Gamma))
				coloring[xy] = color.NRGBA{brightness, brightness, brightness, 0xff}
			}
		}
		return coloring
	},
}

var cf_luma_clip_percent_avg = ColorFunc{ //nolint:unused
	Desc: `Brightness clips at given percent of max`,
	F: func(histogram CalcResults, params ColorFuncParams) ColorResults {
		coloring := make(ColorResults)
		max := (params.Clip / 100) * histogram.Avg()
		for xy, v := range histogram {
			val := float64(v.Val)
			if params.Showclip && val >= max+1 {
				coloring[xy] = color.NRGBA{0xff, 0xd4, 0x79, 0xff}
			} else {
				brightness := uint8(255 * GammaScale(val, max, params.Gamma))
				coloring[xy] = color.NRGBA{brightness, brightness, brightness, 0xff}
			}
		}
		return coloring
	},
}

var cf_luma_clip_percent_max = ColorFunc{ //nolint:unused
	Desc: `Brightness clips at given percent of max`,
	F: func(histogram CalcResults, params ColorFuncParams) ColorResults {
		coloring := make(ColorResults)
		max := (params.Clip / 100) * histogram.Max()
		for xy, v := range histogram {
			val := float64(v.Val)
			if params.Showclip && val >= max+1 {
				coloring[xy] = color.NRGBA{0xff, 0xd4, 0x79, 0xff}
			} else {
				brightness := uint8(255 * GammaScale(val, max, params.Gamma))
				coloring[xy] = color.NRGBA{brightness, brightness, brightness, 0xff}
			}
		}
		return coloring
	},
}

var cf_escaped_1bit = ColorFunc{ //nolint:unused
	Desc: `Escaped points are white (1bit color)`,
	F: func(histogram CalcResults, params ColorFuncParams) ColorResults {
		coloring := make(ColorResults)
		for xy, v := range histogram {
			if v.Escaped {
				coloring[xy] = color.NRGBA{0xff, 0xff, 0xff, 0xff}
			} else {
				coloring[xy] = color.NRGBA{0, 0, 0, 0xff}
			}
		}
		return coloring
	},
}

var cf_escaped_clip_value = ColorFunc{ //nolint:unused
	Desc: `Blue brightness depends on number of iterations to escape`,
	F: func(histogram CalcResults, params ColorFuncParams) ColorResults {
		coloring := make(ColorResults)
		max := params.Clip
		for xy, v := range histogram {
			val := float64(v.Val)
			if v.Escaped {
				luma := GammaScale(val, max, params.Gamma)
				rg := uint8(math.Min(255*luma*luma, 255))
				b := uint8(math.Min(255*math.Sqrt(luma), 255))
				coloring[xy] = color.NRGBA{rg, rg, b, 0xff}
			} else {
				coloring[xy] = color.NRGBA{0, 0, 0, 0xff}
			}
		}
		return coloring
	},
}

var cf_escaped_clip_percent_avg = ColorFunc{ //nolint:unused
	Desc: `Blue brightness depends on number of iterations to escape`,
	F: func(histogram CalcResults, params ColorFuncParams) ColorResults {
		coloring := make(ColorResults)
		max := (params.Clip / 100) * histogram.AvgEscaped()
		for xy, v := range histogram {
			val := float64(v.Val)
			if v.Escaped {
				luma := GammaScale(val, max, params.Gamma)
				rg := uint8(math.Min(255*luma*luma, 255))
				b := uint8(math.Min(255*math.Sqrt(luma), 255))
				coloring[xy] = color.NRGBA{rg, rg, b, 0xff}
			} else {
				coloring[xy] = color.NRGBA{0, 0, 0, 0xff}
			}
		}
		return coloring
	},
}

var cf_escaped_clip_percent_max = ColorFunc{ //nolint:unused
	Desc: `Blue brightness depends on number of iterations to escape`,
	F: func(histogram CalcResults, params ColorFuncParams) ColorResults {
		coloring := make(ColorResults)
		max := (params.Clip / 100) * histogram.MaxEscaped()
		for xy, v := range histogram {
			val := float64(v.Val)
			if v.Escaped {
				luma := GammaScale(val, max, params.Gamma)
				rg := uint8(math.Min(255*luma*luma, 255))
				b := uint8(math.Min(255*math.Sqrt(luma), 255))
				coloring[xy] = color.NRGBA{rg, rg, b, 0xff}
			} else {
				coloring[xy] = color.NRGBA{0, 0, 0, 0xff}
			}
		}
		return coloring
	},
}
