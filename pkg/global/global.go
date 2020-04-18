package global

var (
	Width      = 512
	Height     = 512
	FPS        = 30
	FrameCount = FPS * 4
)

// calculated float versions
var (
	Total = 5.
	DT    = 1.0 / 30.
	W     = 512.0
	H     = 512.0
	CX    = 512 / 2.0
	CY    = 512 / 2.0
)

// Assume Width, Height, FPS and FrameCount are correct
func InitGlobals() {
	Total = float64(FrameCount) / float64(FPS)
	W, H = float64(W), float64(H)
	CX, CY = W/2.0, H/2.0
}
