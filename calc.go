package main

import (
	"fmt"
	"image/color"
	"log"
	"math/cmplx"
	"math/rand"
	"runtime"
	"time"

	"github.com/brainsik/bae/plane"
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
	Z  complex128
	XY plane.ImagePoint
}

// ColorResults maps image plane coordinates to a color.
type ColorResults map[plane.ImagePoint]color.NRGBA

// CalcParams contains all the parameters needed to generate an image.
type CalcParams struct {
	Plane *plane.Plane

	Style      CalcStyle
	ZF         ZFunc
	C          complex128
	Iterations int
	Limit      float64

	CalcArea         plane.PlaneView
	RPoints, IPoints int

	Concurrency int

	CF  ColorFunc
	CFP ColorFuncParams
}

func (cs CalcStyle) String() string {
	return CalcStyleName[int(cs)]
}

func (cp CalcPoint) String() string {
	return fmt.Sprintf("{%v, %v}", cp.Z, cp.XY)
}

func (cp *CalcParams) String() string {
	return fmt.Sprintf(
		"CalcParams{\n%v\nStyle: %v\n%v\n%v\n%v\nc: %v\niterations: %v\nlimit: %v\n"+
			"calc area: %v\n"+
			"real points: %v (%v -> %v | %v)\nimag points: %v (%v -> %v | %v)\n"+
			"concurrency: %d\n}",
		cp.Plane, cp.Style, cp.ZF, cp.CF, cp.CFP, cp.C, cp.Iterations, cp.Limit, cp.CalcArea,
		cp.RPoints, real(cp.CalcArea.Min), real(cp.CalcArea.Max), cp.CalcArea.RealLen(),
		cp.IPoints, imag(cp.CalcArea.Min), imag(cp.CalcArea.Max), cp.CalcArea.ImagLen(),
		cp.Concurrency)
}

// NewCalcParams returns a new CalcParams object based on the given one.
// Some defaults are set for zeroed fields.
func NewCalcParams(cp CalcParams) *CalcParams {
	limit := cp.Limit
	if limit == 0 {
		limit = cmplx.Abs(cp.Plane.GetSize())
	}

	return &CalcParams{
		Plane: cp.Plane,

		Style:      cp.Style,
		ZF:         cp.ZF,
		C:          cp.C,
		Iterations: cp.Iterations,
		Limit:      limit,

		CalcArea: cp.CalcArea,
		RPoints:  cp.RPoints,
		IPoints:  cp.IPoints,

		Concurrency: cp.Concurrency,

		CF:  cp.CF,
		CFP: cp.CFP,
	}
}

// NewAllPoints is a helper function that sets the calculation area of the
// CalcParams to the entire plane where the orbit of every point in the image
// is calculated.
func (cp CalcParams) NewAllPoints(iterations int, cf ColorFunc, cfp ColorFuncParams) CalcParams {
	if iterations <= 0 {
		iterations = cp.Iterations
	}

	return CalcParams{
		// modified
		CF:         cf,
		CFP:        cfp,
		Iterations: iterations,
		CalcArea:   cp.Plane.GetView(),
		RPoints:    cp.Plane.ImageWidth(),
		IPoints:    cp.Plane.ImageHeight(),

		// unchanged
		Plane:       cp.Plane,
		ZF:          cp.ZF,
		C:           cp.C,
		Limit:       cp.Limit,
		Concurrency: cp.Concurrency,
	}
}

// MakePlaneProblemSet returns a problem set for an even distribution of points in the calc_area.
func (cp *CalcParams) MakePlaneProblemSet() (problems []CalcPoint) {
	t_start := time.Now()

	// TODO: Return an Error instead?
	if cp.CalcArea.RealLen() > 0 && cp.RPoints <= 1 {
		panic("Undefined how to make a single point problem on a non-single point line." +
			"Either add more points or set the length of real(calc_area) to be a single point.")
	}
	if cp.CalcArea.ImagLen() > 0 && cp.IPoints <= 1 {
		panic("Undefined how to make a single point problem on a non-single point line." +
			"Either add more points or set the length of imag(calc_area) to be a single point.")
	}

	r_step := cp.CalcArea.RealLen() / float64(cp.RPoints-1)
	i_step := cp.CalcArea.ImagLen() / float64(cp.IPoints-1)

	r := real(cp.CalcArea.Min)
	for r_pt := 0; r_pt < cp.RPoints; r_pt++ {
		i := imag(cp.CalcArea.Min)
		for i_pt := 0; i_pt < cp.IPoints; i_pt++ {
			z := complex(r, i)
			xy := cp.Plane.ToImagePoint(z)
			problems = append(problems, CalcPoint{Z: z, XY: xy})
			i += i_step
		}
		r += r_step
	}

	log.Printf("Took %dms to make problem set\n", time.Since(t_start).Milliseconds())
	return
}

// MakeImageProblemSet returns a problem set representing every point in the image plane.
func (cp *CalcParams) MakeImageProblemSet() (problems []CalcPoint) {
	t_start := time.Now()
	for x := 0; x < cp.Plane.ImageWidth(); x++ {
		for y := 0; y < cp.Plane.ImageHeight(); y++ {
			xy := plane.ImagePoint{X: x, Y: y}
			z := cp.Plane.ToComplexPoint(xy)
			problems = append(problems, CalcPoint{Z: z, XY: xy})
		}
	}
	log.Printf("Took %dms to make problem set\n", time.Since(t_start).Milliseconds())
	return
}

// Calculate does the actual calculations for each point in the problem set.
func (cp *CalcParams) Calculate(problems []CalcPoint) (histogram CalcResults) {
	t_start := time.Now()
	calc_id := fmt.Sprintf("%p", problems)
	showed_progress := make(map[int]bool)

	img_width := cp.Plane.ImageWidth()
	img_height := cp.Plane.ImageHeight()

	// rz_min, rz_max := real(cp.plane.view.min), real(cp.plane.view.max)
	// iz_min, iz_max := imag(cp.plane.view.min), imag(cp.plane.view.max)

	var total_its, num_escaped, num_periodic uint
	histogram = make(CalcResults)
	f_zc := cp.ZF.F

	for progress, pt := range problems {
		var z, c complex128
		if cp.Style == Mandelbrot {
			z = complex(0, 0)
			c = pt.Z
		} else {
			z = pt.Z
			c = cp.C
		}

		rag := make(map[complex128]bool)
		for its := 0; its < cp.Iterations; its++ {
			total_its++

			z = f_zc(z, c)
			xy := cp.Plane.ToImagePoint(z)
			// if real(z) < rz_min || real(z) > rz_max || imag(z) < iz_min || imag(z) > iz_max {
			// 	continue
			// }

			// Escaped?
			if cmplx.Abs(z) > cp.Limit {
				if cp.Style == Attractor {
					histogram.Add(xy, z, 1).Escaped = true
				} else {
					histogram.Add(pt.XY, pt.Z, 1).Escaped = true
				}
				num_escaped++
				// fmt.Printf("Point %v escaped after %v iterations\n", z0, its)
				break
			}

			// Periodic?
			if rag[z] {
				if cp.Style == Attractor {
					histogram.Add(xy, z, 1).Periodic = true
				} else {
					histogram.Add(pt.XY, pt.Z, 1).Periodic = true
				}
				num_periodic++
				// fmt.Printf("Point %v become periodic after %v iterations\n", z0, its)
				break
			}
			rag[z] = true

			if cp.Style == Attractor {
				// Only add to histogram if pixel is in the image plane.
				if xy.X >= 0 && xy.X <= img_width && xy.Y >= 0 && xy.Y <= img_height {
					histogram.Add(xy, z, 1)
				}
			} else {
				histogram.Add(pt.XY, pt.Z, 1)
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
	max_its := len(problems) * cp.Iterations
	fmt.Printf("[%v] ✅ Finish %s %6.0fs (%.0f its/s) • %d its (%1.f%%) • %d escaped, %d periodic\n",
		TimestampMilli(), calc_id, t_total,
		float64(total_its)/t_total, total_its, 100*float64(total_its)/float64(max_its),
		num_escaped, num_periodic)

	return
}

func TimestampMilli() string {
	return time.Now().Format(time.StampMilli)
}

// CalculateParallel breaks the problem set into chunks and runs conncurrent Calculate routines.
func (cp *CalcParams) CalculateParallel() (histogram CalcResults) {
	concurrency := cp.Concurrency
	if cp.Concurrency == 0 {
		concurrency = int(1.5 * float64(runtime.NumCPU()))
	}

	var problems []CalcPoint
	if cp.Style == Attractor {
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
func (cp *CalcParams) ColorImage() {
	histogram := cp.CalculateParallel()
	histogram.PrintStats()

	t_start := time.Now()
	colors := cp.CF.F(histogram, cp.CFP)
	for pt, rgba := range colors {
		cp.Plane.SetXYColor(pt.X, pt.Y, rgba)
	}
	fmt.Printf("Image processing took %dms\n", time.Since(t_start).Milliseconds())
}
