package mask

import (
	"image"
	"image/color"

	"github.com/bdazl/goffer/pkg/math/float"
)

type SoftCircle struct {
	P image.Point
	R int
}

func (c *SoftCircle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *SoftCircle) Bounds() image.Rectangle {
	return image.Rect(c.P.X-c.R, c.P.Y-c.R, c.P.X+c.R, c.P.Y+c.R)
}

func (c *SoftCircle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.P.X)+0.5, float64(y-c.P.Y)+0.5, float64(c.R)
	dsq := xx*xx + yy*yy
	rsq := rr * rr
	frsq := float64(rsq)
	if dsq < rsq {
		m := rsq - dsq
		to1 := float64(m) / frsq
		log := float.Logistic(to1, 0.22, 1, 34)

		v := float.Clamp(255.0*log, 0.0, 255.0)
		return color.Alpha{uint8(v)}
	}
	return color.Alpha{0}
}
