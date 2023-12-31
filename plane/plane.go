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
	origin, size complex128
	view         PlaneView
	inverted     bool

	r_step, i_step float64
	x_step, y_step float64

	image *image.NRGBA
}

// NewPlane returns a new Plane.
func NewPlane(origin, size complex128, y_pixels int) *Plane {
	view := PlaneView{origin - size/2, origin + size/2}
	aspect_ratio := real(size) / imag(size)
	x_pixels := int(float64(y_pixels) * aspect_ratio)

	// Points per pixel.
	r_step := real(size) / float64(x_pixels)
	i_step := imag(size) / float64(y_pixels)

	// Pixels per point.
	x_step := float64(x_pixels) / real(size)
	y_step := float64(y_pixels) / imag(size)

	img := image.NewNRGBA(
		image.Rectangle{image.Point{0, 0}, image.Point{x_pixels, y_pixels}})

	p := Plane{
		origin:   origin,
		size:     size,
		view:     view,
		inverted: false,

		// Points per pixel.
		r_step: r_step,
		i_step: i_step,

		// Pixels per point.
		x_step: x_step,
		y_step: y_step,

		image: img,
	}
	p.ResetImage()
	return &p
}

// WithInverted returns the same Plane with Inverted true.
func (p *Plane) WithInverted() *Plane {
	p.inverted = true
	return p
}

// NewOrigin returns a new Plane with the given origin.
func (p *Plane) NewOrigin(origin complex128) *Plane {
	new := NewPlane(origin, p.size, p.ImageHeight())
	new.inverted = p.inverted
	return new
}

// NewSize returns a new Plane with the given complex plane size.
func (p *Plane) NewSize(size complex128) *Plane {
	new := NewPlane(p.origin, size, p.ImageHeight())
	new.inverted = p.inverted
	return new
}

// NewImageSize returns a new Plane with the given image size.
func (p *Plane) NewImageSize(width, height int) *Plane {
	aspect := float64(width) / float64(height)

	size := complex(imag(p.size)*aspect, imag(p.size)) // keep i size
	if p.ImageWidth() == width {
		size = complex(real(p.size), real(p.size)*(1.0/aspect)) // keep r size
	}

	new := NewPlane(p.origin, size, height)
	new.inverted = p.inverted
	return new
}

func (p *Plane) String() string {
	return fmt.Sprintf(
		"Plane{Origin:%v, View:%v, Image:%dx%d}",
		p.origin, p.view, p.ImageWidth(), p.ImageHeight())
}

// ResetImage resets the image to all black.
func (p *Plane) ResetImage() {
	draw.Draw(p.image, p.image.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)
}

// Image returns the image buffer.
func (p *Plane) Image() *image.NRGBA {
	return p.image
}

// Origin returns the complex plane origin.
func (p *Plane) Origin() complex128 {
	return p.origin
}

// Size returns the complex plane size.
func (p *Plane) Size() complex128 {
	return p.size
}

// View returns the complex plane bounds.
func (p *Plane) View() PlaneView {
	return p.view
}

// Aspect returns the aspect ratio (width/height) of the planes.
func (p *Plane) Aspect() float64 {
	return real(p.size) / imag(p.size)
}

// IsInverted returns whether the complex plane is inverted.
func (p *Plane) IsInverted() bool {
	return p.inverted
}

// ToComplexPoint returns the point in the complex plane corresponding to the given point in the image plane.
func (p *Plane) ToComplexPoint(px ImagePoint) complex128 {
	if px.X < 0 || px.X > p.ImageWidth() {
		fmt.Printf("Warning: ToComplexPoint(%v) x coordinate is outside image bounds: 0 -> %v\n", px, p.ImageWidth())
	}
	if px.Y < 0 || px.Y > p.ImageHeight() {
		fmt.Printf("Warning: ToComplexPoint(%v) y coordinate is outside image bounds: 0 -> %v\n", px, p.ImageHeight())
	}

	r := real(p.view.Min) + float64(px.X)*p.r_step

	var i float64
	if !p.inverted {
		// i on the complex plane and y on the pixel plane increase in opposite directions
		i = imag(p.view.Max) - float64(px.Y)*p.i_step
	} else {
		// leave inverted
		i = imag(p.view.Min) + float64(px.Y)*p.i_step
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
func (p *Plane) ToImagePoint(z complex128) ImagePoint {
	// rz_min, rz_max := real(p.view.min), real(p.view.max)
	// if real(z) < rz_min || real(z) > rz_max {
	// 	fmt.Printf("Warning: ImagePoint(%v) real value is outside plane bounds: %v -> %v\n", z, rz_min, rz_max)
	// }

	// iz_min, iz_max := imag(p.view.min), imag(p.view.max)
	// if imag(z) < iz_min || imag(z) > iz_max {
	// 	fmt.Printf("Warning: ImagePoint(%v) imag value is outside plane bounds: %v -> %v\n", z, iz_min, iz_max)
	// }

	// reorient view so min is complex(0,0)
	z_adj := z - p.view.Min

	x := int(math.Round(real(z_adj) * p.x_step))
	y := int(math.Round(imag(z_adj) * p.y_step))
	if !p.inverted {
		// Flip y, it increases in the opposite direction as i.
		y = p.ImageHeight() - y
	}

	return ImagePoint{x, y}
}

// SetZColor sets the color in the image plane corresponding to the given plane point.
func (p *Plane) SetZColor(z complex128, rgba color.NRGBA) {
	xy := p.ToImagePoint(z)
	p.image.Set(xy.X, xy.Y, rgba)
}

// SetXYColor sets the color in the image plane corresponding to the given image point.
func (p *Plane) SetXYColor(x, y int, rgba color.NRGBA) {
	p.image.Set(x, y, rgba)
}

// ImageWidth returns the image width.
func (p *Plane) ImageWidth() int {
	return p.image.Rect.Dx()
}

// ImageHeight returns the image width.
func (p *Plane) ImageHeight() int {
	return p.image.Rect.Dy()
}

// WritePNG outputs a PNG file at the given path.
func (p *Plane) WritePNG(path string) {
	png_file, _ := os.Create(path)
	penc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := penc.Encode(png_file, p.image); err != nil {
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
			Origin:    [2]float64{real(p.origin), imag(p.origin)},
			Size:      [2]float64{real(p.size), imag(p.size)},
			View:      [4]float64{real(p.view.Min), imag(p.view.Min), real(p.view.Max), imag(p.view.Max)},
			Inverted:  p.inverted,
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
	y_pixels := v.ImageSize[1]

	*p = *NewPlane(origin, size, y_pixels)
	p.inverted = v.Inverted

	return nil
}
