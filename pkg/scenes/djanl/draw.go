package djanl

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"

	jimage "github.com/HexHacks/goffer/pkg/image"
	"github.com/HexHacks/goffer/pkg/image/mask"
	"github.com/HexHacks/goffer/pkg/math/float"
	"github.com/HexHacks/goffer/pkg/math/spline"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/lucasb-eyer/go-colorful"
)

func (dj *Djanl) Frame(t float64) image.Image {
	img, gc := jimage.New()
	_ = gc

	dj.drawSimpleBg(img)
	//dj.drawBG(img, t)

	/*for _, s := range dj.strokes {
		if rand.Int()%2 == 0 {
			dj.dbgDrawSpline(&s.spline, gc)
		}
	}*/
	dj.drawAnimV1Spline(img, t)

	return img
}

func (dj *Djanl) drawSimpleBg(img draw.Image) {
	pal1 := dj.palette[0]

	//draw.Draw(img, img.Bounds(), image.Transparent, image.ZP, draw.Src)
	draw.Draw(img, img.Bounds(), &image.Uniform{pal1}, image.ZP, draw.Src)
}

func (dj *Djanl) drawBG(img draw.Image, t float64) {
	/* bg */
	const (
		T4 = 0.1724137931034483
		T5 = 0.20689655172413793
		T6 = 0.24137931034482757
		T7 = 0.27586206896551724

		B4, B5, B6, B7 = T4 * 0.1, T5 * 0.1, T6 * .1, T7 * 0.1
	)
	var (
		T    = t / MaxTime
		TInv = 1.0 - T
		pal1 = dj.palette[1]
	)

	// fmt.Printf("T = %v\n", T)

	// Background
	//draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.ZP, draw.Src)
	//draw.Draw(img, img.Bounds(), dj.refImgs[0].img, image.ZP, draw.Src)
	bg := dj.refImgs[0].img
	filt := &jimage.Filter{
		Img: bg,
		FilterFunc: func(x, y int, inC color.Color) color.Color {
			r, g, b, _ := inC.RGBA()
			rr, gg, bb := float64(r)/255.0, float64(g)/255.0, float64(b)/255.0

			col := colorful.Color{R: rr, G: gg, B: bb}

			blend0 := col.BlendRgb(pal1, B4)
			//blend1 := col.BlendRgb(pal1, B5)

			/*megbl := colorful.Color{
				R: col.R - blend0.R,
				G: col.G - blend0.G,
				B: col.B - blend0.B,
			}*/

			//out := col.BlendHsv(blend0, 0.5+0.5*TInv)
			//out = out.BlendHsv(blend1, 0.5+0.5*TInv)
			out := col.BlendRgb(blend0, 0.5+0.5*TInv)

			return out
		},
	}

	/*a := piFourth * T
	r := 100.0
	pt := image.Point{
		X: int(r * math.Cos(a)),
		Y: int(r * math.Sin(a)),
	}*/
	curve := dj.bgCurve.Point(T)
	pt := CToP(curve)

	// pt = image.ZP
	// inf := &jimage.Infinite{Image: filt}
	draw.Draw(img, img.Bounds(), filt, pt, draw.Src)
}

func (dj *Djanl) drawAnimV1Spline(img draw.Image, tNominal float64) {
	const (
		// section length
		LMax = 0.1
		LMin = 0.001
	)

	var (
		t    = tNominal / MaxTime
		tFut = 1.0 - t

		//freq = 3.0
		//secA = math.Sin(t*twoPi*freq)*0.5 + 0.5
		//secA = beatFunc(tNominal)
		secL = LMax * 0.3 //secA*(LMax-LMin) + LMin

		Wint  = Width / 512
		Pixel = float64(Wint)
	)

	compensation := func(thrshld float64) float64 {
		if thrshld < secL {
			return secL - thrshld
		}
		return 0.0
	}

	fl := compensation(t) // compensate for when left <= len
	fr := compensation(tFut)

	ll := t - secL + fl // pt left of t
	lr := t + secL - fr // pt right of t

	L := lr - ll

	// Colors
	lgt := colorful.HappyColor()
	h, s, v := lgt.Hsv()
	drkr := colorful.Hsv(h, s, v*0.1)

	// Brushes
	lgtb := newColBrush(lgt, Wint)
	darkb := newColBrush(drkr, Wint)

	scnt := 2000.0
	for _, s := range dj.strokes {
		/*if i != 0 { // dbg
			continue
		}*/
		pts := s.Range(ll, lr, L/scnt)
		for _, pt := range pts {
			//v := math.Sin(lr*twoPi)*
			//lightb.SetR(

			drkpt := CToP(pt)
			lgtpt := CToP(pt + complex(Pixel, Pixel/2))
			darkb.DrawColor(img, drkpt, drkr)
			lgtb.DrawColor(img, lgtpt, lgt)
		}
	}
}

func (dj *Djanl) drawAnimV0(img draw.Image, tNominal float64) {
	const (
		// section length
		secL = 0.1
	)

	var (
		t    = tNominal / MaxTime
		tFut = 1.0 - t
	)

	compensation := func(thrshld float64) float64 {
		if thrshld < secL {
			return secL - thrshld
		}
		return 0.0
	}

	fl := compensation(t) // compensate for when left <= len
	fr := compensation(tFut)

	ll := t - secL + fl // pt left of t
	lr := t + secL - fr // pt right of t

	L := lr - ll

	scnt := 200
	for _, s := range dj.strokes {
		for i := 0; i < scnt; i++ {
			ti := float64(i) / float64(scnt-1)

			curveT := float.Clamp(ll+L*ti, 0.0, 1.0)

			// radius
			maxR := W * correct * float64(s.brush.defMaskP.X) / 4.0
			fr := maxR * (ti + 0.1)
			s.SetR(int(fr))

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
	cp := image.Point{int(cutoutR), int(cutoutR)}
	mask := &mask.Circle{P: cp, R: cp.X}
	for i := 0; i < dj.refImgCount(); i++ {
		ref := dj.randRefImg()

		// drawFullSrc(img, ref, randPoint(img))
		drawFullSrcMask(img, ref, mask, randPoint(img.Bounds().Max))
	}
}

func (dj *Djanl) dbgDrawSpline(s *spline.Spline, gc *draw2dimg.GraphicContext) {
	const (
		samples = 1000.0
	)

	rcoli0 := rand.Int() % len(dj.palette)
	rcoli1 := rand.Int() % len(dj.palette)
	gc.SetFillColor(dj.palette[rcoli0])
	gc.SetStrokeColor(dj.palette[rcoli1])
	gc.SetLineWidth(3)

	pts := s.Range(0, 1, 1.0/samples)
	jimage.DrawLinesImgCoords(gc, pts, samples)
}

func (dj *Djanl) dbgDrawStroke(s *stroke, gc *draw2dimg.GraphicContext) {
	const (
		samples = 1000
	)

	var (
		fsampmax = float64(samples - 1)
	)

	gc.SetFillColor(dj.palette[2])
	gc.SetStrokeColor(dj.palette[3])
	gc.SetLineWidth(5)

	pts := make([]complex128, samples)
	for i := 0; i < samples; i++ {
		t := float64(i) / fsampmax
		pts[i] = s.curve.Point(t)
	}

	jimage.DrawLinesImgCoords(gc, pts, samples)
}

func drawFullSrc(dst draw.Image, src image.Image, dp image.Point) {
	// dp = destination point
	srcR := src.Bounds()
	dstR := srcR.Sub(srcR.Min).Add(dp)

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

	draw.DrawMask(dst, dstR, src, srcR.Min, mask, image.ZP, draw.Over)
}
