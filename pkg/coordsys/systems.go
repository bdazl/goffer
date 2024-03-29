package coordsys

import (
	"github.com/bdazl/goffer/pkg/global"
	"gonum.org/v1/gonum/spatial/r2"
)

// Conversion functions

// Img: Image coordinates [(0, 0), (w, h)]
// Unit: Unit Domain, where:
//   1. origo is in the center of the image
//   2. [(-1, -1), (1, 1)]
func ImgToUnit(pos r2.Vec) r2.Vec {
	return r2.Vec{
		X: 2.0*pos.X/global.W - 1.0,
		Y: -2.0*pos.Y/global.H + 1.0,
	}
}

// Inverse of ImgToCdom (so go from a unit domain
func UnitToImg(unit r2.Vec) r2.Vec {
	w, h := global.W, global.H
	return r2.Vec{
		X: (w*unit.X + w) / 2.0,
		Y: (h - h*unit.Y) / 2.0,
	}
}

func ImgToUnitC(c complex128) complex128 {
	return complex(
		2.0*real(c)/global.W-1.0,
		-2.0*imag(c)/global.H+1.0,
	)
}

// Inverse of ImgToCdom (so go from a unit domain
func UnitToImgC(c complex128) complex128 {
	x, y := real(c), imag(c)
	w, h := global.W, global.H
	return complex(
		(w*x+w)/2.0,
		(h-h*y)/2.0,
	)
}
