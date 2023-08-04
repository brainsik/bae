package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestImageSize(t *testing.T) {
	x_pixels := 128
	p := NewPlane(complex(0, 0), complex(8, 4), x_pixels)

	expected := ImageSize{x_pixels, x_pixels / 2}
	result := p.ImageSize()

	if result != expected {
		t.Error(result, expected)
	}
}

func TestPlanePoint(t *testing.T) {
	origin := complex(0, 0)
	size := complex(8, 8)
	x_pixels := 128
	p := NewPlane(origin, size, x_pixels)

	testCases := []struct {
		name     string
		point    ImagePoint
		expected complex128
	}{
		{"center", ImagePoint{x_pixels / 2, x_pixels / 2}, complex(0, 0)},
		{"left-top", ImagePoint{0, x_pixels}, complex(-4, 4)},
		{"right-bottom", ImagePoint{x_pixels, 0}, complex(4, -4)},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) {
			result := p.PlanePoint(tc.point)
			if result != tc.expected {
				t.Error(result, tc.expected)
			}
		})
	}
}

func TestImagePoint(t *testing.T) {
	origin := complex(0, 0)
	size := complex(8, 8)
	x_pixels := 128
	p := NewPlane(origin, size, x_pixels)

	testCases := []struct {
		name     string
		point    complex128
		expected ImagePoint
	}{
		{"center", complex(0, 0), ImagePoint{x_pixels / 2, x_pixels / 2}},
		{"left-top", complex(-4, -4), ImagePoint{0, x_pixels}},
		{"right-bottom", complex(4, 4), ImagePoint{x_pixels, 0}},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) {
			result := p.ImagePoint(tc.point)
			if result != tc.expected {
				t.Error(result, tc.expected)
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
	expected := `{"origin":[2,2],"size":[8,4],"view":[-2,4,6,0],"image_size":[128,64]}`
	if string(result) != expected {
		t.Error(string(result), expected)
	}
}

func TestPlaneJSONUnmarshaler(t *testing.T) {
	data := []byte(`{"origin":[2,2],"size":[8,4],"view":[-2,4,6,0],"image_size":[128,64]}`)

	var result Plane
	err := json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("json.Unmarshal Error: %v", err)
	}
	result_width := result.ImageSize().width

	expected := NewPlane(complex(2, 2), complex(8, 4), 128)
	expected_width := expected.ImageSize().width

	if result.origin != expected.origin {
		t.Error(result.origin, expected.origin)
	}
	if result.size != expected.size {
		t.Error(result.size, expected.size)
	}
	if result_width != expected_width {
		t.Error(result_width, expected_width)
	}
}
