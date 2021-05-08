package scenes

import (
	"image"
	"image/draw"

	"github.com/HexHacks/goffer/pkg/animation"
	jimage "github.com/HexHacks/goffer/pkg/image"
	jcmplx "github.com/HexHacks/goffer/pkg/math/cmplx"
	"github.com/HexHacks/goffer/pkg/math/fourier"
	"github.com/HexHacks/goffer/pkg/palette"

	"github.com/llgcode/draw2d/draw2dimg"
)

type EpiSmooth struct {
	pts []complex128

	anim *animation.Animation
	img  image.Image
}

func (e *EpiSmooth) Init() {
	const (
		curveCnt = 3
		ptCnt    = 200
		pts      = float64(ptCnt)
		cCnt     = 21
	)

	ca := jcmplx.RandomSlice(cCnt, 0.3)
	cb := jcmplx.RandomSlice(cCnt, 0.2)

	cat := make([]complex128, cCnt)
	copy(cat, ca)

	cdelta := make([]complex128, cCnt)
	copy(cdelta, cb)
	jcmplx.Sub(cdelta, ca)
	jcmplx.Scale(cdelta, complex(1.0/(curveCnt-1), 0))

	for cc := 0; cc < curveCnt; cc++ {
		for i := 0; i < ptCnt; i++ {
			t := float64(i) / (pts - 1.0)
			out := fourier.P(t, cat)

			e.pts = append(e.pts, out)
		}
		jcmplx.Add(cat, cdelta)
	}

	common := func() *draw2dimg.GraphicContext {
		img, gc := jimage.New()
		draw.Draw(img,
			img.Bounds(),
			&image.Uniform{palette.Palette[0]},
			image.ZP, draw.Src)

		e.img = img
		return gc
	}

	e.anim = animation.New([]animation.Animator{
		func(t float64) {
			const (
				tail = 20
			)
			gc := common()
			a := 0
			cnt := int(t * float64(len(e.pts)))
			if cnt > tail {
				inc := cnt - tail
				a = inc
				cnt -= inc
			}
			jimage.DrawLines(gc, e.pts[a:], cnt)
		},
	},
		[]float64{1.0},
	)

}

func (e *EpiSmooth) Frame(t float64) image.Image {
	e.anim.Frame(t)
	return e.img
}
