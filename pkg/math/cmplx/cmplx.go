package cmplx

import "gonum.org/v1/gonum/spatial/r2"

// Complex functions
// w: width of image
// h: height of image
// zw: width of complex region to stretch
// zh: height of complex region to stretch
func ToImage(c complex128, w, h, zw, zh float64) (float64, float64) {
	return (real(c) + zw) * w / (2.0 * zw), (-imag(c) + zh) * h / (2.0 * zh)
}

func ToVec(c complex128) r2.Vec {
	return r2.Vec{
		X: real(c),
		Y: imag(c),
	}
}
