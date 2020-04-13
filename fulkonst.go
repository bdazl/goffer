package main

import (
	"image"
	"math"
	"math/rand"

	kit "github.com/llgcode/draw2d/draw2dkit"

	"gonum.org/v1/gonum/spatial/r2"
)

type frameFulkonstOne struct{}

func (f frameFulkonstOne) Init() {}
func (f frameFulkonstOne) Frame(t float64) *image.Paletted {
	var (
		palette = Palette1
		w, h    = float64(Width), float64(Height)
		cx, cy  = w / 2.0, h / 2.0
	)

	img, gc := drawCommon(palette)

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

	return gifEncodeFrame(img, palette)
}

type frameFulkonstTwo struct {
	Graph
}

func (f *frameFulkonstTwo) Init() {
	const (
		subCount = 3
		dist     = 52.0
		hdist    = dist / 2.0
	)

	f.Graph.Nodes = []Node{
		&FixNode{Pos: r2.Vec{X: CX + hdist, Y: CY + hdist}},
		&FixNode{Pos: r2.Vec{X: CX - hdist, Y: CY + hdist}},
		&FixNode{Pos: r2.Vec{X: CX - hdist, Y: CY - hdist}},
		&FixNode{Pos: r2.Vec{X: CX + hdist, Y: CY - hdist}},
	}

	for r := 0; r < 4; r++ {
		prevI := r
		prevP := f.Graph.Nodes[r].P()
		for i := 0; i < subCount; i++ {
			newNode := &FluidNode{
				// random direction
				Pos: jexpV(rand.Float64(), dist, prevP.X, prevP.Y),
				Vel: jexpV(rand.Float64(), 2.0, 0.0, 0.0),
				StepFun: func(fn *FluidNode) {
					// TODO
					edges := f.Graph.GetNodeEdges(fn)
					var me Edge
					for _, e := range edges {
						if e.A != fn {
							me = e
							break
						}
					}

					fn.Pos = fn.Pos.Add(fn.Vel)
					parentP := me.A.P()

					// normalize lenght
					fn.Pos = parentP.Add(normalize(fn.Pos.Sub(parentP)).Scale(me.L))
				},
			}
			ni := f.Graph.AddNode(newNode)
			f.Graph.AddEdge(prevI, ni, length(prevP, newNode.Pos))

			prevI = ni
			prevP = newNode.Pos
		}
	}
}

func (f *frameFulkonstTwo) Frame(t float64) *image.Paletted {
	const (
		csiz = 5.0
	)
	var (
		palette = Palette1
	)

	img, gc := drawCommon(palette)

	for _, e := range f.Graph.Edges {
		a := e.A.P()
		b := e.B.P()

		gc.MoveTo(a.X, a.Y)
		gc.LineTo(b.X, b.Y)
		gc.Stroke()
	}

	for _, n := range f.Graph.Nodes {
		v := n.P()
		kit.Circle(gc, v.X, v.Y, csiz)
		gc.FillStroke()

		n.Step()
	}

	return gifEncodeFrame(img, palette)
}

func randPt(cx, cy, w, h float64) r2.Vec {
	return r2.Vec{
		X: rand.Float64()*w + cx/2.0,
		Y: rand.Float64()*h + cy/2.0,
	}
}
