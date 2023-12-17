package scenes

import (
	"image"
	"math"

	"github.com/bdazl/goffer/pkg/global"
	jimage "github.com/bdazl/goffer/pkg/image"
	jmath "github.com/bdazl/goffer/pkg/math"
	"github.com/bdazl/goffer/pkg/math/float"
	"github.com/bdazl/goffer/pkg/palette"
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
		lines = 30.0
		cY    = -15
	)

	var (
		w, h = global.W, global.H
		tot  = global.Total
		dY   = h / (lines - 1)
		dX   = w / res
		dXzo = dX / (w - 1)
	)

	img, gc := jimage.New()
	//gc.SetMatrixTransform(draw2d.NewScaleMatrix(w/2.0, h/2.0))
	//gc.ComposeMatrixTransform(draw2d.NewTranslationMatrix(-w/2.0, h/2.0))

	fun := func(x, y float64) float64 {
		//saw := math.Mod(x, 0.25) - 1.0
		marker := 10 + 20.*x - 10.*float.Morlet(4, 3*(y*2-1)) - 20*float.Smoothstep(y, tot, t)
		wave := 10.0 * math.Sin(jmath.Tau*t*0.2+y*2.0) * float.Morlet(6.0, marker)
		return wave
	}

	yy := 0
	gc.SetLineWidth(2.0)
	for y := dY; y < h; y += dY {
		gc.SetStrokeColor(palette.Palette[1+yy%2])
		gc.MoveTo(0.0, cY+y-fun(-1.0, y/(h-1.0)))
		for x := dX; x < w; x += dX {
			xzo := 2*x/(w-1) - 1

			//hx, hy := x, y+fun(xzo)
			nx, ny := x+dX, cY+y-fun(xzo+dXzo, y/(h-1.0))

			gc.LineTo(nx, ny)
		}
		gc.Stroke()
		yy++
	}

	return img
}
