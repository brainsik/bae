package main

import (
	"fmt"
	"image/color"
	"math/cmplx"
	"math/rand"
	"runtime"
	"time"
)

// CalcStyle is an enum representing the type of calculation that will be performed.
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

// CalcPoint is the mapping between coordinate types.
type CalcPoint struct {
	z  complex128
	xy ImagePoint
}

// ColorResults maps image plane coordinates to a color.
type ColorResults map[ImagePoint]color.NRGBA

// CalcParams contains all the parameters needed to generate an image.
type CalcParams struct {
	plane *Plane

	style CalcStyle
	zf    ZFunc
	c     complex128

	cf  ColorFunc
	cfp ColorFuncParams

	iterations int
	limit      float64

	calc_area          PlaneView
	r_points, i_points int
}

func (cs CalcStyle) String() string {
	return CalcStyleName[int(cs)]
}

func (cp CalcPoint) String() string {
	return fmt.Sprintf("{%v, %v}", cp.z, cp.xy)
}

func (cp *CalcParams) String() string {
	return fmt.Sprintf(
		"CalcParams{\n%v\nStyle: %v\n%v\n%v\n%v\nc: %v\niterations: %v\nlimit: %v\n"+
			"calc area: %v\n"+
			"real points: %v (%v -> %v | %v)\nimag points: %v (%v -> %v | %v)\n}",
		cp.plane, cp.style, cp.zf, cp.cf, cp.cfp, cp.c, cp.iterations, cp.limit, cp.calc_area,
		cp.r_points, real(cp.calc_area.min), real(cp.calc_area.max), cp.calc_area.RealLen(),
		cp.i_points, imag(cp.calc_area.min), imag(cp.calc_area.max), cp.calc_area.ImagLen())
}

// NewAllPoints is a helper function that sets the calculation area of the
// CalcParams to the entire plane where the orbit of every point in the image
// is calculated.
func (cp CalcParams) NewAllPoints(iterations int, cf ColorFunc, cfp ColorFuncParams) CalcParams {
	if iterations <= 0 {
		iterations = cp.iterations
	}

	return CalcParams{
		// modified
		cf:         cf,
		cfp:        cfp,
		iterations: iterations,
		calc_area:  cp.plane.view,
		r_points:   cp.plane.ImageSize().width,
		i_points:   cp.plane.ImageSize().height,

		// unchanged
		plane: cp.plane,
		zf:    cp.zf,
		c:     cp.c,
		limit: cp.limit,
	}
}

// MakePlaneProblemSet returns a problem set for an even distribution of points in the calc_area.
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

// MakeImageProblemSet returns a problem set representing every point in the image plane.
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

// Calculate does the actual calculations for each point in the problem set.
func (cp *CalcParams) Calculate(problems []CalcPoint) (histogram CalcResults) {
	t_start := time.Now()
	calc_id := fmt.Sprintf("%p", problems)
	showed_progress := make(map[int]bool)

	img_width := cp.plane.ImageSize().width
	img_height := cp.plane.ImageSize().height

	// rz_min, rz_max := real(cp.plane.view.min), real(cp.plane.view.max)
	// iz_min, iz_max := imag(cp.plane.view.min), imag(cp.plane.view.max)

	var total_its, num_escaped, num_periodic uint
	histogram = make(CalcResults)
	f_zc := cp.zf.f

	for progress, pt := range problems {
		var z, c complex128
		if cp.style == Mandelbrot {
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
			// if real(z) < rz_min || real(z) > rz_max || imag(z) < iz_min || imag(z) > iz_max {
			// 	continue
			// }

			// Escaped?
			if cmplx.Abs(z) > cp.limit {
				if cp.style == Attractor {
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
				if cp.style == Attractor {
					histogram.Add(xy, z, 1).periodic = true
				} else {
					histogram.Add(pt.xy, pt.z, 1).periodic = true
				}
				num_periodic++
				// fmt.Printf("Point %v become periodic after %v iterations\n", z0, its)
				break
			}
			rag[z] = true

			if cp.style == Attractor {
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

// CalculateParallel breaks the problem set into chunks and runs conncurrent Calculate routines.
func (cp *CalcParams) CalculateParallel(concurrency int) (histogram CalcResults) {
	if concurrency == 0 {
		concurrency = int(1.5 * float64(runtime.NumCPU()))
	}

	var problems []CalcPoint
	if cp.style == Attractor {
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
			result_ch <- cp.Calculate(p_chunk)
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

// ColorImage sets image colors based on the results from Calculate.
func (cp *CalcParams) ColorImage(concurrency int) {
	histogram := cp.CalculateParallel(concurrency)
	histogram.PrintStats()

	t_start := time.Now()
	colors := cp.cf.f(histogram, cp.cfp)
	for pt, rgba := range colors {
		cp.plane.image.Set(pt.x, pt.y, rgba)
	}
	fmt.Printf("Image processing took %dms\n", time.Since(t_start).Milliseconds())
}
