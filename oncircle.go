package main

import (
	//"fmt"
	"image"
	"image/color"
	"math"

	//"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/spatial/r2"

	"github.com/llgcode/draw2d/draw2dimg"
	kit "github.com/llgcode/draw2d/draw2dkit"
)

type CNode struct {
	P     r2.Vec
	Size  float64
	Color int
	Next  *CNode
}

func newCNode(p r2.Vec) *CNode {
	return &CNode{
		P:     p,
		Size:  5.0,
		Color: 1,
	}
}

type OCirc struct {
	Root *CNode
	r    float64 // radius
	t    r2.Vec  // translation
}

func newOCirc(pts int, r float64, t r2.Vec) *OCirc {
	var (
		ptsf = float64(pts)
	)
	out := &OCirc{
		Root: newCNode(jexpV(0.0, r, t.X, t.Y)),
		r:    r,
		t:    t,
	}

	last := out.Root
	for i := 1; i < pts; i++ {
		last.Next = newCNode(jexpV(float64(i)/ptsf, r, t.X, t.Y))
		last = last.Next
	}

	return out
}

func (c OCirc) Render(gc *draw2dimg.GraphicContext, palette color.Palette) {
	// draw lines
	gc.SetFillColor(palette[3])
	gc.MoveTo(c.Root.P.X, c.Root.P.Y)
	next := c.Root.Next
	for next != nil {
		//fmt.Println(next.P.X, next.P.Y)
		gc.LineTo(next.P.X, next.P.Y)
		if next.Next == nil {
			//fmt.Println("next nil:", next.P.X, next.P.Y)
			gc.LineTo(c.Root.P.X, c.Root.P.Y)
		}

		next = next.Next
	}
	gc.Stroke()

	// draw circles
	next = c.Root
	for next != nil {
		p := next.P

		gc.SetFillColor(palette[next.Color])
		kit.Circle(gc, p.X, p.Y, next.Size)
		gc.FillStroke()

		next = next.Next
	}

}

func (c OCirc) Len() int {
	i := 0
	next := c.Root
	for next != nil {
		i++
		next = next.Next
	}
	return i
}

func (c *OCirc) Apply(fun func(n *CNode, f float64)) {
	i := 0.0
	l := float64(c.Len() - 1)
	next := c.Root
	for next != nil {
		f := i / l
		fun(next, f)

		next = next.Next
		i = i + 1.0
	}
}

type OnCircle0 struct {
	Circs []*OCirc
}

func (o *OnCircle0) Init() {
	const (
		circs = 2
		r0    = 70.0
		distB = 20.0
		pts   = 13
	)

	o.Circs = make([]*OCirc, circs)
	for i := range o.Circs {
		circ := newOCirc(pts, r0*float64(i+1), C)

		// set colors
		p := 0
		next := circ.Root
		for next != nil {
			if (p+i)%2 == 0 {
				next.Color = 3
			}

			p++
			next = next.Next
		}

		o.Circs[i] = circ
	}
}

func (o *OnCircle0) Frame(t float64) *image.Paletted {
	const (
		amp   = 0.1
		sigma = 2.0
		is    = 10.0
		fs    = 8.0
	)
	var (
		palette = Palette1
	)

	img, gc := drawCommon(palette)

	for _, c := range o.Circs {
		//fi := float64(i)
		//fis := fi * is // :)
		c.Apply(func(n *CNode, f float64) {
			//f = f + 0.1*math.Sin(TwoPi*t*0.1)
			tzo := t / (Total - DT)
			saw := 2.0*math.Mod((f+tzo)*2, 2.0) - 1.0
			r := c.r + amp*c.r*MorletBnd(sigma, saw)
			n.P = jexpV(f, r, CX, CY)
		})
		c.Render(gc, palette)
	}

	return gifEncodeFrame(img, palette)
}
