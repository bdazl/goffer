package rand

import (
	"math"
	"math/rand"
)

const (
	twoPi = math.Pi * 2.0
)

func Normal() float64 {
	Z0, Z1 := BoxMuller()
	_ = Z1
	return Z0
}

func BoxMuller() (float64, float64) {
	var (
		U1        = rand.Float64()
		U2        = rand.Float64()
		sqrt2LnU1 = -2.0 * math.Log2(U1)

		c2pU2 = math.Cos(twoPi * U2)
		s2pU2 = math.Sin(twoPi * U2)

		Z1 = sqrt2LnU1 * c2pU2
		Z2 = sqrt2LnU1 * s2pU2
	)

	return Z1, Z2
}
