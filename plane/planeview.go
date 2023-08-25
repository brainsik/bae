package plane

import "fmt"

// PlaneView represents the bounds of complex plane.
// Min is the left-bottom point and max is the right-top point.
type PlaneView struct {
	Min, Max complex128
}

func (pv PlaneView) String() string {
	return fmt.Sprintf("%v ↗︎ %v", pv.Min, pv.Max)
}

// RealLen returns the real axis length.
func (pv PlaneView) RealLen() float64 {
	return real(pv.Max) - real(pv.Min)
}

// ImagLen returns the imaginary axis length.
func (pv PlaneView) ImagLen() float64 {
	return imag(pv.Max) - imag(pv.Min)
}
