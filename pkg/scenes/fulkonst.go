package main

import (
	"image"
	"math"
	"math/rand"

	kit "github.com/llgcode/draw2d/draw2dkit"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/spatial/r2"
)

type frameFulkonstOne struct{}

func (f frameFulkonstOne) Init() {}
func (f frameFulkonstOne) Frame(t float64) image.Image {
	var (
		w, h   = float64(Width), float64(Height)
		cx, cy = w / 2.0, h / 2.0
	)

	img, gc := drawCommon(Palette)

	rad := TwoPi * t / Total
	amp := w / 2.0
	cosp, cosc := amp*math.Cos(rad)+cx, amp*math.Cos(rad/2.0)+cx
	sinp, sinc := amp*math.Sin(rad)+cy, amp*math.Sin(rad/2.0)+cy

	// Draw a closed shape
	gc.BeginPath()
	gc.MoveTo(cx, cy)
	gc.LineTo(w, h/2.0)
	gc.QuadCurveTo(cosc, sinc, cosp, sinp)
	gc.Close()
	gc.FillStroke()

	return img
}

type frameFulkonstTwo struct {
	Graph *simple.DirectedGraph
	Roots []*FixNode
	Edges map[int64]*WeightEdge
}

func (f *frameFulkonstTwo) Init() {
	const (
		subCount = 3
		dist     = 52.0
		hdist    = dist / 2.0
	)

	f.Graph = simple.NewDirectedGraph()
	f.Edges = make(map[int64]*WeightEdge)

	f.Roots = []*FixNode{
		NewFixNode(r2.Vec{X: CX + hdist, Y: CY + hdist}),
		NewFixNode(r2.Vec{X: CX - hdist, Y: CY + hdist}),
		NewFixNode(r2.Vec{X: CX - hdist, Y: CY - hdist}),
		NewFixNode(r2.Vec{X: CX + hdist, Y: CY - hdist}),
	}

	for _, r := range f.Roots {
		f.Graph.AddNode(r)
	}

	for r := 0; r < 4; r++ {
		var prev graph.Node = f.Roots[r]
		for i := 0; i < subCount; i++ {
			prevP := prev.(Positioner).Pos()
			newNode := NewFluidNode(jexpV(rand.Float64(), dist, prevP.X, prevP.Y),
				jexpV(rand.Float64(), 2.0, 0.0, 0.0), // vel
				func(t float64, fn *FluidNode) {
					var parent Positioner

					e := f.Edges[fn.ID()]
					parent = e.From().(Positioner)
					parentP := parent.Pos()

					newPos := fn.Base.P.Add(fn.V)

					// normalize lenght
					fn.Base.P = parentP.Add(normalize(newPos.Sub(parentP)).Scale(e.L))
				},
			)
			f.Graph.AddNode(newNode)

			e := NewWeightEdge(prev, newNode, length(prevP, newNode.Pos()))
			f.Graph.SetEdge(e)
			f.Edges[newNode.ID()] = e

			prev = newNode
		}
	}
}

func (f *frameFulkonstTwo) Frame(t float64) image.Image {
	const (
		csiz = 5.0
	)

	img, gc := drawCommon(Palette)

	eiter := f.Graph.Edges()
	for eiter.Next() {
		e := eiter.Edge()
		a := e.From().(Positioner).Pos()
		b := e.To().(Positioner).Pos()

		gc.MoveTo(a.X, a.Y)
		gc.LineTo(b.X, b.Y)
		gc.Stroke()
	}

	niter := f.Graph.Nodes()
	for niter.Next() {
		n := niter.Node()
		v := n.(Positioner).Pos()

		kit.Circle(gc, v.X, v.Y, csiz)
		gc.FillStroke()

		if u, ok := n.(Updater); ok {
			u.Update(t)
		}
	}

	return img
}

func randPt(cx, cy, w, h float64) r2.Vec {
	return r2.Vec{
		X: rand.Float64()*w + cx/2.0,
		Y: rand.Float64()*h + cy/2.0,
	}
}
