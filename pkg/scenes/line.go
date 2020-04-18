package scenes

import (
	"image"
	"math"

	"github.com/HexHacks/goffer/pkg/global"
	jimage "github.com/HexHacks/goffer/pkg/image"
	jmath "github.com/HexHacks/goffer/pkg/math"
	"github.com/HexHacks/goffer/pkg/math/float"
	"github.com/HexHacks/goffer/pkg/palette"
	"github.com/llgcode/draw2d"
	//	"github.com/llgcode/draw2d"
)

type Lines struct {
	I draw2d.Matrix
}

func (l *Lines) Init() {
	palette.Palette = palette.Palette1
	l.I = draw2d.NewIdentityMatrix()
}

func (l *Lines) Frame(t float64) image.Image {
	const (
		res   = 1000.0
		lines = 20.0
		cY    = -15
	)

	var (
		w, h = global.W, global.H
		dY   = h / (lines - 1)
		dX   = w / res
		dXzo = dX / (w - 1)
	)

	img, gc := jimage.New()
	//gc.SetMatrixTransform(draw2d.NewScaleMatrix(w/2.0, h/2.0))
	//gc.ComposeMatrixTransform(draw2d.NewTranslationMatrix(-w/2.0, h/2.0))

	fun := func(x, y float64) float64 {
		saw := math.Mod(x, 0.3)
		marker := 20 + saw*20 - float.Gaussian(y, 1.0, 0.5, 5.0/(t+.01))*20. - t
		return 10.0 * math.Sin(jmath.Tau*t*0.3+y) * float.Morlet(6.0, marker)
	}

	gc.SetLineWidth(1.5)
	for y := dY; y < h; y += dY {
		gc.MoveTo(0.0, cY+y-fun(-1.0, y/(h-1.0)))
		for x := dX; x < w; x += dX {
			xzo := 2*x/(w-1) - 1

			//hx, hy := x, y+fun(xzo)
			nx, ny := x+dX, cY+y-fun(xzo+dXzo, y/(h-1.0))

			gc.LineTo(nx, ny)
		}
		gc.Stroke()
	}

	return img
}
