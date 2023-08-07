package main

import (
	"fmt"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"time"
)

type CalcStyle int

const (
	Attractor CalcStyle = iota
	Julia
	Mandelbrot
)

var CalcStyleName = map[int]string{
	int(Attractor):  "Attractor",
	int(Julia):      "Julia",
	int(Mandelbrot): "Mandelbrot",
}

type CalcPoint struct {
	z  complex128
	xy ImagePoint
}

type CalcResult struct {
	z   complex128
	val uint

	escaped, periodic bool
}

type CalcResults map[ImagePoint]*CalcResult

type ColorResults map[ImagePoint]color.NRGBA

type CalcParams struct {
	plane *Plane

	zf ZFunc
	cf ColorFunc

	c complex128

	iterations int
	limit      float64

	calc_area          PlaneView
	r_points, i_points int
}

func (cp CalcPoint) String() string {
	return fmt.Sprintf("{%v, %v}", cp.z, cp.xy)
}

func (cr *CalcResult) Add(n uint) {
	prev := cr.val
	cr.val += n

	if cr.val <= prev && n > 0 {
		fmt.Printf("CalcResult %v Overflowed: setting value to max\n", cr)
		cr.val = math.MaxUint
	}
}

func (cr CalcResults) Add(xy ImagePoint, z complex128, val uint) *CalcResult {
	cr_xy, ok := cr[xy]
	if !ok {
		cr[xy] = &CalcResult{z, val, false, false}
	} else {
		cr_xy.Add(1)
	}
	return cr[xy]
}

func (cr CalcResults) Merge(src CalcResults) {
	for k, v := range src {
		dst, ok := cr[k]
		if !ok {
			cr[k] = v
		} else {
			dst.val += v.val
			dst.escaped = dst.escaped || v.escaped
			dst.periodic = dst.periodic || v.periodic
		}
	}
}

func (cr CalcResults) PrintStats() {
	var min, max uint
	var periodic, escaped uint
	var sum, old_sum uint

	min = math.MaxUint
	values := []uint{}
	for _, v := range cr {
		old_sum = sum
		sum += v.val
		if sum < old_sum {
			fmt.Printf("CalcResults.PrintStats(): Sum overflowed!\n")
			return
		}

		if v.val < min {
			min = v.val
		}

		if v.val > max {
			max = v.val
		}

		if v.escaped {
			escaped++
		} else if v.periodic {
			periodic++
		}

		values = append(values, v.val)
	}
	avg := float64(sum) / float64(len(values))

	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	median := values[len(values)/2]

	escaped_pct := 100 * float64(escaped) / float64(len(cr))
	periodic_pct := 100 * float64(periodic) / float64(len(cr))

	fmt.Printf("[CalcResults] total: %d, max: %d, avg: %.1f, median: %d, min: %d, escaped: %d (%.1f%%), periodic: %d (%.1f%%)\n",
		len(cr), max, avg, median, min, escaped, escaped_pct, periodic, periodic_pct)
	return
}

func (cr CalcResults) Max() (max uint) {
	for _, v := range cr {
		if v.val > max {
			max = v.val
		}
	}
	return
}

func (cp *CalcParams) String() string {
	return fmt.Sprintf(
		"AttractorParams{\n%v\n%v\n%v\nc: %v\niterations: %v\nlimit: %v\n"+
			"calc area: %v\n"+
			"real points: %v (%v -> %v | %v)\nimag points: %v (%v -> %v | %v)\n}",
		cp.plane, cp.zf, cp.cf, cp.c, cp.iterations, cp.limit, cp.calc_area,
		cp.r_points, real(cp.calc_area.min), real(cp.calc_area.max), cp.calc_area.RealLen(),
		cp.i_points, imag(cp.calc_area.min), imag(cp.calc_area.max), cp.calc_area.ImagLen())
}

func (cp CalcParams) NewAllPoints(iterations int, cf ColorFunc) CalcParams {
	if iterations <= 0 {
		iterations = cp.iterations
	}

	return CalcParams{
		// modified
		cf:         cf,
		iterations: iterations,
		calc_area:  cp.plane.view,

		// unchanged
		plane:    cp.plane,
		zf:       cp.zf,
		c:        cp.c,
		limit:    cp.limit,
		r_points: cp.plane.ImageSize().width,
		i_points: cp.plane.ImageSize().height,
	}
}

func (cp *CalcParams) MakePlaneProblemSet() (problems []CalcPoint) {
	// TODO: Return an Error instead?
	if cp.calc_area.RealLen() > 0 && cp.r_points <= 1 {
		panic("Undefined how to make a single point problem on a non-single point line." +
			"Either add more points or set the length of real(calc_area) to be a single point.")
	}
	if cp.calc_area.ImagLen() > 0 && cp.i_points <= 1 {
		panic("Undefined how to make a single point problem on a non-single point line." +
			"Either add more points or set the length of imag(calc_area) to be a single point.")
	}

	r_step := cp.calc_area.RealLen() / float64(cp.r_points-1)
	i_step := cp.calc_area.ImagLen() / float64(cp.i_points-1)

	r := real(cp.calc_area.min)
	for r_pt := 0; r_pt < cp.r_points; r_pt++ {
		i := imag(cp.calc_area.min)
		for i_pt := 0; i_pt < cp.i_points; i_pt++ {
			z := complex(r, i)
			xy := cp.plane.ImagePoint(z)
			problems = append(problems, CalcPoint{z: z, xy: xy})
			i += i_step
		}
		r += r_step
	}
	return
}

func (cp *CalcParams) MakeImageProblemSet() (problems []CalcPoint) {
	for x := 0; x < cp.plane.ImageSize().width; x++ {
		for y := 0; y < cp.plane.ImageSize().height; y++ {
			xy := ImagePoint{x: x, y: y}
			z := cp.plane.PlanePoint(xy)
			problems = append(problems, CalcPoint{z: z, xy: xy})
		}
	}
	return
}

func (cp *CalcParams) Calculate(problems []CalcPoint, style CalcStyle) (histogram CalcResults) {
	t_start := time.Now()
	calc_id := fmt.Sprintf("%p", problems)
	showed_progress := make(map[int]bool)

	img_width := cp.plane.ImageSize().width
	img_height := cp.plane.ImageSize().height

	var total_its, num_escaped, num_periodic uint
	histogram = make(CalcResults)
	f_zc := cp.zf.f

	for progress, pt := range problems {
		var z, c complex128
		if style == Mandelbrot {
			z = complex(0, 0)
			c = pt.z
		} else {
			z = pt.z
			c = cp.c
		}

		rag := make(map[complex128]bool)
		for its := 0; its < cp.iterations; its++ {
			total_its++

			z = f_zc(z, c)
			xy := cp.plane.ImagePoint(z)

			// Escaped?
			if cmplx.Abs(z) > cp.limit {
				if style == Attractor {
					histogram.Add(xy, z, 1).escaped = true
				} else {
					histogram.Add(pt.xy, pt.z, 1).escaped = true
				}
				num_escaped++
				// fmt.Printf("Point %v escaped after %v iterations\n", z0, its)
				break
			}

			// Periodic?
			if rag[z] {
				if style == Attractor {
					histogram.Add(xy, z, 1).periodic = true
				} else {
					histogram.Add(pt.xy, pt.z, 1).periodic = true
				}
				num_periodic++
				// fmt.Printf("Point %v become periodic after %v iterations\n", z0, its)
				break
			}
			rag[z] = true

			if style == Attractor {
				// Only add to histogram if pixel is in the image plane.
				if xy.x >= 0 && xy.x <= img_width && xy.y >= 0 && xy.y <= img_height {
					histogram.Add(xy, z, 1)
				}
			} else {
				histogram.Add(pt.xy, pt.z, 1)
			}
		}

		// Show progress.
		elapsed := time.Since(t_start).Seconds()
		if int(elapsed+1)%6 == 0 && !showed_progress[int(elapsed)] {
			percent := float64(progress) / float64(len(problems))
			t_remaining := (elapsed / percent) - elapsed
			fmt.Printf("[%v] ⌚️ Workin %v %6.0fs remaining\n",
				time.Now().Format(time.StampMilli), calc_id, t_remaining)
			showed_progress[int(elapsed)] = true
		}
	}
	t_total := time.Since(t_start).Seconds()
	max_its := len(problems) * cp.iterations
	fmt.Printf("[%v] ✅ Finish %s %6.0fs (%.0f its/s) • %d its (%1.f%%) • %d escaped, %d periodic\n",
		TimestampMilli(), calc_id, t_total,
		float64(total_its)/t_total, total_its, 100*float64(total_its)/float64(max_its),
		num_escaped, num_periodic)

	return
}

func (cp *CalcParams) CalculateParallel(concurrency int, style CalcStyle) (histogram CalcResults) {
	if concurrency == 0 {
		concurrency = int(1.5 * float64(runtime.NumCPU()))
	}

	var problems []CalcPoint
	if style == Attractor {
		problems = cp.MakePlaneProblemSet()
	} else {
		problems = cp.MakeImageProblemSet()
	}

	if len(problems) < concurrency {
		concurrency = len(problems)
	}

	chunk_size := len(problems) / concurrency

	// Randomly distribute points to prevent calcuation hot spots.
	rand.Shuffle(len(problems), func(a, b int) {
		problems[a], problems[b] = problems[b], problems[a]
	})

	fmt.Printf("%v\n\n", cp)
	fmt.Printf("Logical CPUs: %v (will use %v concurrent routines)\n", runtime.NumCPU(), concurrency)
	fmt.Printf("Orbits to calculate: %d (~%d per routine)\n", len(problems), chunk_size)

	// Calculate points.
	result_ch := make(chan CalcResults)
	for chunk_n := 0; chunk_n < concurrency; chunk_n++ {
		chunk_start := chunk_n * chunk_size
		chunk_end := chunk_start + chunk_size

		// final chunk includes remainder of fractional division
		if chunk_n+1 == concurrency {
			chunk_end = len(problems)
		}

		p_chunk := problems[chunk_start:chunk_end]
		fmt.Printf("[%v] 🚀 Launch %p | orbits %d - %d\n",
			time.Now().Format(time.StampMilli), p_chunk, chunk_start, chunk_end)
		go func() {
			result_ch <- cp.Calculate(p_chunk, style)
		}()
	}

	// Collect results.
	histogram = make(CalcResults)
	chunks_received := 0
	for r_chunk := range result_ch {
		histogram.Merge(r_chunk)
		chunks_received++
		if chunks_received >= concurrency {
			break
		}
	}

	return
}

func (cp *CalcParams) ColorImage(concurrency int, style CalcStyle) {
	histogram := cp.CalculateParallel(concurrency, style)
	histogram.PrintStats()

	colors := cp.cf.f(histogram)
	for pt, rgba := range colors {
		cp.plane.image.Set(pt.x, pt.y, rgba)
	}
}

func (cp *CalcParams) WritePNG(filename string) {
	png_file, _ := os.Create(filename)
	penc := png.Encoder{CompressionLevel: png.BestCompression}
	penc.Encode(png_file, cp.plane.image)
	fmt.Printf("Wrote %s\n", png_file.Name())
	png_file.Close()
}
