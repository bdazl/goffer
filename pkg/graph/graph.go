package graph

import (
	"sync/atomic"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/spatial/r2"
)

var (
	nextNodeID int64 = 0
)

type Positioner interface {
	Pos() r2.Vec
}

type Velocitor interface {
	Vel() r2.Vec
}

type Updater interface {
	Update(float64)
}

type WeightEdge struct {
	*Edge
	L float64
}

func NewWeightEdge(f, t graph.Node, l float64) *WeightEdge {
	return &WeightEdge{
		Edge: NewEdge(f, t),
		L:    l,
	}
}

type FixNode struct {
	Base *Node
	P    r2.Vec
}

func (n FixNode) Pos() r2.Vec { return n.P }
func (n FixNode) ID() int64   { return n.Base.ID() }

func NewFixNode(pos r2.Vec) *FixNode {
	return &FixNode{
		Base: NewNode(),
		P:    pos,
	}
}

type FluidNode struct {
	Base    *FixNode
	V       r2.Vec
	StepFun func(float64, *FluidNode)
}

func (f *FluidNode) Pos() r2.Vec      { return f.Base.P }
func (f *FluidNode) Vel() r2.Vec      { return f.V }
func (f *FluidNode) ID() int64        { return f.Base.ID() }
func (f *FluidNode) Update(t float64) { f.StepFun(t, f) }

func NewFluidNode(pos, vel r2.Vec, step func(float64, *FluidNode)) *FluidNode {
	return &FluidNode{
		Base:    NewFixNode(pos),
		V:       vel,
		StepFun: step,
	}
}

type Edge struct {
	F, T graph.Node
}

func (e *Edge) From() graph.Node         { return e.F }
func (e *Edge) To() graph.Node           { return e.T }
func (e *Edge) ReversedEdge() graph.Edge { return &Edge{F: e.T, T: e.F} }

func NewEdge(f, t graph.Node) *Edge {
	return &Edge{F: f, T: t}
}

type Node struct {
	id int64
}

func (n Node) ID() int64 {
	return n.id
}

func NewNode() *Node {
	return &Node{
		id: NewNodeID(),
	}
}

func NewNodeID() int64 {
	return atomic.AddInt64(&nextNodeID, 1)
}
