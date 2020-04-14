package main

// cmd line arguments
var (
	FPS        = 30
	FrameCount = FPS * 4
	Width      = 512
	Height     = 512
)

var (
	OutputFileType = MP4
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
	Total = float64(FrameCount) / float64(FPS)
	W, H = float64(Width), float64(Height)
	CX, CY = W/2.0, H/2.0
	C = V(CX, CY)
}
