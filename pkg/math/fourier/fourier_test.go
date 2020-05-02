package fourier

import (
	"fmt"
	"math"
	"testing"

	jmath "github.com/HexHacks/goffer/pkg/math"
	"github.com/stretchr/testify/assert"
)

const (
	eps = 2e-4
)

func TestFourier(t *testing.T) {
	const (
		pi = math.Pi
	)

	{
		// Constant function (F(x) = 1)
		cnt := 10
		path := make([]complex128, cnt)
		for i := range path {
			path[i] = complex(1, 0)
		}

		expc := make([]complex128, 11)
		expc[5] = complex(1, 0)

		coef := Coefficients(path, 5)

		/*for i, c := range coef {
			a, b := real(expc[i]), real(c)
			fmt.Printf("%.5f - %.5f = %.5f\n", a, b, a-b)
		}*/

		for i, c := range coef {
			assert.InDelta(t, 0.0, imag(c), eps)
			assert.InDelta(t, real(expc[i]), real(c), eps)
		}

		for i := range path {
			fl := float64(i) / float64(len(path)-1)
			p := P(fl, coef)

			assert.InDelta(t, 1.0, real(p), eps)
			assert.InDelta(t, 0.0, imag(p), eps)
		}
	}
	{
		// F(x) = x
		cnt := 100
		path := make([]complex128, cnt)
		for i := range path {
			fi := float64(i)
			t := fi / float64(cnt-1)
			path[i] = complex(jmath.Tau*t-math.Pi, 0)
		}

		C := func(n int) complex128 {
			if n == 0 {
				return complex(0, 0)
			}

			fn := float64(n)
			return complex(0, math.Pow(-1.0, fn)/fn)
		}

		o := 10
		cnt = 2*o + 1
		hcnt := cnt / 2
		expc := make([]complex128, cnt)
		for i := range expc {
			n := i - hcnt
			expc[i] = C(n)
		}

		coef := Coefficients(path, o)

		for i, c := range coef {
			a, b := imag(expc[i]), imag(c)
			fmt.Printf("%.5f - %.5f = %.5f\n", a, b, a-b)
		}

		for i, c := range coef {
			assert.InDelta(t, 0.0, real(c), eps)
			assert.InDelta(t, imag(expc[i]), imag(c), 2e-2)
		}

		// Seems to converge, but endpoints are quite bad.
		// Not sure if there is anything to do here
		/*for i, ep := range path {
			fl := float64(i) / float64(len(path)-1)
			p := P(fl, coef)

			a, b := real(ep), real(p)
			fmt.Printf("%.5f - %.5f = %.5f\n", a, b, a-b)

			assert.InDelta(t, real(ep), real(p), eps)
			assert.InDelta(t, imag(ep), imag(p), eps)
		}*/
	}
}
