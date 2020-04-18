package cmplx

// Complex functions
// w: width of image
// h: height of image
// zw: width of complex region to stretch
// zh: height of complex region to stretch
func ToImage(c complex128, w, h, zw, zh float64) (float64, float64) {
	return (real(c) + zw) * w / (2.0 * zw), (-imag(c) + zh) * h / (2.0 * zh)
}
