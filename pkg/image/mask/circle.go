package draw

import (
	"image"
	"image/color"
)

type CircleMask struct {
	p image.Point
	r int
}

func (c *CircleMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *CircleMask) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *CircleMask) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
