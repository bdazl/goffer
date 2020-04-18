package main

var (
	OutputFileType = GIF
)

// calculated
var (
	Total = 5.
	DT    = 1.0 / 30.
	W     = 512.0
	H     = 512.0
	CX    = 512 / 2.0
	CY    = 512 / 2.0
	C     = V(CX, CY)
)

func initGlobals() {
	P = Projects[ActiveProject]

	Total = float64(FrameCount) / float64(FPS)
	W, H = float64(Width), float64(Height)
	CX, CY = W/2.0, H/2.0
	C = V(CX, CY)
}
