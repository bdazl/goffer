package scenes

import (
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
	origPts []r2.Vec
	svgPts  []complex128
	coeff   []complex128
	fft     *fourier.CmplxFFT

	// running the fourier version will leave us with some
	// points that we want to draw lines through
	pts []r2.Vec
}

func (se *SvgEpicycle) Init() {
	fil, err := os.Open("assets/gubb.svg")
	panicOn(err)
	defer fil.Close()

	se.pts = make([]r2.Vec, 0, global.FrameCount)

	s, err := svg.ParseSvg(fil)
	panicOn(err)

	se.origPts = ExtractPoints(s)
	se.svgPts = ExpandCurve(se.origPts, 10)

	se.fft = fourier.NewCmplxFFT(len(se.svgPts))
	se.coeff = se.fft.Coefficients(nil, se.svgPts)
}

func ExpandCurve(curve []r2.Vec, perBezier int) []complex128 {
	pB := float64(perBezier)

	l := len(curve)
	out := make([]complex128, 0, l*perBezier)
	for i := 0; i < l; i = i + 3 {
		if l <= i+3 {
			break
		}
		a, b, c, d := curve[i], curve[i+1], curve[i+2], curve[i+3]

		for j := 0.0; j < pB; j = j + 1.0 {
			p := cBezier(j/pB, a, b, c, d)
			out = append(out, complex(p.X, p.Y))
		}
	}
	return out
}

func ExtractPoints(s *svg.Svg) []r2.Vec {
	curve := make([]r2.Vec, 0)
	rel := r2.Vec{}
	for _, op := range s.Groups[0].Paths[0].Operations {
		if op.Type == svg.Move {
			pt := op.Points[0]
			rel.X, rel.Y = pt.X, pt.Y
			continue
		}

		for _, pt := range op.Points {
			curve = append(curve, r2.Vec{
				X: pt.X + rel.X,
				Y: pt.Y + rel.X,
			})
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

func (se *SvgEpicycle) Fourier(s float64) r2.Vec {
	out := 0 + 0i
	w := complex(s*jmath.Tau/global.Total, 0)
	for i, c := range se.coeff {
		n := complex(float64(i-len(se.coeff)/2), 0)

		freq := complex(se.fft.Freq(i), 0)

		next := c * cmplx.Exp(w*freq*n)
		out = out + next
	}

	return r2.Vec{
		real(out),
		imag(out),
	}
}

func (se *SvgEpicycle) Frame(s float64) image.Image {
	img, gc := jimage.New()
	draw.Draw(img, img.Bounds(), &image.Uniform{palette.Palette[0]}, image.ZP, draw.Src)

	// make new point and add it to list
	v := se.Fourier(s)
	se.pts = append(se.pts, v)

	DrawLines(gc, CmplxSliceToVec(se.svgPts))

	return img
}

func DrawLines(gc *draw2dimg.GraphicContext, pts []r2.Vec) {
	// draw line through all points
	gc.MoveTo(pts[0].X, pts[0].Y)
	for i := 0; i < len(pts); i++ {
		gc.LineTo(pts[i].X, pts[i].Y)
	}
	gc.Close()
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
