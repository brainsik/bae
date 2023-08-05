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

	r_start, r_end     float64
	i_start, i_end     float64
	r_points, i_points int
}

func (ap *AttractorParams) String() string {
	return fmt.Sprintf(
		"AttractorParams{\n%v\n%v\n%v\nc: %v\niterations: %v\nlimit: %v\n"+
			"real points: %v (%v -> %v)\nimag points: %v (%v -> %v)\n}",
		ap.plane, ap.zf, ap.cf, ap.c, ap.iterations, ap.limit,
		ap.r_points, ap.r_start, ap.r_end, ap.i_points, ap.i_start, ap.i_end)
}

func (ap AttractorParams) NewAllPoints(iterations int, cf ColorFunc) AttractorParams {
	if iterations <= 0 {
		iterations = ap.iterations
	}

	return AttractorParams{
		cf:         cf,
		iterations: iterations,

		plane:    ap.plane,
		zf:       ap.zf,
		c:        ap.c,
		limit:    ap.limit,
		r_start:  real(ap.plane.view.min),
		r_end:    real(ap.plane.view.max),
		i_start:  imag(ap.plane.view.max),
		i_end:    imag(ap.plane.view.min),
		r_points: ap.plane.ImageSize().width,
		i_points: ap.plane.ImageSize().height,
	}
}

func (ap *AttractorParams) MakeProblemSet() (problems []CalcPoint) {
	r_step := (ap.r_end - ap.r_start) / float64(ap.r_points)
	i_step := (ap.i_end - ap.i_start) / float64(ap.i_points)
	for r := ap.r_start; r <= ap.r_end; r += r_step {
		for i := ap.i_start; i <= ap.i_end; i += i_step {
			z := complex(r, i)
			problems = append(problems, CalcPoint{z: z, xy: ap.plane.ImagePoint(z)})
		}
	}

	// Randomly distribute points to prevent calcuation hot spots from high iteration
	// points being near one another.
	rand.Shuffle(len(problems), func(a, b int) {
		problems[a], problems[b] = problems[b], problems[a]
	})
	return
}

func (ap *AttractorParams) Calculate(problems []CalcPoint) (histogram CalcResults) {
	t_start := time.Now()
	calc_id := fmt.Sprintf("%p", problems)
	showed_progress := make(map[int]bool)

	var total_its, num_escaped, num_periodic uint
	histogram = make(CalcResults)
	f_zc := ap.zf.f

	fmt.Printf("[%v] 🧠 Workin %v\n", time.Now().Format(time.StampMilli), calc_id)
	for progress, pt := range problems {
		// z0 := pt.z{}
		z := pt.z
		// Iterate point.
		rag := make(map[complex128]bool)
		for its := 0; its < ap.iterations; its++ {
			// Actual calculation.
			z = f_zc(z, ap.c)
			total_its++

			// Escaped?
			if cmplx.Abs(z) > ap.limit {
				num_escaped++
				// fmt.Printf("Point %v escaped after %v iterations\n", z0, its)
				break
			}

			// Periodic?
			if rag[z] {
				num_periodic++
				// fmt.Printf("Point %v become periodic after %v iterations\n", z0, its)
				break
			}
			rag[z] = true

			// Add iteration to histogram.
			xy := ap.plane.ImagePoint(z)
			hxy, ok := histogram[xy]
			if !ok {
				histogram[xy] = &CalcResult{z, 1}
			} else {
				hxy.Add(1)
			}
		}

		// Show progress.
		elapsed := time.Since(t_start).Seconds()
		if int(elapsed+1)%5 == 0 && !showed_progress[int(elapsed)] {
			percent := float64(progress) / float64(len(problems))
			t_remaining := (elapsed / percent) - elapsed
			fmt.Printf("[%v] ⌚️ %.0fs remaining for %v\n",
				time.Now().Format(time.StampMilli), t_remaining, calc_id)
			showed_progress[int(elapsed)] = true
		}
	}
	t_total := time.Since(t_start).Seconds()
	max_its := len(problems) * ap.iterations
	fmt.Printf("[%v] ✅ Finish %-30s\t%5.2fs (%.0f its/s) %d iterations (%1.f%%) - %d escaped, %d periodic\n",
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
		fmt.Printf("[%s] Got chunk of %d results\n", TimestampMilli(), len(r_chunk))
		for xy, cr := range r_chunk {
			hxy, ok := histogram[xy]
			if !ok {
				histogram[xy] = cr
			} else {
				hxy.Add(cr.val)
			}
		}
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
