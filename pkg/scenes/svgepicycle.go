package scenes

import (
	"fmt"
	"image"
	"image/draw"
	"math/cmplx"
	"os"

	"github.com/HexHacks/goffer/pkg/global"
	jmath "github.com/HexHacks/goffer/pkg/math"
	"github.com/HexHacks/goffer/pkg/palette"
	"github.com/HexHacks/goffer/pkg/svg"
	"github.com/llgcode/draw2d/draw2dimg"
	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/spatial/r2"

	jimage "github.com/HexHacks/goffer/pkg/image"
)

type SvgEpicycle struct {
	operations []svg.Operation
	origPts    []r2.Vec
	svgPts     []complex128
	coeff      []complex128
	fft        *fourier.CmplxFFT

	// running the fourier version will leave us with some
	// points that we want to draw lines through
	pts []r2.Vec
}

func (se *SvgEpicycle) Init() {

	fil, err := os.Open("assets/gubb_abs.svg")
	panicOn(err)
	defer fil.Close()

	// out point when drawing
	se.pts = make([]r2.Vec, 0, global.FrameCount)

	s, err := svg.ParseSvg(fil)
	panicOn(err)

	se.operations = s.Groups[0].Paths[0].Operations
	se.origPts = ExtractPoints(s)
	se.svgPts = ExpandCurve(se.origPts, 100, 3.8)

	se.fft = fourier.NewCmplxFFT(len(se.svgPts))
	se.coeff = se.fft.Coefficients(nil, se.svgPts)
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

func (se *SvgEpicycle) InverseFTransform(s float64) r2.Vec {
	out := 0 + 0i
	w := complex(s*jmath.Tau/global.Total, 0)
	for k, c := range se.coeff {
		kk := complex(float64(k), 0)
		freq := complex(se.fft.Freq(k), 0)

		next := c * cmplx.Exp(w*freq*kk)
		out = out + next
	}

	// scale
	out = out / complex(float64(len(se.coeff)), 0)

	return r2.Vec{
		real(out),
		imag(out),
	}
}

func (se *SvgEpicycle) Frame(s float64) image.Image {
	img, gc := jimage.New()
	draw.Draw(img, img.Bounds(), &image.Uniform{palette.Palette[0]}, image.ZP, draw.Src)

	// make new point and add it to list
	v := se.InverseFTransform(s)
	se.pts = append(se.pts, v)

	//se.DrawOp(gc)
	DrawLines(gc, CmplxSliceToVec(se.svgPts))
	//DrawLines(gc, se.pts)

	return img
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
	gc.MoveTo(pts[0].X, pts[0].Y)
	for i := 0; i < len(pts); i++ {
		gc.LineTo(pts[i].X, pts[i].Y)
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
