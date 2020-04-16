package main

import (
	"fmt"
	"image"
	"math"
	"math/cmplx"
	"math/rand"

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

	oldMax, newMax float64

	Changes float64
	Doubles float64
	Entered bool
	Zo      float64
}

func (p *PtBend0) Init() {
	const (
		lineDist   = 0.8
		lineCount  = 10
		lineCountf = float64(lineCount)
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

			cx := lineDist*float64(x)/(lineCountf-1.0) - hlineDist
			cy := lineDist*float64(y)/(lineCountf-1.0) - hlineDist
			pt := complex(cx, cy)
			p.Now[idx] = pt
			addMax(&p.oldMax, p.Now[idx])

			p.Next[idx] = ptBend(pt)
			addMax(&p.newMax, p.Next[idx])
		}
	}
}

func (p *PtBend0) TransitionTo(pse PtStateEnum) {
	p.NowT = pse
	p.NextT = nextPtState(pse)

	p.oldMax = p.newMax
	p.newMax = 0

	switch pse {
	case Bend:
		for i, now := range p.Now {
			p.Next[i] = ptBend(cmplx.Sqrt(now))
			addMax(&p.newMax, p.Next[i])
		}

	case Duplicate:
		p.newMax = p.oldMax // nothing will change
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
		p.newMax = p.oldMax // noting will change
		// conjugate the new points from duplication step
		for i := len(p.Next) / 2; i < len(p.Next); i++ {
			p.Next[i] = cmplx.Conj(p.Next[i])
		}

	case Translate:
		a, b := rand.Float64()*2.0-1.0, rand.Float64()*2.0-1.0

		// conjugate the new points from duplication step
		for i := 0; i < len(p.Next); i++ {
			//p.Next[i] = p.Next[i] + 1 + 1i
			p.Next[i] = p.Next[i] + complex(a, b)
			addMax(&p.newMax, p.Next[i])
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
		p.oldMax = p.newMax
	}
	fmt.Println("post transition:", p.NowT, p.NextT)
}

func addMax(f *float64, c complex128) {
	a, b := real(c), imag(c)
	a = math.Max(a, b)
	*f = math.Max(*f, a)
}

func ptBend(c complex128) complex128 {

	//return cmplx.Sqrt(c) // first bend
	// ctwo := (c + 2)
	//return ctwo * ctwo * (c - 1 - 2i) * (c + 1i)
	fact := (c + 1)
	return fact * fact
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
		sPow  = 1.1
	)
	var (
		// find max values
		sstep = smoothstep(0.0, 1.0, p.Zo)
		ma    = (p.oldMax + sstep*(p.newMax-p.oldMax)) * 1.3
		//nolen = float64(len(p.Now))
		//nelen = float64(len(p.Next))
		siz = 2.5
	)

	for i, now := range p.Now {
		next := p.Next[i]

		// lerp point
		lerp := now + complex(sstep, 0.0)*(next-now)

		// lerp size

		scr := ComplexToScreen(lerp, ma)

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
