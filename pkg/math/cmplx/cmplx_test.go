package cmplx

import (
	"math"
	"testing"

	jmath "github.com/HexHacks/goffer/pkg/math"
	"github.com/stretchr/testify/assert"
)

const (
	eps = 2e-10
)

func TestIntegrate(t *testing.T) {
	{
		c := []complex128{1 + 0i, 1 + 1i, 1 + 2i, 1 + 3i}
		o := Integrate(0, 5, c)

		// area of triangle
		iexp := 3.0 * 5.0 / 2.0

		assert.InDelta(t, 5.0, real(o), eps)
		assert.InDelta(t, iexp, imag(o), eps)
	}
	{
		l := 10
		c := make([]complex128, l)
		for i := range c {
			fi := float64(i)
			f := fi * jmath.Tau / float64(l-1)
			c[i] = complex(math.Cos(f), math.Sin(f))
		}
		o := Integrate(0.0, jmath.Tau, c)

		assert.InDelta(t, 0.0, real(o), eps)
		assert.InDelta(t, 0.0, imag(o), eps)
	}
}
