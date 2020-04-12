package main

import (
	"gonum.org/v1/gonum/spatial/r2"
)

type Node interface {
	P() r2.Vec
	Step()
}

type Edge struct {
	A, B Node
	L    float64
}

type FixNode struct {
	Pos r2.Vec
}

func (n FixNode) P() r2.Vec { return n.Pos }
func (n FixNode) Step()     {}

type FluidNode struct {
	Pos     r2.Vec
	Vel     r2.Vec
	StepFun func(*FluidNode)
}

func (n *FluidNode) P() r2.Vec { return n.Pos }
func (n *FluidNode) V() r2.Vec { return n.Vel }
func (n *FluidNode) Step()     { n.StepFun(n) }

type Graph struct {
	Nodes []Node
	Edges []Edge
}

func (g *Graph) AddNode(n Node) int {
	g.Nodes = append(g.Nodes, n)
	return len(g.Nodes) - 1
}

func (g *Graph) AddEdge(u, v int, l float64) {
	if u < 0 || u >= len(g.Nodes) {
		panic("bad u")
	}
	if v < 0 || v >= len(g.Nodes) {
		panic("bad v")
	}

	nu, nv := g.Nodes[u], g.Nodes[v]

	// if it already exists, update weight
	for i, e := range g.Edges {
		if (e.A == nu && e.B == nv) ||
			(e.B == nu && e.A == nv) {
			g.Edges[i].L = l
			return
		}
	}

	g.Edges = append(g.Edges, Edge{A: nu, B: nv, L: l})
}

func (g *Graph) GetNodeEdges(n Node) []Edge {
	out := make([]Edge, 0, len(g.Nodes))
	for _, e := range g.Edges {
		if e.A == n || e.B == n {
			out = append(out, e)
		}
	}
	return out
}
