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
	drawFullSrcMask(onto, b.img, &b.mask, dp)
}
