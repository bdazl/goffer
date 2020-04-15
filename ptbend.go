package main

import (
	"fmt"
	"image"
	"math/cmplx"

	"github.com/llgcode/draw2d/draw2dimg"
	kit "github.com/llgcode/draw2d/draw2dkit"
	"gonum.org/v1/gonum/spatial/r2"
)

type PtStateEnum int

const (
	Bend PtStateEnum = iota
	Duplicate
	FlipNewHalf
	Translate
)

func nextPtState(pse PtStateEnum) PtStateEnum {
	switch pse {
	case Bend:
		return Duplicate
	case Duplicate:
		return FlipNewHalf
	case FlipNewHalf:
		return Translate
	case Translate:
		return Bend
	default:
		panic("bad stat")
	}
}

func (s PtStateEnum) String() string {
	switch s {
	case Bend:
		return "bend"
	case Duplicate:
		return "dupl"
	case FlipNewHalf:
		return "flip"
	case Translate:
		return "transl"
	}
	return "no idea"
}

type PtBend0 struct {
	Now, Next   []complex128
	NowT, NextT PtStateEnum

	Changes float64
	Doubles float64
	Entered bool
	Zo      float64
}

func (p *PtBend0) Init() {
	const (
		lineCount  = 10
		lineCountf = float64(lineCount)
		lineDist   = 20.0
		hlineDist  = lineDist / 2.0
		ptCount    = lineCount * lineCount
	)

	p.NextT = Duplicate
	p.Doubles = 1.0

	p.Now = make([]complex128, ptCount)
	p.Next = make([]complex128, ptCount)

	for y := 0; y < lineCount; y++ {
		for x := 0; x < lineCount; x++ {
			idx := y*lineCount + x

			cx := float64(lineDist*x)/(lineCountf-1.0) - hlineDist
			cy := float64(lineDist*y)/(lineCountf-1.0) - hlineDist
			pt := complex(cx, cy)
			p.Now[idx] = pt
			p.Next[idx] = cmplx.Sqrt(pt)
		}
	}
}

func (p *PtBend0) TransitionTo(pse PtStateEnum) {
	p.NowT = pse
	p.NextT = nextPtState(pse)

	switch pse {
	case Bend:
		for i, now := range p.Now {
			p.Next[i] = cmplx.Sqrt(now)
		}

	case Duplicate:
		p.Doubles = p.Doubles + 1.0

		// extend next
		p.Next = make([]complex128, 2*len(p.Now))
		for i := range p.Next {
			p.Next[i] = p.Now[i%len(p.Now)]
		}

		// swap back to now (duplicate point on the same spot)
		p.Now = make([]complex128, len(p.Next))
		copy(p.Now, p.Next)

		// flip new points of next
		for i := len(p.Next) / 2; i < len(p.Next); i++ {
			p.Next[i] = -cmplx.Conj(p.Next[i])
		}

	case FlipNewHalf:
		// conjugate the new points from duplication step
		for i := len(p.Next) / 2; i < len(p.Next); i++ {
			p.Next[i] = cmplx.Conj(p.Next[i])
		}

	case Translate:
		// conjugate the new points from duplication step
		for i := 0; i < len(p.Next); i++ {
			p.Next[i] = p.Next[i] + (1 + 1i)
		}
	}
}

func (p *PtBend0) NextState() {
	// previous points always goes to now state
	copy(p.Now, p.Next)

	fmt.Println("pre transition:", p.NowT, p.NextT)
	if p.NowT == p.NextT {
		p.TransitionTo(p.NowT)
		p.Changes = p.Changes + 1.0
	} else {
		// do a pause (by fooling ourselves that we are in the same state)
		p.NowT = p.NextT
	}
	fmt.Println("post transition:", p.NowT, p.NextT)
}

func (p *PtBend0) Step() {
	const (
		stateTime = 1.0 // seconds
	)

	p.Zo = p.Zo + DT
	if p.Zo >= stateTime {
		p.NextState()
		p.Zo = 0.0
	}
}

func (p *PtBend0) Render(gc *draw2dimg.GraphicContext) {
	const (
		sBase = 20.0
	)
	for i, now := range p.Now {
		next := p.Next[i]

		// TODO: not linearly interpolate?
		lerp := now + complex(smoothstep(0.0, 1.0, p.Zo), 0.0)*(next-now)

		siz := sBase / (p.Doubles * 2.0)
		if p.NextT == Duplicate {
			zo := p.Zo
			if p.NowT == p.NextT {
				zo = 1.0
			}

			siz = sBase / (2 * (p.Doubles + smoothstep(0.0, 1.0, zo)))
		}

		scr := ComplexToScreen(lerp, siz)

		gc.SetFillColor(Palette[1+i%3])
		kit.Circle(gc, scr.X, scr.Y, siz)
		gc.FillStroke()
	}
}

func ComplexToScreen(c complex128, ma float64) r2.Vec {
	return r2.Vec{
		X: (real(c) + ma) * W / (2.0 * ma),
		Y: (imag(c) + ma) * H / (2.0 * ma),
	}
}

func (p *PtBend0) Frame(t float64) image.Image {
	img, gc := drawCommon(Palette)

	p.Step()
	p.Render(gc)

	return img
}
