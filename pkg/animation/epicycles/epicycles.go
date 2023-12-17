package epicycles

import (
	"image"
	"image/color"
	"image/draw"
	"math/cmplx"
	"sort"

	"github.com/bdazl/goffer/pkg/animation"
	"github.com/bdazl/goffer/pkg/coordsys"
	jimage "github.com/bdazl/goffer/pkg/image"
	"github.com/bdazl/goffer/pkg/math/fourier"
	"github.com/bdazl/goffer/pkg/palette"

	"github.com/llgcode/draw2d/draw2dimg"
	kit "github.com/llgcode/draw2d/draw2dkit"
)

type Epicycles struct {
	C      []complex128 // fourier coefficients
	sorted Coeffs
	pts    []complex128

	anim    *animation.Animation
	currImg image.Image

	palette color.Palette
}

type CoefSort struct {
	I int        // index
	C complex128 // coefficient
}

type Coeffs []CoefSort

// Sort by coef magnitude, largest first
func (a Coeffs) Len() int           { return len(a) }
func (a Coeffs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Coeffs) Less(i, j int) bool { return cmplx.Abs(a[i].C) > cmplx.Abs(a[j].C) }

func New(coef []complex128, ptCount int) *Epicycles {
	epi := &Epicycles{
		C: coef,
	}

	epi.sorted = MakeSorted(coef)

	epi.pts = make([]complex128, ptCount)
	lcf := float64(ptCount)
	for i := 0; i < ptCount; i++ {
		fi := float64(i)
		epi.pts[i] = fourier.P(fi/lcf, epi.C)
	}

	common := func() *draw2dimg.GraphicContext {
		img, gc := jimage.New()
		draw.Draw(img,
			img.Bounds(),
			&image.Uniform{palette.Palette[0]},
			image.ZP, draw.Src)

		epi.currImg = img
		return gc
	}

	anim0 := func(t float64) {
		gc := common()
		cnt := int(t * float64(len(epi.sorted)))
		epi.DrawEpiCircles(gc, 0, cnt)
	}

	anim1 := func(t float64) {
		gc := common()
		epi.DrawEpiCircles(gc, t, len(epi.sorted))

		cnt := int(t * float64(len(epi.pts)))
		DrawCmplxLines(gc, epi.pts, cnt)
	}

	anim2 := func(t float64) {
		gc := common()

		if t < 0.5 {
			cnt := int((1 - t*2) * float64(len(epi.sorted)))
			epi.DrawEpiCircles(gc, 1.0, cnt)
		}

		DrawCmplxLines(gc, epi.pts, len(epi.pts))
	}

	epi.anim = animation.New([]animation.Animator{
		anim0, anim1, anim2,
	},
		[]float64{0.15, 0.5, 0.35},
	)

	return epi
}

func (e *Epicycles) Frame(t float64) image.Image {
	e.anim.Frame(t)
	return e.currImg
}

func MakeSorted(C []complex128) Coeffs {
	sorted := make(Coeffs, len(C))
	for i, c := range C {
		sorted[i].I = i
		sorted[i].C = c
	}

	sort.Sort(sorted)
	return sorted
}

func DrawCmplxLines(gc *draw2dimg.GraphicContext, pts []complex128, count int) {
	if count > len(pts) {
		count = len(pts)
	}

	gc.SetStrokeColor(palette.Palette[3])
	// draw line through all points
	start := coordsys.UnitToImgC(pts[0])
	p := start
	gc.MoveTo(real(p), imag(p))
	for i := 0; i < count; i++ {
		p := coordsys.UnitToImgC(pts[i])
		gc.LineTo(real(p), imag(p))
	}

	if count == len(pts) {
		gc.LineTo(real(start), imag(start))
	}

	gc.Stroke()
}

func (se Epicycles) DrawEpiCircles(gc *draw2dimg.GraphicContext, t float64, count int) {

	if count >= len(se.sorted) {
		count = len(se.sorted) - 1
	}

	h := len(se.C) / 2
	center := complex(0, 0)
	for i := 0; i < count; i++ {
		srt := se.sorted[i]
		p := fourier.Pat(t, se.C, srt.I)

		// don't draw the static one
		if srt.I != h {
			DrawCCirc(gc, center, p)
		}

		center += p
	}

	center = complex(0, 0)
	for i := 0; i < count; i++ {
		srt := se.sorted[i]
		p := fourier.Pat(t, se.C, srt.I)

		if srt.I != h {
			DrawCLine(gc, center, p)
		}

		center += p
	}

}

func DrawCCirc(gc *draw2dimg.GraphicContext, center, coef complex128) {
	cent := coordsys.UnitToImgC(complex(0, 0))
	t, c := coordsys.UnitToImgC(center), coordsys.UnitToImgC(coef)
	tc := c - cent

	gc.SetStrokeColor(palette.Palette[1])
	kit.Circle(gc, real(t), imag(t), cmplx.Abs(tc))
	gc.Stroke()
}

func DrawCLine(gc *draw2dimg.GraphicContext, center, coef complex128) {
	cent := coordsys.UnitToImgC(complex(0, 0))
	t, c := coordsys.UnitToImgC(center), coordsys.UnitToImgC(coef)
	tc := c - cent

	gc.SetStrokeColor(palette.Palette[2])
	gc.MoveTo(real(t), imag(t))
	gc.LineTo(real(t+tc), imag(t+tc))
	gc.Stroke()
}
