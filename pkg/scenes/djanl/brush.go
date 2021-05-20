package djanl

import (
	"image"
	"image/draw"

	"github.com/HexHacks/goffer/pkg/image/mask"
)

type brush struct {
	img      image.Image
	mask     *mask.SoftCircle
	defMaskP image.Point
}

func newBrush(img image.Image) brush {
	cr := int(cutoutR)
	cp := image.Point{X: cr, Y: cr}
	return brush{
		img: img,
		//mask:     &mask.Circle{P: cp, R: cp.X},
		mask:     &mask.SoftCircle{P: cp, R: cr},
		defMaskP: cp,
	}
}

func (b *brush) SetR(r int) {
	cp := image.Point{X: r, Y: r}
	b.mask.P = cp
	b.mask.R = r
}

func (b *brush) Draw(onto draw.Image, dp image.Point) {
	var (
	//bnds = b.img.Bounds()
	//hx   = bnds.Size().X / 2
	//hy   = bnds.Size().Y / 2
	)
	ndp := image.Point{
		X: dp.X, //- hx,
		Y: dp.Y, //- hy,
	}
	drawFullSrcMask(onto, b.img, b.mask, ndp)
}
