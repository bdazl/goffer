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
)

var (
	W  = global.Width
	H  = global.Height
	CX = W / 2
	CY = H / 2
	//(C  = image.Point{CX, CY}

	blue        = color.RGBA{0, 0, 255, 255}
	red         = color.RGBA{220, 10, 10, 255}
	uniformBlue = &image.Uniform{blue}

	cutoutCnt = 20
	cutoutR   = image.Rect(0, 0, 50, 50)

	twoPi = math.Pi * 2.0
)

type Djanl struct {
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
		pts := randPts(5)
		dj.strokes[i] = newStroke(ref, pts)
	}
}

func randPts(n int) []complex128 {
	out := make([]complex128, n)
	for i := 0; i < n; i++ {
		pt := image.Point{W, H}

		out[i] = randComplexPoint(pt)
	}
	return out
}

// DRAW -------------------
func (dj *Djanl) Frame(t float64) image.Image {
	img, _ := jimage.New()

	// Background: Transparent
	draw.Draw(img, img.Bounds(), image.Transparent, image.ZP, draw.Src)

	scnt := 100
	for _, s := range dj.strokes {
		for i := 0; i < scnt; i++ {
			ti := float64(i) / float64(scnt-1)
			a := ti * twoPi

			r := float64(s.mask.R) * 1.2
			/*rot*/ _ = image.Point{
				int(math.Cos(a) * r * 0.5),
				int(math.Sin(a) * r * 0.5),
			}
			// s.mask.P = s.defMaskP.Add(rot)
			s.Draw(img, ti)
		}
	}

	// Cirklar ifyllda slumpmässiga portioner av styckade referensbilder
	/*cp := cutoutR.Max.Div(2)
	mask := &mask.Circle{P: cp, R: cp.X}
	for i := 0; i < dj.refImgCount(); i++ {
		ref := dj.randRefImg()

		// drawFullSrc(img, ref, randPoint(img))
		drawFullSrcMask(img, ref, mask, randPoint(img.Bounds().Max))
	}*/

	return img
}

func (dj *Djanl) randRefImg() image.Image {
	r0, r1 := rand.Int(), rand.Int()

	ref := dj.refImgs[r0%len(dj.refImgs)]
	return ref.cutouts[r1%len(ref.cutouts)]
}

func (dj *Djanl) refImgCount() int {
	return len(dj.refImgs) * cutoutCnt
}

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
