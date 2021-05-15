package djanl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	eps = 2e-4
)

func TestBeatFunc(t *testing.T) {

	tests := []struct {
		in  float64
		exp float64
	}{
		{0, 1.0},
		{0.6024, 1.0},
		{0.3, 0.0},
		{0.9036, 0.0},
		{1.205, 1.0},
	}

	for _, tst := range tests {
		t.Run(fmt.Sprintf("g(%v)=%v", tst.in, tst.exp), func(nt *testing.T) {
			out := beatFunc(tst.in)
			assert.InDelta(t, tst.exp, out, eps)
		})
	}
}
