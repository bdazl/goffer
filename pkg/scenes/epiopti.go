package scenes

import (
	"fmt"
	"image"
	"math"
	"math/cmplx"
	"math/rand"

	"github.com/bdazl/goffer/pkg/animation/epicycles"
	jmath "github.com/bdazl/goffer/pkg/math"
	"github.com/bdazl/goffer/pkg/math/fourier"

	"gonum.org/v1/gonum/optimize"
)

const (
	bigNum = 100000.0
)

type EpiOpti struct {
	epi *epicycles.Epicycles
}

func (e *EpiOpti) Init() {
	const (
		ptCnt = 200
		cCnt  = 21
	)

	c := make([]complex128, cCnt)
	for i := range c {
		a, b := rand.Float64(), rand.Float64()
		c[i] = complex(a*0.3-0.15, b*0.3-0.15)
	}

	workBuffer := make([]complex128, cCnt)
	p := optimize.Problem{
		Func: func(x []float64) float64 {
			FloatSliceToCmplx(workBuffer, x)

			sum := complex(0, 0)
			for _, p := range workBuffer {
				sum += p
			}

			// I want the initial condition to be bigger than zero but less than
			// the size of the screen.. I need to figure out how conditional optimization works
			extra := 0.0
			l := cmplx.Abs(sum)
			if l < 0.5 || l > math.Sqrt2 {
				extra = (2.0 + l) * (2.0 + l)
			}

			// draw epicycle circle points from workBuffer as coefs
			epiPts := drawEpiPoints(workBuffer)

			// complare points to minimal distance from points as generated by curve
			return extra + cmpEpicycle(workBuffer, epiPts)
		},
	}

	x := CmplxSliceToFloat(c)
	result, err := optimize.Minimize(p, x, nil, nil)
	panicOn(err)
	panicOn(result.Status.Err())

	fmt.Printf("result.Status: %v\n", result.Status)
	fmt.Printf("result.X: %0.4g\n", result.X)
	fmt.Printf("result.F: %0.4g\n", result.F)
	fmt.Printf("result.Stats.FuncEvaluations: %d\n", result.Stats.FuncEvaluations)

	FloatSliceToCmplx(workBuffer, result.X)

	e.epi = epicycles.New(workBuffer, ptCnt)
}

func cmpEpicycle(C []complex128, expPts []complex128) float64 {
	out := 0.0
	cnt := float64(len(expPts))
	for i, p := range expPts {
		t := float64(i) / cnt
		P := fourier.P(t, C)
		abs := cmplx.Abs(P - p)
		out += abs * abs
	}
	return out
}

func drawEpiPoints(C []complex128) []complex128 {
	const (
		ptsInCircle = 10
		ptsInLine   = 4
		ptsPerCoef  = ptsInCircle + ptsInLine
	)

	sorted := epicycles.MakeSorted(C)

	ptCount := ptsPerCoef * len(C)

	out := make([]complex128, 0, ptCount)
	sum := complex(0, 0)
	for _, s := range sorted {
		addLine(&out, sum, sum+s.C, ptsInLine)
		addCircle(&out, sum, s.C, ptsInCircle)
		sum += s.C
	}

	return out
}

func addLine(out *[]complex128, a, b complex128, count int) {
	cnt := float64(count)
	for i := 0; i < count; i++ {
		f := float64(i) / cnt

		p := a + complex(f, 0)*(b-a)
		*out = append(*out, p)
	}
}

func addCircle(out *[]complex128, pt, c complex128, count int) {
	cnt := float64(count)
	for i := 0; i < count; i++ {
		f := float64(i) / cnt

		e := complex(0, jmath.Tau*f)
		p := pt + c*cmplx.Exp(e)
		*out = append(*out, p)
	}
}

func CmplxSliceToFloat(C []complex128) []float64 {
	out := make([]float64, len(C)*2)
	for i, c := range C {
		out[i*2] = real(c)
		out[i*2+1] = imag(c)
	}
	return out
}

func FloatSliceToCmplx(out []complex128, in []float64) {
	if len(out)*2 != len(in) {
		panic("bad input")
	}

	for i := range out {
		out[i] = complex(in[i*2], in[i*2+1])
	}
}

func (e *EpiOpti) Frame(t float64) image.Image {
	return e.epi.Frame(t)
}
