package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestImageSize(t *testing.T) {
	x_pixels := 128
	p := NewPlane(complex(0, 0), complex(8, 4), x_pixels)

	expect := ImageSize{x_pixels, x_pixels / 2}
	result := p.ImageSize()

	if result != expect {
		t.Error(result, expect)
	}
}

func TestPlanePoint(t *testing.T) {
	origin := complex(-1, 1)
	size := complex(1, 2)
	x_pixels := 128
	p := NewPlane(origin, size, x_pixels)

	aspect := imag(size) / real(size)
	y_pixels := int(float64(x_pixels) * aspect)

	testCases := []struct {
		name   string
		point  ImagePoint
		expect complex128
	}{
		{"xy-center", ImagePoint{x_pixels / 2, y_pixels / 2}, origin},
		{"xy-left-top", ImagePoint{0, 0}, complex(-1.5, 2)},
		{"xy-right-bottom", ImagePoint{x_pixels, y_pixels}, complex(-0.5, 0)},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) { //nolint:gosimple
			result := p.PlanePoint(tc.point)
			if result != tc.expect {
				t.Error(result, tc.expect)
			}
		})
	}
}

func TestImagePoint(t *testing.T) {
	origin := complex(-1, 1)
	size := complex(2, 1)
	x_pixels := 100
	p := NewPlane(origin, size, x_pixels)

	aspect := imag(size) / real(size)
	y_pixels := int(float64(x_pixels) * aspect)

	testCases := []struct {
		name   string
		point  complex128
		expect ImagePoint
	}{
		{"z-center", origin, ImagePoint{x_pixels / 2, y_pixels / 2}},
		{"z-left-top", complex(-2, 1.5), ImagePoint{0, 0}},
		{"z-right-bottom", complex(0, 0.5), ImagePoint{x_pixels, y_pixels}},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) { //nolint:gosimple
			result := p.ImagePoint(tc.point)
			if result != tc.expect {
				t.Error(result, tc.expect)
			}
		})
	}
}

func TestImageToPlaneToImagePoint(t *testing.T) {
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
			z := p.PlanePoint(tc.point)
			result := p.ImagePoint(z)
			if result != tc.point {
				t.Error(result, tc.point)
			}
		})
	}
}

func TestPlaneJSONMarshaler(t *testing.T) {
	origin := complex(2, 2)
	size := complex(8, 4)
	x_pixels := 128
	p := NewPlane(origin, size, x_pixels)

	result, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("json.Marshal Error: %v", err)
	}
	expect := `{"origin":[2,2],"size":[8,4],"view":[-2,0,6,4],"image_size":[128,64]}`
	if string(result) != expect {
		t.Error(string(result), expect)
	}
}

func TestPlaneJSONUnmarshaler(t *testing.T) {
	data := []byte(`{"origin":[2,2],"size":[8,4],"view":[-2,0,6,4],"image_size":[128,64]}`)

	var result Plane
	err := json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("json.Unmarshal Error: %v", err)
	}
	result_width := result.ImageSize().width

	expect := NewPlane(complex(2, 2), complex(8, 4), 128)
	expect_width := expect.ImageSize().width

	if result.origin != expect.origin {
		t.Error(result.origin, expect.origin)
	}
	if result.size != expect.size {
		t.Error(result.size, expect.size)
	}
	if result_width != expect_width {
		t.Error(result_width, expect_width)
	}
}
