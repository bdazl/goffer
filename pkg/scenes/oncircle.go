package scenes

import (
	//"fmt"
	"image"
	"image/color"
	"math"

	"github.com/bdazl/goffer/pkg/global"
	jimage "github.com/bdazl/goffer/pkg/image"
	"github.com/bdazl/goffer/pkg/math/float"
	jr2 "github.com/bdazl/goffer/pkg/math/r2"
	"github.com/bdazl/goffer/pkg/palette"

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
		Root: newCNode(jr2.ExpV(0.0, r, t.X, t.Y)),
		r:    r,
		t:    t,
	}

	last := out.Root
	for i := 1; i < pts; i++ {
		last.Next = newCNode(jr2.ExpV(float64(i)/ptsf, r, t.X, t.Y))
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
		circs = 3
		r0    = 70.0
		distB = 20.0
		pts   = 13
	)
	var (
		c = jr2.V(global.CX, global.CY)
	)

	o.Circs = make([]*OCirc, circs)
	for i := range o.Circs {
		circ := newOCirc(pts, r0*float64(i+1), c)

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

func (o *OnCircle0) Frame(t float64) image.Image {
	const (
		amp   = 0.09
		sigma = 2.0
		is    = 0.33
		isv   = 0.1
	)
	var (
		tzo = t / (global.Total - global.DT)
	)

	img, gc := jimage.New()

	for i, c := range o.Circs {
		fi := float64(i)
		fis := fi * is // :)
		c.Apply(func(n *CNode, f float64) {

			// make the wavelet repeat saw: [0, 1] -> [-1, 1] (repeating)
			saw := 2.0*math.Mod((f+tzo)*2, 2.0) - 1.0
			r := c.r + amp*c.r*float.MorletBnd(sigma, saw)

			f = f + fis

			n.P = jr2.ExpV(f, r, global.CX, global.CY)
		})
		c.Render(gc, palette.Palette)
	}

	return img
}
