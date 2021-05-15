package djanl

import (
	"image"
	"image/draw"

	"github.com/HexHacks/goffer/pkg/image/mask"
)

type brush struct {
	img      image.Image
	mask     mask.Circle
	defMaskP image.Point
}

func newBrush(img image.Image) brush {
	cp := cutoutR.Max.Div(2)
	return brush{
		img:      img,
		mask:     mask.Circle{P: cp, R: cp.X},
		defMaskP: cp,
	}
}

func (b *brush) Draw(onto draw.Image, dp image.Point) {
	var (
		bnds = b.img.Bounds()
		hx   = bnds.Size().X / 2
		hy   = bnds.Size().Y / 2
	)
	ndp := image.Point{
		X: dp.X - hx,
		Y: dp.Y - hy,
	}
	drawFullSrcMask(onto, b.img, &b.mask, ndp)
}
