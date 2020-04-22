package global

import (
	gr1 "github.com/golang/geo/r1"
	gr2 "github.com/golang/geo/r2"
)

var (
	Width      = 512
	Height     = 512
	FPS        = 30
	FrameCount = FPS * 4
)

// calculated float versions
var (
	Total   = 5.
	DT      = 1.0 / 30.
	W       = 512.0
	H       = 512.0
	CX      = 512 / 2.0
	CY      = 512 / 2.0
	WinRect gr2.Rect
)

// Assume Width, Height, FPS and FrameCount are correct
func InitGlobals() {
	Total = float64(FrameCount) / float64(FPS)
	W, H = float64(W), float64(H)
	CX, CY = W/2.0, H/2.0
	WinRect = gr2.Rect{
		X: gr1.Interval{Lo: 0.0, Hi: W},
		Y: gr1.Interval{Lo: 0.0, Hi: H},
	}
}
