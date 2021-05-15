package djanl

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"path"

	"github.com/HexHacks/goffer/pkg/global"
	"github.com/HexHacks/goffer/pkg/math/float"
	"github.com/lucasb-eyer/go-colorful"
)

func (dj *Djanl) Init() {
	resetGlobals()

	palette, err := colorful.HappyPalette(4)
	panicOn(err)

	dj.palette = palette

	dj.initRefImages()
	dj.initStrokes()
}

func resetGlobals() {
	W = global.Width
	H = global.Height
	CX = W / 2
	CY = H / 2

	Dur = global.Total
	//bezierPoints = int(math.Ceil(Dur)) * 2
}

func (dj *Djanl) initRefImages() {
	imgs := loadInputImages()
	dj.refImgs = make([]refImage, len(imgs))
	for i, img := range imgs {
		dj.refImgs[i] = newRefImg(img)
	}
}

func (dj *Djanl) initStrokes() {
	var (
		count = dj.refImgCount()
	)
	dj.strokes = make([]stroke, count)
	for i := 0; i < count; i++ {
		ref := dj.randRefImg()
		pts := randPts(bezierPoints)
		dj.strokes[i] = newStroke(ref, pts)
	}
}

func randPts(n int) []complex128 {
	var (
		start = randI(0, twoPi)
		circs = []float64{
			randI(10, 50),
			randI(70, 150),
			randI(100, 500),
			randI(500, 800),
			randI(810, 950),
		}
	)

	prevcirc := rand.Int() % len(circs)
	out := make([]complex128, n)
	for i := 0; i < n; i++ {
		s := float64(i) / float64(n-1)

		t := s * Dur // [0, Dur)

		// Zero when not on beat
		// One when on beat
		f := beatFunc(t) // [0, 1)

		nextcirc := rand.Int() % len(circs)
		c0, c1 := circs[prevcirc], circs[nextcirc]
		cl := c1 - c0

		a := start + randI(-0.2, 0.2) + s*twoPi //randI(math.Pi, twoPi)
		r := c0 + f*cl                          //randI(0, 1000)

		r = float.Clamp(r, 10.0, 950.0)
		cnt := image.Point{CX, CY}

		out[i] = cmplxCircle(cnt, a, r)

		//w4 := float64(W) * 2.0 / 3.0 // * float64(i%2)
		prevcirc = nextcirc

		//out[i] = randComplexPTwoCircles(pt, zeroTFS, zeroTFS+20)
		//fmt.Printf("s: %v, t: %v, f: %v\nout: %v\n", s, t, f, out[i])
	}
	return out
}

func randPtsV0(n int) []complex128 {
	out := make([]complex128, n)
	for i := 0; i < n; i++ {
		pt := image.Point{W, H}

		out[i] = randComplexPoint(pt)
	}
	return out
}

func loadInputImages() []ImageSub {
	filenames := []string{
		"spejs3.png",
		"spejs5.png",
		// "bad_marketing.png",
		"heap_dream.png",
	}

	out := make([]ImageSub, len(filenames))
	for i, fn := range filenames {
		reader, err := os.Open(path.Join("assets", fn))
		panicOn(err)

		img, typ, err := image.Decode(reader)
		sub, ok := img.(ImageSub)
		if !ok {
			fmt.Println("file:", fn, "type:", typ)
			panic("could not convert to sub image")
		}
		out[i] = sub
	}

	return out
}
