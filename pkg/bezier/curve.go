package bezier

type point struct {
	Point, Control complex128
}

// Curve implements Bezier curve calculation according to the algorithm of Robert D. Miller.
//
// Graphics Gems 5, 'Quick and Simple Bézier Curve Drawing', pages 206-209.
type Curve []point

// NewCurve returns a Curve initialized with the control points in cp.
func New(cp ...complex128) Curve {
	if len(cp) == 0 {
		return nil
	}
	c := make(Curve, len(cp))
	for i, p := range cp {
		c[i].Point = p
	}

	var w float64
	for i, p := range c {
		switch i {
		case 0:
			w = 1
		case 1:
			w = float64(len(c)) - 1
		default:
			w *= float64(len(c)-i) / float64(i)
		}
		c[i].Control = complex(real(p.Point)*w, imag(p.Point)*w)
	}

	return c
}

// Point returns the point at t along the curve, where 0 ≤ t ≤ 1.
func (c Curve) Point(t float64) complex128 {
	c[0].Point = c[0].Control
	u := t
	for i, p := range c[1:] {
		c[i+1].Point = complex(
			real(p.Control)*float64(u),
			imag(p.Control)*float64(u),
		)
		u *= t
	}

	var (
		t1 = 1 - t
		tt = t1
	)
	p := c[len(c)-1].Point
	for i := len(c) - 2; i >= 0; i-- {
		ttf := float64(tt)
		p += complex(real(c[i].Point)*ttf, imag(c[i].Point)*ttf)
		tt *= t1
	}

	return p
}

// Curve returns a slice of complex128, p, filled with points along the Bézier curve described by c.
// If the length of p is less than 2, the curve points are undefined. The length of p is not
// altered by the call.
func (c Curve) Curve(p []complex128) []complex128 {
	for i, nf := 0, float64(len(p)-1); i < len(p); i++ {
		p[i] = c.Point(float64(i) / nf)
	}
	return p
}
