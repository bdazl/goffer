package djanl

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"path"

	"github.com/HexHacks/goffer/pkg/math/float"
	"github.com/lucasb-eyer/go-colorful"
)

type ptFunc = func(s float64) complex128

func (dj *Djanl) Init() {
	resetGlobals()

	palette, err := colorful.HappyPalette(4)
	panicOn(err)

	dj.palette = palette

	dj.initRefImages()
	dj.initStrokes()
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

		norm, max := normalize(pts)

		dj.strokes[i] = newStroke(ref, norm, max)
	}
}

func randPts(n int) []complex128 {
	var (
		start = randI(0, twoPi)
		cnt   = image.Point{CX, CY}
		radi  = []float64{
			//randI(10, 50),
			//randI(70, 150),
			//randI(100, 500),
			randI(350, 400),
			randI(400, 500),
			//
			randI(600, 800),
			randI(810, 950),
		}
		radStart = rand.Int() % len(radi)
	)

	prevR := radStart
	baseline := func(s, a, b, d float64) complex128 {
		x := start + s*twoPi

		rr, currR := radVariation(radi, prevR, s)
		prevR = currR

		r0, r1 := randI(-0.01, 0.01), randI(-0.01, 0.01)

		return lissajous(cnt, x, rr+r0, rr+r1, a, b, d)
	}

	// Lissajous parameters
	// x, A, B, a, b, δ
	funcs := []ptFunc{
		func(s float64) complex128 {
			return baseline(s, 1, 1, piHalf)
		},
		func(s float64) complex128 {
			return baseline(s, 1, 2, piHalf)
		},
		func(s float64) complex128 {
			return baseline(s, 3, 2, piHalf)
		},
	}

	l := n / len(funcs)
	variations := make([][]complex128, len(funcs))
	for i, f := range funcs {
		variations[i] = ptLoop(l, f)

		// UGLY :(
		if i == 0 {
			// Add trail from center to first variation
			c := complex(float64(CX), float64(CY))
			//fst := variations[i][0]
			cnt := l / 2
			extra := make([]complex128, cnt)
			for y := 0; y < cnt; y++ {
				//s := float64(y) / float64(cnt-1)
				//sc := complex(s, s)
				extra[y] = c //c + sc*(fst-c)
			}

			variations[i] = append(extra, variations[i]...)
		}
	}

	flat := flattenPts(variations)
	return extend(flat, 50)
}

func extend(in []complex128, extra int) []complex128 {
	var (
		lin    = len(in)
		lextra = lin + extra
	)
	out := make([]complex128, lextra)
	copy(out, in)

	for i := lin; i < lextra; i++ {
		out[i] = in[lin-1]
	}
	return out
}

func randPtsV1(n int) []complex128 {
	// Cirklar
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

		t := s * MaxTime // [0, Dur)

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

		prevcirc = nextcirc
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

func radVariation(vals []float64, prev int, x float64) (float64, int) {
	curr := rand.Int() % len(vals)
	r0, r1 := vals[prev], vals[curr]
	rr := x*(r1-r0) + r0
	return rr, curr
}

func ptLoop(n int, f func(float64) complex128) []complex128 {
	out := make([]complex128, n)
	for i := 0; i < n; i++ {
		s := float64(i) / float64(n-1)
		out[i] = f(s)
	}
	return out
}

func flattenPts(args [][]complex128) []complex128 {
	out := make([]complex128, 0)
	for _, lst := range args {
		out = append(out, lst...)
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
