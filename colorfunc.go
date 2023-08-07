package main

import (
	"fmt"
	"image/color"
	"math"
	"math/cmplx"
)

type ColorFunc struct {
	desc string
	f    func(CalcResults) ColorResults
}

func (cf ColorFunc) String() string {
	return fmt.Sprintf("ColorFunc: %s", cf.desc)
}

func scale(val uint, max float64) float64 {
	scaled := math.Min(float64(val), max) / max
	// Gamma 2.2 correction.
	return real(cmplx.Pow(complex(scaled, 0), 1.0/2.2))
}

var escaped_1bit = ColorFunc{
	desc: `1bit coloring: escaped points are white`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		for xy, v := range histogram {
			if v.escaped {
				coloring[xy] = color.NRGBA{255, 255, 255, 255}
			} else {
				coloring[xy] = color.NRGBA{0, 0, 0, 255}
			}
		}
		return
	},
}

var escaped_blue = ColorFunc{
	desc: `1bit coloring: escaped points are white`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		max := float64(histogram.Max())
		for xy, v := range histogram {
			if v.escaped {
				brightness := uint8(255 * scale(v.val, max))
				coloring[xy] = color.NRGBA{brightness / 4, brightness / 4, brightness, 255}
			} else {
				coloring[xy] = color.NRGBA{0, 0, 0, 255}
			}
		}
		return
	},
}

var luma = ColorFunc{
	desc: `Gamma corrected brightness`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		max := float64(histogram.Max())
		for xy, v := range histogram {
			brightness := uint8(255 * scale(v.val, max))
			coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
		}
		return
	},
}

var luma_ceil_12bit = ColorFunc{
	desc: `Gamma corrected brightness with original values clipped at 1023`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		for xy, v := range histogram {
			brightness := uint8(255 * scale(v.val, 4095))
			coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
		}
		return
	},
}

var luma_ceil_10bit = ColorFunc{
	desc: `Gamma corrected brightness with original values clipped at 1023`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		for xy, v := range histogram {
			brightness := uint8(255 * scale(v.val, 1023))
			coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
		}
		return
	},
}

var luma_ceil_10bit_showclip = ColorFunc{
	desc: `Gamma corrected brightness with original values clipped at 1023 and shown`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		for xy, v := range histogram {
			if v.val >= 1024 {
				coloring[xy] = color.NRGBA{0xFF, 0xD4, 0x79, 255}
			} else {
				brightness := uint8(255 * scale(v.val, 1023))
				coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
			}
		}
		return
	},
}

var luma_ceil_8bit = ColorFunc{
	desc: `Gamma corrected brightness with original values clipped at 255`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		for xy, v := range histogram {
			brightness := uint8(255 * scale(v.val, 255))
			coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
		}
		return
	},
}

var luma_ceil_6bit = ColorFunc{
	desc: `Gamma corrected brightness with original values clipped at 63`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		for xy, v := range histogram {
			brightness := uint8(255 * scale(v.val, 63))
			coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
		}
		return
	},
}

var luma_ceil_80pct = ColorFunc{
	desc: `Gamma corrected brightness with original values clipped at 80% of max`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		max := 0.80 * float64(histogram.Max())
		fmt.Printf("Scaled max: %.1f\n", max)
		for xy, v := range histogram {
			brightness := uint8(255 * scale(v.val, max))
			coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
		}
		return
	},
}

var luma_ceil_40pct = ColorFunc{
	desc: `Gamma corrected brightness with original values clipped at 40% of max`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		max := 0.40 * float64(histogram.Max())
		fmt.Printf("Scaled max: %.1f\n", max)
		for xy, v := range histogram {
			brightness := uint8(255 * scale(v.val, max))
			coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
		}
		return
	},
}

var luma_ceil_20pct = ColorFunc{
	desc: `Gamma corrected brightness with original values clipped at 20% of max`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		max := 0.20 * float64(histogram.Max())
		fmt.Printf("Scaled max: %.1f\n", max)
		for xy, v := range histogram {
			brightness := uint8(255 * scale(v.val, max))
			coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
		}
		return
	},
}

var luma_ceil_10pct = ColorFunc{
	desc: `Gamma corrected brightness with original values clipped at 10% of max`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		max := 0.1 * float64(histogram.Max())
		fmt.Printf("Scaled max: %.1f\n", max)
		for xy, v := range histogram {
			brightness := uint8(255 * scale(v.val, max))
			coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
		}
		return
	},
}

var luma_ceil_sqrt = ColorFunc{
	desc: `Gamma corrected brightness with original values clipped at sqrt(max)`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		max := math.Sqrt(float64(histogram.Max()))
		fmt.Printf("Scaled max: %.1f\n", max)
		for xy, v := range histogram {
			brightness := uint8(255 * scale(v.val, max))
			coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
		}
		return
	},
}

var luma_ceil_log = ColorFunc{
	desc: `Gamma corrected brightness with original values clipped at log(max)`,
	f: func(histogram CalcResults) (coloring ColorResults) {
		coloring = make(ColorResults)
		max := math.Log(float64(histogram.Max()))
		fmt.Printf("Scaled max: %.1f\n", max)
		for xy, v := range histogram {
			brightness := uint8(255 * scale(v.val, max))
			coloring[xy] = color.NRGBA{brightness, brightness, brightness, 255}
		}
		return
	},
}
