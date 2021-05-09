package scenes

// Made 2021-05-08

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"os"
	"path"

	_ "image/jpeg"
	_ "image/png"

	"github.com/HexHacks/goffer/pkg/bezier"
	"github.com/HexHacks/goffer/pkg/global"
	jimage "github.com/HexHacks/goffer/pkg/image"
	"github.com/HexHacks/goffer/pkg/image/mask"
	"github.com/HexHacks/goffer/pkg/math/float"

	"github.com/lucasb-eyer/go-colorful"
)

const (
	// Real value
	// bpm = 100.0
	bpm         = 100.0
	tempoFreq   = bpm / 60.0
	tempoPeriod = 60.0 / bpm

	cutoutCnt    = 20
	bezierPoints = 10
)

var (
	W  = global.Width
	H  = global.Height
	CX = W / 2
	CY = H / 2
	//(C  = image.Point{CX, CY}
	Dur = global.Total

	blue        = color.RGBA{0, 0, 255, 255}
	red         = color.RGBA{220, 10, 10, 255}
	uniformBlue = &image.Uniform{blue}

	cutoutR = image.Rect(0, 0, 100, 100)

	twoPi = math.Pi * 2.0
)

type Djanl struct {
	palette []colorful.Color
	refImgs []refImage
	strokes []stroke
}

type ImageSub interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}

func resetGlobals() {
	W = global.Width
	H = global.Height
	CX = W / 2
	CY = H / 2

	Dur = global.Total
	uniformBlue = &image.Uniform{blue}
}

type refImage struct {
	img     ImageSub
	cutouts []image.Image
}

func newRefImg(img ImageSub) refImage {
	var (
		mx = cutoutR.Max.X
		my = cutoutR.Max.Y
	)
	cutouts := make([]image.Image, cutoutCnt)
	for i := range cutouts {
		max := img.Bounds().Max
		r0, r1 := rand.Intn(max.X-mx), rand.Intn(max.Y-my)

		sub := image.Rect(r0, r1, r0+mx, r1+my)
		cutouts[i] = img.SubImage(sub)
	}
	return refImage{
		img:     img,
		cutouts: cutouts,
	}
}

type stroke struct {
	brush
	curve bezier.Curve
}

func newStroke(img image.Image, pts []complex128) stroke {
	return stroke{
		brush: newBrush(img),
		curve: bezier.New(pts...),
	}
}

func (s *stroke) Draw(dst draw.Image, t float64) {
	ptc := s.curve.Point(t)
	pt := image.Point{int(real(ptc)), int(imag(ptc))}
	s.brush.Draw(dst, pt)
}

type brush struct {
	img      image.Image
	mask     mask.Circle
	defMaskP image.Point
}

func newBrush(img image.Image) brush {
	cp := cutoutR.Max.Div(2)
	return brush{
		img:      img,
		mask:     mask.Circle{P: cp, R: cp.X},
		defMaskP: cp,
	}
}

func (b *brush) Draw(onto draw.Image, dp image.Point) {
	drawFullSrcMask(onto, b.img, &b.mask, dp)
}

// INIT ---------------
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
		dj.strokes[i] = newStroke(ref, pts)
	}
}

func randPts(n int) []complex128 {
	out := make([]complex128, n)
	for i := 0; i < n; i++ {
		s := float64(i) / float64(n-1)
		t := s * Dur

		inp := math.Mod(t, tempoPeriod)

		// Zero when not on beat
		// One when on beat
		f := beatFunc(inp)            // [0, 1)
		zeroTFS := f * float64(W) / 4 // [0, W/4)

		pt := image.Point{CX, CY}

		out[i] = randComplexPTwoCircles(pt, zeroTFS, float64(CX))
	}
	return out
}

func beatFunc(t float64) float64 {
	const (
		T = tempoPeriod

		a = 1
		b = 0
		c = 18

		m = T

		cmin = 18
		cmax = 200
	)
	g := func(x float64) float64 {
		return float.Gaussian(x, a, b, c)
	}

	m0 := math.Mod(t, m)
	m1 := math.Mod(-t, m)
	return g(m0) + g(m1)
}

func randPtsV0(n int) []complex128 {
	out := make([]complex128, n)
	for i := 0; i < n; i++ {
		pt := image.Point{W, H}

		out[i] = randComplexPoint(pt)
	}
	return out
}

// DRAW -----------------------------------------------------------------------------

func (dj *Djanl) Frame(t float64) image.Image {
	img, _ := jimage.New()

	/* bg */
	bg := dj.palette[0]

	// Background
	// draw.Draw(img, img.Bounds(), image.Transparent, image.ZP, draw.Src)
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.ZP, draw.Src)

	dj.drawAnimV0(img, t)
	//dj.drawImageV2(img)
	//dj.drawImageV1(img)
	//dj.drawImageV0(img)

	return img
}

func (dj *Djanl) drawAnimV0(img draw.Image, tNominal float64) {
	const (
		// section length
		secL = 0.1
	)

	var (
		t    = tNominal / Dur
		tFut = 1.0 - t
	)

	// Compensation for start
	var fl float64 = 0
	if t < secL {
		fl = secL - t
	}

	// Compensation for end
	var fr float64 = 0
	if tFut < secL {
		fr = secL - tFut
	}

	ll := t - secL + fl
	lr := t + secL - fr

	scnt := 100
	for _, s := range dj.strokes {
		for i := 0; i < scnt; i++ {
			ti := float64(i) / float64(scnt-1)
			curveT := (lr-ll)*ti + ll

			// radius
			oR := s.brush.defMaskP.X / 2
			s.mask.R = oR * i / (scnt - 1)

			s.Draw(img, curveT)
		}
	}
}

func (dj *Djanl) drawImageV2(img draw.Image) {
	// Ormar som växer
	scnt := 100
	for _, s := range dj.strokes {
		for i := 0; i < scnt; i++ {
			ti := float64(i) / float64(scnt-1)

			// Radius
			oR := s.brush.defMaskP.X / 2
			s.mask.R = oR * i / (scnt - 1)
			s.Draw(img, ti)
		}
	}
}

func (dj *Djanl) drawImageV1(img draw.Image) {
	// Långa feta ormar
	scnt := 100
	for _, s := range dj.strokes {
		for i := 0; i < scnt; i++ {
			ti := float64(i) / float64(scnt-1)
			s.Draw(img, ti)
		}
	}
}

func (dj *Djanl) drawImageV0(img draw.Image) {
	// Cirklar ifyllda slumpmässiga portioner av styckade referensbilder
	cp := cutoutR.Max.Div(2)
	mask := &mask.Circle{P: cp, R: cp.X}
	for i := 0; i < dj.refImgCount(); i++ {
		ref := dj.randRefImg()

		// drawFullSrc(img, ref, randPoint(img))
		drawFullSrcMask(img, ref, mask, randPoint(img.Bounds().Max))
	}
}

// ---------------- DRAW

//  RefImg - -----------------------
func (dj *Djanl) randRefImg() image.Image {
	r0, r1 := rand.Int(), rand.Int()

	ref := dj.refImgs[r0%len(dj.refImgs)]
	return ref.cutouts[r1%len(ref.cutouts)]
}

func (dj *Djanl) refImgCount() int {
	return len(dj.refImgs) * cutoutCnt
}

// ---------------------- RefImg

// RANDOM -------------------
func randPoint(max image.Point) image.Point {
	r0, r1 := rand.Int(), rand.Int()
	return image.Point{
		r0 % max.X,
		r1 % max.Y,
	}
}

func randComplexPoint(max image.Point) complex128 {
	r0, r1 := rand.Float64(), rand.Float64()
	return complex(
		r0*float64(max.X),
		r1*float64(max.Y),
	)
}

func randComplexPTwoCircles(cnt image.Point, minR, maxR float64) complex128 {
	a := randI(0, twoPi)
	r := randI(minR, maxR)

	return complex(
		math.Cos(a)*r,
		math.Sin(a)*r,
	)
}

func randI(a, b float64) float64 {
	r := rand.Float64()
	atb := b - a
	return a + (atb * r)
}

func drawFullSrc(dst draw.Image, src image.Image, dp image.Point) {
	// dp = destination point
	srcR := src.Bounds()
	dstR := srcR.Sub(srcR.Min).Add(dp)

	// Draw a red rectangle
	draw.Draw(dst, dstR, src, srcR.Min, draw.Src)
}

func drawFullSrcMask(
	dst draw.Image,
	src image.Image,
	mask image.Image,
	dp image.Point) {
	// dp = destination point
	srcR := src.Bounds()
	dstR := srcR.Sub(srcR.Min).Add(dp)

	// Draw a red rectangle
	draw.DrawMask(dst, dstR, src, srcR.Min, mask, image.ZP, draw.Over)
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

func PT(x, y int) image.Point {
	return image.Point{x, y}
}

func PTs(n int) image.Point {
	return image.Point{n, n}
}
