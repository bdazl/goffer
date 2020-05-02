package scenes

import (
	"fmt"
	"image"
	"image/draw"
	"math/cmplx"
	"os"
	"sort"

	"github.com/HexHacks/goffer/pkg/animation"
	"github.com/HexHacks/goffer/pkg/coordsys"
	jimage "github.com/HexHacks/goffer/pkg/image"
	"github.com/HexHacks/goffer/pkg/math/fourier"
	"github.com/HexHacks/goffer/pkg/palette"
	"github.com/HexHacks/goffer/pkg/svg"

	"github.com/llgcode/draw2d/draw2dimg"
	kit "github.com/llgcode/draw2d/draw2dkit"
	"gonum.org/v1/gonum/spatial/r2"
)

type SvgEpicycle struct {
	F          int
	operations []svg.Operation
	origPts    []r2.Vec
	svgPts     []complex128
	shift      []complex128
	coeff      []complex128
	sorted     coeffs

	pts []complex128

	anim    *animation.Animation
	currImg image.Image
}

type coefSort struct {
	idx  int
	coef complex128
}

type coeffs []coefSort

// Sort by coef magnitude, largest first
func (a coeffs) Len() int           { return len(a) }
func (a coeffs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a coeffs) Less(i, j int) bool { return cmplx.Abs(a[i].coef) > cmplx.Abs(a[j].coef) }

func (se *SvgEpicycle) Init() {
	const (
		perBezier = 100
		imgScale  = 3.8
		coefOrder = 50
	)

	fil, err := os.Open("assets/gubb_abs.svg")
	panicOn(err)
	defer fil.Close()

	s, err := svg.ParseSvg(fil)
	panicOn(err)

	se.operations = s.Groups[0].Paths[0].Operations
	se.origPts = ExtractPoints(s)
	se.svgPts = ExpandCurve(se.origPts, perBezier, imgScale)
	se.shift = ShiftPoints(se.svgPts)
	se.coeff = fourier.Coefficients(se.shift, coefOrder)

	se.sorted = make(coeffs, len(se.coeff))
	for i, c := range se.coeff {
		se.sorted[i].idx = i
		se.sorted[i].coef = c
	}

	sort.Sort(se.sorted)

	loopCount := len(se.svgPts)
	se.pts = make([]complex128, loopCount)
	lcf := float64(loopCount)
	for i := 0; i < loopCount; i++ {
		fi := float64(i)
		se.pts[i] = fourier.P(fi/lcf, se.coeff)
	}

	common := func() *draw2dimg.GraphicContext {
		img, gc := jimage.New()
		draw.Draw(img,
			img.Bounds(),
			&image.Uniform{palette.Palette[0]},
			image.ZP, draw.Src)

		se.currImg = img
		return gc
	}

	anim0 := func(t float64) {
		gc := common()
		cnt := int(t * float64(len(se.sorted)))
		se.DrawEpiCircles(gc, 0, cnt)
	}

	anim1 := func(t float64) {
		gc := common()
		se.DrawEpiCircles(gc, t, len(se.sorted))

		cnt := int(t * float64(len(se.pts)))
		DrawCmplxLines(gc, se.pts, cnt)
	}

	anim2 := func(t float64) {
		gc := common()

		if t < 0.5 {
			cnt := int((1 - t*2) * float64(len(se.sorted)))
			se.DrawEpiCircles(gc, 1.0, cnt)
		}

		DrawCmplxLines(gc, se.pts, len(se.pts))
	}

	se.anim = animation.New([]animation.Animator{
		anim0, anim1, anim2,
	},
		[]float64{0.15, 0.5, 0.35},
	)
}

func ExpandCurve(curve []r2.Vec, perBezier int, scale float64) []complex128 {
	pB := float64(perBezier)

	l := len(curve)
	out := make([]complex128, 0, l*perBezier)
	for i := 0; i < l; i = i + 3 {
		if l <= i+3 {
			fmt.Printf("remaining points can't be expanded, i: %v, len: %v\n", i, l)
			break
		}
		a, b, c, d := curve[i], curve[i+1], curve[i+2], curve[i+3]

		for j := 0.0; j < pB; j = j + 1.0 {
			p := cBezier(j/pB, a, b, c, d)
			out = append(out, complex(p.X*scale, p.Y*scale))
		}
	}
	return out
}

func ExtractPoints(s *svg.Svg) []r2.Vec {
	curve := make([]r2.Vec, 0)
	rel := r2.Vec{}
	for _, op := range s.Groups[0].Paths[0].Operations {
		switch op.Type {
		case svg.Move:
			pt := op.Points[0]
			rel.X, rel.Y = pt.X, pt.Y
			curve = append(curve, r2.Vec{
				X: pt.X,
				Y: pt.Y,
			})
		case svg.MoveRel:
			pt := op.Points[0]
			rel.X, rel.Y = rel.X+pt.X, rel.X+pt.Y
		case svg.Cubic:
			for _, pt := range op.Points {
				curve = append(curve, r2.Vec{
					X: pt.X,
					Y: pt.Y,
				})
				rel.X = pt.X
				rel.Y = pt.Y
			}
		case svg.CubicRel:
			for _, pt := range op.Points {
				curve = append(curve, r2.Vec{
					X: pt.X + rel.X,
					Y: pt.Y + rel.X,
				})
			}
		}
		// seems like this should be reset...
		//rel.X, rel.Y = 0.0, 0.0
	}
	return curve
}

func ShiftPoints(pts []complex128) []complex128 {
	/*var center r2.Vec
	for _, p := range pts {
		center.Add(p)
	}

	center.Scale(-1.0 / float64(len(pts)))*/

	out := make([]complex128, len(pts))
	for i := range out {
		out[i] = coordsys.ImgToUnitC(pts[i])
	}

	return out
}

// t [0, 1]
func qBezier(t float64, a, b, c r2.Vec) r2.Vec {
	omT := 1.0 - t
	p0 := a.Scale(omT * omT)
	p1 := b.Scale(2 * t * omT)
	p2 := b.Scale(t * t)
	return p0.Add(p1).Add(p2)
}

func cBezier(t float64, a, b, c, d r2.Vec) r2.Vec {
	omT := 1.0 - t
	p0 := qBezier(t, a, b, c).Scale(omT)
	p1 := qBezier(t, b, c, d).Scale(t)
	return p0.Add(p1)
}

func (se *SvgEpicycle) Frame(t float64) image.Image {
	/*img, gc := jimage.New()
	draw.Draw(img, img.Bounds(), &image.Uniform{palette.Palette[0]}, image.ZP, draw.Src)

	perFrame := len(se.pts) / global.FrameCount
	DrawCmplxLines(gc, se.pts, se.F*perFrame)
	se.DrawEpiCircles(gc, t)

	se.F++*/

	se.anim.Frame(t)

	return se.currImg
}

func (se *SvgEpicycle) DrawOp(gc *draw2dimg.GraphicContext) {
	rel := r2.Vec{}
	for _, op := range se.operations {
		switch op.Type {
		case svg.Move:
			pt := op.Points[0]
			gc.MoveTo(pt.X, pt.Y)
			rel.X, rel.Y = pt.X, pt.Y
		case svg.MoveRel:
			pt := op.Points[0]
			gc.MoveTo(pt.X, pt.Y)
		case svg.Cubic:
			for i := 0; i < len(op.Points); i += 3 {
				x0, x1, x := op.Points[i], op.Points[i+1], op.Points[i+2]
				gc.CubicCurveTo(x0.X, x0.Y, x1.X, x1.Y, x.X, x.Y)
			}
		case svg.CubicRel:
			for i := 0; i < len(op.Points); i += 3 {
				x0, x1, x := op.Points[i], op.Points[i+1], op.Points[i+2]
				x0.X, x0.Y = x0.X+rel.X, x0.Y+rel.Y
				x1.X, x1.Y = x1.X+rel.X, x1.Y+rel.Y
				x.X, x.Y = x.X+rel.X, x.Y+rel.Y
				gc.CubicCurveTo(x0.X, x0.Y, x1.X, x1.Y, x.X, x.Y)
			}
		}
	}
	gc.Stroke()
}

func DrawLines(gc *draw2dimg.GraphicContext, pts []r2.Vec) {
	// draw line through all points
	p := coordsys.UnitToImg(pts[0])
	gc.MoveTo(p.X, p.Y)
	for i := 0; i < len(pts); i++ {
		p = coordsys.UnitToImg(pts[i])
		gc.LineTo(p.X, p.Y)
	}
	gc.Stroke()
}

func DrawCmplxLines(gc *draw2dimg.GraphicContext, pts []complex128, count int) {
	if count >= len(pts) {
		count = len(pts) - 1
	}

	gc.SetStrokeColor(palette.Palette[3])
	// draw line through all points
	p := coordsys.UnitToImgC(pts[0])
	gc.MoveTo(real(p), imag(p))
	for i := 0; i < count; i++ {
		p := coordsys.UnitToImgC(pts[i])
		gc.LineTo(real(p), imag(p))
	}
	gc.Stroke()
}

func (se *SvgEpicycle) DrawEpiCircles(gc *draw2dimg.GraphicContext, t float64, count int) {

	if count >= len(se.sorted) {
		count = len(se.sorted) - 1
	}

	h := len(se.coeff) / 2
	center := complex(0, 0)
	for i := 0; i < count; i++ {
		srt := se.sorted[i]
		p := fourier.Pat(t, se.coeff, srt.idx)

		// don't draw the static one
		if srt.idx != h {
			DrawCCirc(gc, center, p)
		}

		center += p
	}

	center = complex(0, 0)
	for i := 0; i < count; i++ {
		srt := se.sorted[i]
		p := fourier.Pat(t, se.coeff, srt.idx)

		if srt.idx != h {
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

func CmplxSliceToVec(cslic []complex128) []r2.Vec {
	out := make([]r2.Vec, len(cslic))
	for i, c := range cslic {
		out[i].X = real(c)
		out[i].Y = imag(c)
	}
	return out
}
