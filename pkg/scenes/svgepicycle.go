package scenes

import (
	"image"
	"math/cmplx"
	"os"

	"github.com/HexHacks/goffer/pkg/global"
	jmath "github.com/HexHacks/goffer/pkg/math"
	"github.com/HexHacks/goffer/pkg/svg"
	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/spatial/r2"

	jimage "github.com/HexHacks/goffer/pkg/image"
)

type SvgEpicycle struct {
	svgPts []complex128
	coeff  []complex128
	fft    *fourier.CmplxFFT

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

	curve := ExtractPoints(s)
	se.svgPts = ExpandCurve(curve, 10)

	se.fft = fourier.NewCmplxFFT(len(se.svgPts))
	se.coeff = se.fft.Coefficients(nil, se.svgPts)
}

func ExpandCurve(curve []r2.Vec, perBezier int) []complex128 {
	pB := float64(perBezier)

	l := len(curve)
	out := make([]complex128, 0, l*perBezier)
	for i := 0; i < l; i = i + 2 {
		if l <= i+2 {
			break
		}
		a, b, c := curve[i], curve[i+1], curve[i+2]

		for j := 0.0; j < pB; j = j + 1.0 {
			p := bezier(j/pB, a, b, c)
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
		rel.X, rel.Y = 0.0, 0.0
	}
	return curve
}

// t [0, 1]
func bezier(t float64, a, b, c r2.Vec) r2.Vec {
	v, w := b.Sub(a), c.Sub(b)

	n0, n1 := a.Add(v.Scale(t)), b.Add(w.Scale(t))
	return n0.Add(n1.Sub(n0).Scale(t))
}

func (se *SvgEpicycle) Fourier(s float64) r2.Vec {
	out := 0 + 0i
	w := complex(s*jmath.Tau/global.Total, 0)
	for i, c := range se.coeff {
		n := complex(float64(se.fft.ShiftIdx(i)), 0)
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

	// make new point and add it to list
	v := se.Fourier(s)
	se.pts = append(se.pts, v)

	// draw line through all points
	gc.MoveTo(se.pts[0].X, se.pts[0].Y)
	for i := 0; i < len(se.pts); i++ {
		gc.LineTo(se.pts[i].X, se.pts[i].Y)
	}
	gc.Close()
	gc.Stroke()

	return img
}
