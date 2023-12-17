package animation

import (
	"github.com/bdazl/goffer/pkg/global"
)

type Animator func(float64)

type Animation struct {
	Sequence  []Animator // Animate this slice, in order
	Times     []float64
	currId    int
	currStart float64
}

func New(sequence []Animator, times []float64) *Animation {
	return &Animation{
		Sequence: sequence,
		Times:    times,
	}
}

func (a *Animation) Frame(t float64) {
	curr := a.Sequence[a.currId]
	tim := a.Times[a.currId] * global.Total

	nodeT := t - a.currStart
	if nodeT >= tim {
		next := a.Next(t)
		next(0.0) // because t - t = 0
	} else {
		curr(nodeT / tim)
	}
}

func (a *Animation) Next(t float64) Animator {
	a.currStart = t
	a.currId++
	if a.currId >= len(a.Sequence) {
		a.currId = 0 // safe side
	}

	next := a.Sequence[a.currId]

	return next
}
