package scenes

import (
	"fmt"
	"image"
	"os"

	"github.com/bdazl/goffer/pkg/animation/epicycles"
	"github.com/bdazl/goffer/pkg/coordsys"
	"github.com/bdazl/goffer/pkg/math/fourier"
	jr2 "github.com/bdazl/goffer/pkg/math/r2"
	"github.com/bdazl/goffer/pkg/svg"

	"github.com/llgcode/draw2d/draw2dimg"
	"gonum.org/v1/gonum/spatial/r2"
)

type SvgEpicycle struct {
	F          int
	operations []svg.Operation
	origPts    []r2.Vec
	svgPts     []complex128
	shift      []complex128
	coeff      []complex128

	pts []complex128

	epi *epicycles.Epicycles
}

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

	se.epi = epicycles.New(se.coeff, len(se.shift))
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
			p := jr2.CubeBezier(j/pB, a, b, c, d)
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

func (se *SvgEpicycle) Frame(t float64) image.Image {
	return se.epi.Frame(t)
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

func CmplxSliceToVec(cslic []complex128) []r2.Vec {
	out := make([]r2.Vec, len(cslic))
	for i, c := range cslic {
		out[i].X = real(c)
		out[i].Y = imag(c)
	}
	return out
}
