package plane

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestWithInverted(t *testing.T) {
	p := NewPlane(complex(0, 0), complex(8, 4), 100)
	pi := p.WithInverted()

	if p != pi {
		t.Fatalf("Expected planes to be the same object, they are not: %p %p", p, pi)
	}
	if !pi.inverted {
		t.Errorf("Expected the plane to be inverted, it is not.")
	}
}

func TestNewOrigin(t *testing.T) {
	origin := complex(0, 0)
	p1 := NewPlane(origin, complex(8, 4), 100)

	expect_origin := origin + complex(1, -1)
	p2 := p1.NewOrigin(expect_origin)
	result_origin := p2.origin

	if p1 == p2 {
		t.Fatalf("Expected a new plane to be created.")
	}
	if result_origin != expect_origin {
		t.Errorf("Expected new size to be %v, got %v", expect_origin, result_origin)
	}
}

func TestNewSize(t *testing.T) {
	size := complex(8, 4)
	y_pixels := 100
	p1 := NewPlane(complex(0, 0), size, y_pixels)

	expect_size := size - complex(4, 0)
	aspect := real(expect_size) / imag(expect_size)
	expect_x_pixels := int(float64(y_pixels) * aspect)

	p2 := p1.NewSize(expect_size)
	result_size := p2.size
	result_x_pixels := p2.ImageWidth()

	if p1 == p2 {
		t.Fatalf("Expected a new plane to be created.")
	}
	if result_size != expect_size {
		t.Errorf("Expected new size to be %v, got %v", expect_size, result_size)
	}
	if result_x_pixels != expect_x_pixels {
		t.Errorf("Expected new image width to be %v, got %v", expect_x_pixels, result_x_pixels)
	}
}

func TestNewImageSize(t *testing.T) {
	size := complex(8, 4)
	y_pixels := 100
	p1 := NewPlane(complex(0, 0), size, y_pixels)

	p2 := p1.NewImageSize(y_pixels, y_pixels)
	result_size := p2.size
	result_x_pixels := p2.ImageWidth()
	result_y_pixels := p2.ImageHeight()

	expect_size := complex(4, 4)
	expect_x_pixels := y_pixels
	expect_y_pixels := y_pixels

	if p1 == p2 {
		t.Fatalf("Expected a new plane to be created.")
	}
	if result_size != expect_size {
		t.Errorf("Expected new size to be %v, got %v", expect_size, result_size)
	}
	if result_x_pixels != expect_x_pixels {
		t.Errorf("Expected new image width to be %v, got %v", expect_x_pixels, result_x_pixels)
	}
	if result_y_pixels != expect_y_pixels {
		t.Errorf("Expected new image height to be %v, got %v", expect_y_pixels, result_y_pixels)
	}
}

func TestView(t *testing.T) {
	p := NewPlane(complex(1, 1), complex(2, 2), 100)
	expect := PlaneView{Min: complex(0, 0), Max: complex(2, 2)}

	result := p.View()
	if &(p.view) == &result {
		t.Fatalf("Expected and result should be different objects.")
	}
	if result != expect {
		t.Errorf("Expected %v, got %v", expect, result)
	}
}

func TestImageHeight(t *testing.T) {
	y_pixels := 128
	p := NewPlane(complex(0, 0), complex(8, 4), y_pixels)

	expect := y_pixels
	result := p.ImageHeight()

	if result != expect {
		t.Error(result, expect)
	}
}

func TestImageWidth(t *testing.T) {
	y_pixels := 128
	p := NewPlane(complex(0, 0), complex(8, 4), y_pixels)

	aspect := real(p.size) / imag(p.size)
	expect := int(float64(y_pixels) * aspect)
	result := p.ImageWidth()

	if result != expect {
		t.Error(result, expect)
	}
}

func TestToComplexPoint(t *testing.T) {
	origin := complex(-1, 1)
	size := complex(1, 2)
	y_pixels := 128
	p := NewPlane(origin, size, y_pixels)
	p_inverted := NewPlane(origin, size, y_pixels).WithInverted()

	aspect := real(size) / imag(size)
	x_pixels := int(float64(y_pixels) * aspect)

	testCases := []struct {
		name   string
		plane  *Plane
		point  ImagePoint
		expect complex128
	}{
		{"xy-center", p, ImagePoint{x_pixels / 2, y_pixels / 2}, origin},
		{"xy-left-top", p, ImagePoint{0, 0}, complex(-1.5, 2)},
		{"xy-right-bottom", p, ImagePoint{x_pixels, y_pixels}, complex(-0.5, 0)},
		{"invert-xy-center", p_inverted, ImagePoint{x_pixels / 2, y_pixels / 2}, origin},
		{"invert-xy-left-top", p_inverted, ImagePoint{0, 0}, complex(-1.5, 0)},
		{"invert-xy-right-bottom", p_inverted, ImagePoint{x_pixels, y_pixels}, complex(-0.5, 2)},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) { //nolint:gosimple
			result := tc.plane.ToComplexPoint(tc.point)
			if result != tc.expect {
				t.Error(result, tc.expect)
			}
		})
	}
}

func TestToImagePoint(t *testing.T) {
	origin := complex(-1, 1)
	size := complex(2, 1)
	y_pixels := 100
	p := NewPlane(origin, size, y_pixels)
	p_inverted := NewPlane(origin, size, y_pixels).WithInverted()

	aspect := real(size) / imag(size)
	x_pixels := int(float64(y_pixels) * aspect)

	testCases := []struct {
		name   string
		plane  *Plane
		point  complex128
		expect ImagePoint
	}{
		{"z-center", p, origin, ImagePoint{x_pixels / 2, y_pixels / 2}},
		{"z-left-top", p, complex(-2, 1.5), ImagePoint{0, 0}},
		{"z-right-bottom", p, complex(0, 0.5), ImagePoint{x_pixels, y_pixels}},
		{"invert-z-center", p_inverted, origin, ImagePoint{x_pixels / 2, y_pixels / 2}},
		{"invert-z-left-top", p_inverted, complex(-2, 1.5), ImagePoint{0, y_pixels}},
		{"invert-z-right-bottom", p_inverted, complex(0, 0.5), ImagePoint{x_pixels, 0}},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) { //nolint:gosimple
			result := tc.plane.ToImagePoint(tc.point)
			if result != tc.expect {
				t.Error(result, tc.expect)
			}
		})
	}
}

func TestToImageToComplexToImagePoint(t *testing.T) {
	origin := complex(0, 0)
	// Use numbers that ensure points to pixel ratio does not divide cleanly
	size := complex(1.59, 1.59)
	x_pixels, y_pixels := 127, 127
	p := NewPlane(origin, size, x_pixels)

	testCases := []struct {
		name  string
		point ImagePoint
	}{
		{"xy-center", ImagePoint{x_pixels / 2, y_pixels / 2}},
		{"xy-left-top", ImagePoint{0, 0}},
		{"xy-right-bottom", ImagePoint{x_pixels, y_pixels}},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) { //nolint:gosimple
			z := p.ToComplexPoint(tc.point)
			result := p.ToImagePoint(z)
			if result != tc.point {
				t.Error(result, tc.point)
			}
		})
	}
}

func TestPlaneJSONMarshaler(t *testing.T) {
	origin := complex(2, 2)
	size := complex(8, 4)
	y_pixels := 64
	p := NewPlane(origin, size, y_pixels).WithInverted()

	result, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("json.Marshal Error: %v", err)
	}
	expect := `{"origin":[2,2],"size":[8,4],"view":[-2,0,6,4],"inverted":true,"image_size":[128,64]}`
	if string(result) != expect {
		t.Error(string(result), expect)
	}
}

func TestPlaneJSONUnmarshaler(t *testing.T) {
	data := []byte(`{"origin":[2,2],"size":[8,4],"view":[-2,0,6,4],"inverted":true, "image_size":[128,64]}`)

	var result Plane
	err := json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("json.Unmarshal Error: %v", err)
	}
	result_width := result.ImageWidth()

	expect := NewPlane(complex(2, 2), complex(8, 4), 64).WithInverted()
	expect_width := expect.ImageWidth()

	if result.origin != expect.origin {
		t.Error(result.origin, expect.origin)
	}
	if result.size != expect.size {
		t.Error(result.size, expect.size)
	}
	if result_width != expect_width {
		t.Error(result_width, expect_width)
	}
	if result_width != expect_width {
		t.Error(result.inverted, expect.inverted)
	}
}
