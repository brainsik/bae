package plane

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

// Plane represents an area in the complex plane and it's corresponding graphical image.
type Plane struct {
	Origin, Size complex128
	View         PlaneView
	Inverted     bool

	r_step, i_step float64
	x_step, y_step float64

	Image *image.NRGBA
}

// PlaneView represents a rectangle in the complex plane. Min is the left-bottom point and max is the right-top point.
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

// NewPlane returns a new Plane.
func NewPlane(origin, size complex128, x_pixels int) *Plane {
	view := PlaneView{origin - size/2, origin + size/2}
	_aspect_ratio := imag(size) / real(size)
	y_pixels := int(float64(x_pixels) * _aspect_ratio)

	// Points per pixel.
	r_step := real(size) / float64(x_pixels)
	i_step := imag(size) / float64(y_pixels)

	// Pixels per point.
	x_step := float64(x_pixels) / real(size)
	y_step := float64(y_pixels) / imag(size)

	img := image.NewNRGBA(
		image.Rectangle{image.Point{0, 0}, image.Point{x_pixels, y_pixels}})
	draw.Draw(img, img.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)

	return &Plane{
		Origin:   origin,
		Size:     size,
		View:     view,
		Inverted: false,

		// Points per pixel.
		r_step: r_step,
		i_step: i_step,

		// Pixels per point.
		x_step: x_step,
		y_step: y_step,

		Image: img,
	}
}

func NewPlaneInverted(origin, size complex128, x_pixels int) *Plane {
	p := NewPlane(origin, size, x_pixels)
	p.Inverted = true
	return p
}

func (p *Plane) String() string {
	return fmt.Sprintf(
		"Plane{Origin:%v, View:%v, Image:%dx%d}",
		p.Origin, p.View, p.ImageWidth(), p.ImageHeight())
}

// PlanePoint returns the point in the complex plane corresponding to the given point in the image plane.
func (p *Plane) PlanePoint(px ImagePoint) complex128 {
	if px.X < 0 || px.X > p.ImageWidth() {
		fmt.Printf("Warning: PlanePoint(%v) x coordinate is outside image bounds: 0 -> %v\n", px, p.ImageWidth())
	}
	if px.Y < 0 || px.Y > p.ImageHeight() {
		fmt.Printf("Warning: PlanePoint(%v) y coordinate is outside image bounds: 0 -> %v\n", px, p.ImageHeight())
	}

	r := real(p.View.Min) + float64(px.X)*p.r_step

	var i float64
	if !p.Inverted {
		// i on the complex plane and y on the pixel plane increase in opposite directions
		i = imag(p.View.Max) - float64(px.Y)*p.i_step
	} else {
		// leave inverted
		i = imag(p.View.Min) + float64(px.Y)*p.i_step
	}
	return complex(r, i)
}

// ImagePoint represents coordinates in the image plane.
type ImagePoint struct {
	X, Y int
}

func (ip ImagePoint) String() string {
	return fmt.Sprintf("(%v, %v)", ip.X, ip.Y)
}

// PlanePoint returns the point in the image plane corresponding to the given point in the complex plane.
func (p *Plane) ImagePoint(z complex128) ImagePoint {
	// rz_min, rz_max := real(p.view.min), real(p.view.max)
	// if real(z) < rz_min || real(z) > rz_max {
	// 	fmt.Printf("Warning: ImagePoint(%v) real value is outside plane bounds: %v -> %v\n", z, rz_min, rz_max)
	// }

	// iz_min, iz_max := imag(p.view.min), imag(p.view.max)
	// if imag(z) < iz_min || imag(z) > iz_max {
	// 	fmt.Printf("Warning: ImagePoint(%v) imag value is outside plane bounds: %v -> %v\n", z, iz_min, iz_max)
	// }

	// reorient view so min is complex(0,0)
	z_adj := z - p.View.Min

	x := int(math.Round(real(z_adj) * p.x_step))
	y := int(math.Round(imag(z_adj) * p.y_step))
	if !p.Inverted {
		// Flip y, it increases in the opposite direction as i.
		y = p.ImageHeight() - y
	}

	return ImagePoint{x, y}
}

// Set sets the color in the image plane corresponding to the given complex plane point.
func (p *Plane) Set(z complex128, rgba color.NRGBA) {
	xy := p.ImagePoint(z)
	p.Image.Set(xy.X, xy.Y, rgba)
}

// ImageWidth returns the image width.
func (p *Plane) ImageWidth() int {
	return p.Image.Rect.Dx()
}

// ImageHeight returns the image width.
func (p *Plane) ImageHeight() int {
	return p.Image.Rect.Dy()
}

// WritePNG outputs a PNG file at the given path.
func (p *Plane) WritePNG(path string) {
	png_file, _ := os.Create(path)
	penc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := penc.Encode(png_file, p.Image); err != nil {
		fmt.Printf("Error encoding PNG: %v\n", err)
	}
	fmt.Printf("Wrote %s\n", png_file.Name())
	png_file.Close()
}

// planeJSON is the JSON representation of a Plane.
type planeJSON struct {
	Origin    [2]float64 `json:"origin"`
	Size      [2]float64 `json:"size"`
	View      [4]float64 `json:"view,omitempty"`
	Inverted  bool       `json:"inverted"`
	ImageSize [2]int     `json:"image_size"`
}

func (p *Plane) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		planeJSON{
			Origin:    [2]float64{real(p.Origin), imag(p.Origin)},
			Size:      [2]float64{real(p.Size), imag(p.Size)},
			View:      [4]float64{real(p.View.Min), imag(p.View.Min), real(p.View.Max), imag(p.View.Max)},
			Inverted:  p.Inverted,
			ImageSize: [2]int{p.ImageWidth(), p.ImageHeight()},
		})
}

func (p *Plane) UnmarshalJSON(data []byte) error {
	var v planeJSON
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	origin := complex(v.Origin[0], v.Origin[1])
	size := complex(v.Size[0], v.Size[1])
	x_pixels := v.ImageSize[0]

	if v.Inverted {
		*p = *NewPlaneInverted(origin, size, x_pixels)
	} else {
		*p = *NewPlane(origin, size, x_pixels)
	}

	return nil
}
