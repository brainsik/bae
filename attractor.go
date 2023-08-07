package main

import (
	"fmt"
	"image/png"
	"math/cmplx"
	"math/rand"
	"os"
	"runtime"
	"time"
)

type AttractorParams struct {
	plane *Plane

	zf ZFunc
	cf ColorFunc

	c complex128

	iterations int
	limit      float64

	calc_area          PlaneView
	r_points, i_points int
}

func (ap *AttractorParams) String() string {
	return fmt.Sprintf(
		"AttractorParams{\n%v\n%v\n%v\nc: %v\niterations: %v\nlimit: %v\n"+
			"calc area: %v\n"+
			"real points: %v (%v -> %v | %v)\nimag points: %v (%v -> %v | %v)\n}",
		ap.plane, ap.zf, ap.cf, ap.c, ap.iterations, ap.limit, ap.calc_area,
		ap.r_points, real(ap.calc_area.min), real(ap.calc_area.max), ap.calc_area.RealLen(),
		ap.i_points, imag(ap.calc_area.min), imag(ap.calc_area.max), ap.calc_area.ImagLen())
}

func (ap AttractorParams) NewAllPoints(iterations int, cf ColorFunc) AttractorParams {
	if iterations <= 0 {
		iterations = ap.iterations
	}

	return AttractorParams{
		// modified
		cf:         cf,
		iterations: iterations,
		calc_area:  ap.plane.view,

		// unchanged
		plane:    ap.plane,
		zf:       ap.zf,
		c:        ap.c,
		limit:    ap.limit,
		r_points: ap.plane.ImageSize().width,
		i_points: ap.plane.ImageSize().height,
	}
}

func (ap *AttractorParams) MakeProblemSet() (problems []CalcPoint) {
	// TODO: Return an Error instead?
	if ap.calc_area.RealLen() > 0 && ap.r_points <= 1 {
		panic("Undefined how to make a single point problem on a non-single point line." +
			"Either add more points or set the length of real(calc_area) to be a single point.")
	}
	if ap.calc_area.ImagLen() > 0 && ap.i_points <= 1 {
		panic("Undefined how to make a single point problem on a non-single point line." +
			"Either add more points or set the length of imag(calc_area) to be a single point.")
	}

	r_step := ap.calc_area.RealLen() / float64(ap.r_points-1)
	i_step := ap.calc_area.ImagLen() / float64(ap.i_points-1)

	r := real(ap.calc_area.min)
	for r_pt := 0; r_pt < ap.r_points; r_pt++ {
		i := imag(ap.calc_area.min)
		for i_pt := 0; i_pt < ap.i_points; i_pt++ {
			z := complex(r, i)
			xy := ap.plane.ImagePoint(z)
			problems = append(problems, CalcPoint{z: z, xy: xy})
			i += i_step
		}
		r += r_step
	}
	return
}

func (ap *AttractorParams) Calculate(problems []CalcPoint) (histogram CalcResults) {
	t_start := time.Now()
	calc_id := fmt.Sprintf("%p", problems)
	showed_progress := make(map[int]bool)

	var total_its, num_escaped, num_periodic uint
	histogram = make(CalcResults)
	f_zc := ap.zf.f

	// fmt.Printf("[%v] 🧠 Workin %v\n", time.Now().Format(time.StampMilli), calc_id)
	for progress, pt := range problems {
		// Iterate point.
		z := pt.z
		rag := make(map[complex128]bool)
		for its := 0; its < ap.iterations; its++ {
			total_its++

			z = f_zc(z, ap.c)
			xy := ap.plane.ImagePoint(z)

			// Escaped?
			if cmplx.Abs(z) > ap.limit {
				histogram.Add(xy, z, 1).escaped = true
				num_escaped++
				// fmt.Printf("Point %v escaped after %v iterations\n", z0, its)
				break
			}

			// Periodic?
			if rag[z] {
				histogram.Add(xy, z, 1).periodic = true
				num_periodic++
				// fmt.Printf("Point %v become periodic after %v iterations\n", z0, its)
				break
			}
			rag[z] = true

			histogram.Add(xy, z, 1)
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
	max_its := len(problems) * ap.iterations
	fmt.Printf("[%v] ✅ Finish %s %6.0fs (%.0f its/s) • %d its (%1.f%%) • %d escaped, %d periodic\n",
		TimestampMilli(), calc_id, t_total,
		float64(total_its)/t_total, total_its, 100*float64(total_its)/float64(max_its),
		num_escaped, num_periodic)

	return
}

func (ap *AttractorParams) CalculateParallel(concurrency int) (histogram CalcResults) {
	if concurrency == 0 {
		concurrency = int(1.5 * float64(runtime.NumCPU()))
	}

	problems := ap.MakeProblemSet()
	chunk_size := len(problems) / concurrency

	// Randomly distribute points to prevent calcuation hot spots.
	rand.Shuffle(len(problems), func(a, b int) {
		problems[a], problems[b] = problems[b], problems[a]
	})

	fmt.Printf("%v\n\n", ap)
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
			result_ch <- ap.Calculate(p_chunk)
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

func (ap *AttractorParams) ColorImage(concurrency int) {
	var histogram CalcResults
	if concurrency == 1 {
		histogram = ap.Calculate(ap.MakeProblemSet())
	} else {
		histogram = ap.CalculateParallel(concurrency)
	}
	histogram.PrintStats()

	colors := ap.cf.f(histogram)

	for pt, rgba := range colors {
		ap.plane.image.Set(pt.x, pt.y, rgba)
	}
}

func (ap *AttractorParams) WritePNG(filename string) {
	png_file, _ := os.Create(filename)
	penc := png.Encoder{CompressionLevel: png.BestCompression}
	penc.Encode(png_file, ap.plane.image)
	fmt.Printf("Wrote %s\n", png_file.Name())
	png_file.Close()
}
