package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

type Plane struct {
	origin, size   complex128
	view           PlaneView
	r_step, i_step float64
	x_step, y_step int
	image          *image.NRGBA
}

type PlaneView struct {
	min, max complex128
}

func (view PlaneView) String() string {
	return fmt.Sprintf("%v, %v", view.min, view.max)
}

func NewPlane(origin, size complex128, x_pixels int) *Plane {
	// orient like we think of a pixel image: start -> top-left, end -> bottom-right
	view := PlaneView{
		complex(real(origin)-real(size)/2.0, imag(origin)+imag(size)/2.0),
		complex(real(origin)+real(size)/2.0, imag(origin)-imag(size)/2.0)}

	_aspect_ratio := imag(size) / real(size)
	y_pixels := int(float64(x_pixels) * _aspect_ratio)

	// Points per pixel.
	r_step := real(size) / float64(x_pixels)
	i_step := imag(size) / float64(y_pixels)

	// Pixels per point.
	x_step := int(float64(x_pixels) / real(size))
	y_step := int(float64(y_pixels) / imag(size))

	img := image.NewNRGBA(
		image.Rectangle{image.Point{0, 0}, image.Point{x_pixels, y_pixels}})
	draw.Draw(img, img.Bounds(), image.NewUniform(color.Black), image.ZP, draw.Src)

	return &Plane{
		origin: origin,
		size:   size,
		view:   view,

		// Points per pixel.
		r_step: r_step,
		i_step: i_step,

		// Pixels per point.
		x_step: x_step,
		y_step: y_step,

		image: img,
	}
}

func (p *Plane) String() string {
	return fmt.Sprintf(
		"Plane{\nOrigin: %v\nView:   %v\nSteps:  %v, %vi\nImage:  %vx%v\n}",
		p.origin, p.view, p.r_step, p.i_step, p.ImageSize().width, p.ImageSize().height)
}

func (p *Plane) PlanePoint(px ImagePoint) complex128 {
	r := real(p.view.min) + float64(px.x)*p.r_step
	// i on the complex plane and y on the pixel plane increase in opposite directions
	i := -imag(p.view.min) + float64(px.y)*p.i_step
	return complex(r, i)
}

type ImagePoint struct {
	x, y int
}

func (p *Plane) ImagePoint(z complex128) ImagePoint {
	x := int((real(z) - real(p.view.min)) * float64(p.x_step))
	// y on the pixel plane and i on the complex plane increase in opposite directions
	y := -int((imag(z) - imag(p.view.min)) * float64(p.y_step))
	return ImagePoint{x, y}
}

func (p *Plane) Set(z complex128, rgba color.NRGBA) {
	xy := p.ImagePoint(z)
	p.image.Set(xy.x, xy.y, rgba)
}

type ImageSize struct {
	width, height int
}

func (p *Plane) ImageSize() ImageSize {
	return ImageSize{p.image.Rect.Dx(), p.image.Rect.Dy()}
}

type planeJSON struct {
	Origin    [2]float64 `json:"origin"`
	Size      [2]float64 `json:"size"`
	View      [4]float64 `json:"view,omitempty"`
	ImageSize [2]int     `json:"image_size"`
}

func (p *Plane) MarshalJSON() ([]byte, error) {
	img_size := p.ImageSize()
	return json.Marshal(
		planeJSON{
			Origin:    [2]float64{real(p.origin), imag(p.origin)},
			Size:      [2]float64{real(p.size), imag(p.size)},
			View:      [4]float64{real(p.view.min), imag(p.view.min), real(p.view.max), imag(p.view.max)},
			ImageSize: [2]int{img_size.width, img_size.height},
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
	*p = *NewPlane(origin, size, x_pixels)

	return nil
}
