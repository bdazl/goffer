package scenes

import (
	"image"
	"math"
	"math/rand"

	"github.com/bdazl/goffer/pkg/global"
	jr2 "github.com/bdazl/goffer/pkg/math/r2"

	jgraph "github.com/bdazl/goffer/pkg/graph"
	jimage "github.com/bdazl/goffer/pkg/image"
	jmath "github.com/bdazl/goffer/pkg/math"

	kit "github.com/llgcode/draw2d/draw2dkit"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/spatial/r2"
)

type frameFulkonstOne struct{}

func (f frameFulkonstOne) Init() {}
func (f frameFulkonstOne) Frame(t float64) image.Image {
	var (
		w, h   = global.W, global.H
		cx, cy = w / 2.0, h / 2.0
	)

	img, gc := jimage.New()

	rad := jmath.Tau * t / global.Total
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
	Roots []*jgraph.FixNode
	Edges map[int64]*jgraph.WeightEdge
}

func (f *frameFulkonstTwo) Init() {
	const (
		subCount = 3
		dist     = 52.0
		hdist    = dist / 2.0
	)
	var (
		cx, cy = global.CX, global.CY
	)

	f.Graph = simple.NewDirectedGraph()
	f.Edges = make(map[int64]*jgraph.WeightEdge)

	f.Roots = []*jgraph.FixNode{
		jgraph.NewFixNode(r2.Vec{X: cx + hdist, Y: cy + hdist}),
		jgraph.NewFixNode(r2.Vec{X: cx - hdist, Y: cy + hdist}),
		jgraph.NewFixNode(r2.Vec{X: cx - hdist, Y: cy - hdist}),
		jgraph.NewFixNode(r2.Vec{X: cx + hdist, Y: cy - hdist}),
	}

	for _, r := range f.Roots {
		f.Graph.AddNode(r)
	}

	for r := 0; r < 4; r++ {
		var prev graph.Node = f.Roots[r]
		for i := 0; i < subCount; i++ {
			prevP := prev.(jgraph.Positioner).Pos()
			newNode := jgraph.NewFluidNode(jr2.ExpV(rand.Float64(), dist, prevP.X, prevP.Y),
				jr2.ExpV(rand.Float64(), 2.0, 0.0, 0.0), // vel
				func(t float64, fn *jgraph.FluidNode) {
					var parent jgraph.Positioner

					e := f.Edges[fn.ID()]
					parent = e.From().(jgraph.Positioner)
					parentP := parent.Pos()

					newPos := fn.Base.P.Add(fn.V)

					// normalize lenght
					fn.Base.P = parentP.Add(jr2.Normalize(newPos.Sub(parentP)).Scale(e.L))
				},
			)
			f.Graph.AddNode(newNode)

			e := jgraph.NewWeightEdge(prev, newNode, jr2.Length(prevP, newNode.Pos()))
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

	img, gc := jimage.New()

	eiter := f.Graph.Edges()
	for eiter.Next() {
		e := eiter.Edge()
		a := e.From().(jgraph.Positioner).Pos()
		b := e.To().(jgraph.Positioner).Pos()

		gc.MoveTo(a.X, a.Y)
		gc.LineTo(b.X, b.Y)
		gc.Stroke()
	}

	niter := f.Graph.Nodes()
	for niter.Next() {
		n := niter.Node()
		v := n.(jgraph.Positioner).Pos()

		kit.Circle(gc, v.X, v.Y, csiz)
		gc.FillStroke()

		if u, ok := n.(jgraph.Updater); ok {
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
