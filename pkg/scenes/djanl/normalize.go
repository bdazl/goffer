package djanl

import "math"

func normalize(pts []complex128) ([]complex128, complex128) {
	// find max value
	max := complex(0, 0)
	for _, p := range pts {
		max = complex(
			math.Max(real(p), real(max)),
			math.Max(imag(p), imag(max)),
		)
	}

	out := make([]complex128, len(pts))
	for i := range out {
		in := pts[i]
		out[i] = complex(
			real(in)/real(max),
			imag(in)/imag(max),
		)
	}

	return out, max
}
