package main

import (
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
)

func nextPtState(pse PtStateEnum) PtStateEnum {
	switch pse {
	case Bend:
		return Duplicate
	case Duplicate:
		return FlipNewHalf
	case FlipNewHalf:
		return Bend
	}
	return Bend
}

type PtState struct {
	PtStateEnum
	Pts []complex128
}

type PtBend0 struct {
	Now, Next PtState
	NowT      PtStateEnum

	Entered bool
	Zo      float64
}

func (p *PtBend0) Init() {
	const (
		lineCount  = 10
		lineCountf = float64(lineCount)
		lineDist   = 15.0
		hlineDist  = lineDist / 2.0
		ptCount    = lineCount * lineCount
	)

	p.Now.Pts = make([]complex128, ptCount)
	p.Next.Pts = make([]complex128, ptCount)

	for y := 0; y < lineCount; y++ {
		for x := 0; x < lineCount; x++ {
			idx := y*lineCount + x

			cx := float64(lineDist*x)/(lineCountf-1.0) - hlineDist
			cy := float64(lineDist*y)/(lineCountf-1.0) - hlineDist
			pt := complex(cx, cy)
			p.Now.Pts[idx] = pt
			p.Next.Pts[idx] = cmplx.Sqrt(pt)
		}
	}
}

func (p *PtBend0) TransitionTo(pse PtStateEnum) {
	switch pse {
	case Bend:
	case Duplicate:
	case FlipNewHalf:
	}
}

func (p *PtBend0) NextState() {
	next := nextPtState(p.NowT)

	p.TransitionTo(next)
	p.Now = p.Next
	p.NowT = next
	p.Zo = 0.0
}

func (p *PtBend0) Step() {
	const (
		stateTime = 1.0 // seconds
	)

	p.Zo = p.Zo + DT
	if p.Zo >= stateTime {
		p.NextState()
	}
}

func (p *PtBend0) Render(gc *draw2dimg.GraphicContext) {
	for i, now := range p.Now.Pts {
		next := p.Next.Pts[i]

		// TODO: not linearly interpolate?
		lerp := now + complex(p.Zo, 0.0)*(next-now)
		scr := ComplexToScreen(lerp)

		gc.SetFillColor(Palette[2])
		kit.Circle(gc, scr.X, scr.Y, 5.0)
		gc.FillStroke()
	}
}

func ComplexToScreen(c complex128) r2.Vec {
	const (
		C = 10.0
	)
	return r2.Vec{
		X: (real(c) + C) * W / (2.0 * C),
		Y: (imag(c) + C) * H / (2.0 * C),
	}
}

func (p *PtBend0) Frame(t float64) image.Image {
	img, gc := drawCommon(Palette)

	p.Step()
	p.Render(gc)

	return img
}
