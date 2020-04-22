package scenes

import (
	"image"
	"math"
	"math/rand"

	"github.com/HexHacks/goffer/pkg/coordsys"
	"github.com/HexHacks/goffer/pkg/global"
	jimage "github.com/HexHacks/goffer/pkg/image"
	jmath "github.com/HexHacks/goffer/pkg/math"
	jcmplx "github.com/HexHacks/goffer/pkg/math/cmplx"
	jr2 "github.com/HexHacks/goffer/pkg/math/r2"
	"github.com/HexHacks/goffer/pkg/palette"
	jpalette "github.com/HexHacks/goffer/pkg/palette"
	"github.com/llgcode/draw2d/draw2dimg"

	gr1 "github.com/golang/geo/r1"
	gr2 "github.com/golang/geo/r2"
	kit "github.com/llgcode/draw2d/draw2dkit"
	"github.com/lucasb-eyer/go-colorful"
	"gonum.org/v1/gonum/spatial/r2"
)

const (
	trajectory = 10
)

type dEqPt struct {
	Pos r2.Vec
	Vel r2.Vec

	AccPos []r2.Vec
}

// returned false implies remove thes point (or repurpose it)
func (pt *dEqPt) Update(t float64) bool {
	// assume velocity is not for us to decide

	// update position
	pt.Pos = pt.Pos.Add(pt.Vel)
	pp := coordsys.UnitToImg(pt.Pos)

	// add to history of positions
	pt.AccPos = append(pt.AccPos, pp)

	// Go through the trailing points and remove any point completely outside of the screen
	xi := gr1.Interval{Lo: pp.X, Hi: pp.Y}
	yi := gr1.Interval{Lo: pp.Y, Hi: pp.Y}

	aLen := len(pt.AccPos)
	cnt := min(aLen-1, trajectory)
	for i := 0; i < cnt-1; i++ {
		a, b := pt.AccPos[aLen-i-1], pt.AccPos[aLen-i-2]

		expand(&xi, a.X)
		expand(&yi, a.Y)

		expand(&xi, b.X)
		expand(&yi, b.Y)
	}
	r := gr2.Rect{X: xi, Y: yi}
	return global.WinRect.Contains(r)

}

// return false when nothing is rendered
func (pt *dEqPt) Render(img *image.RGBA, gc *draw2dimg.GraphicContext) {
	p := coordsys.UnitToImg(pt.Pos)

	// render tail
	aLen := len(pt.AccPos)
	cnt := min(aLen-1, trajectory)
	for i := 0; i < cnt-1; i++ {
		a, b := pt.AccPos[aLen-i-1], pt.AccPos[aLen-i-2]

		gc.MoveTo(a.X, a.Y)
		gc.LineTo(b.X, b.Y)
		gc.Close()
	}

	// render pt
	gc.SetFillColor(palette.Palette[5])
	kit.Circle(gc, p.X, p.Y, 5.0)
	gc.FillStroke()

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expand(i *gr1.Interval, c float64) {
	i.Lo = math.Min(i.Lo, c)
	i.Hi = math.Max(i.Hi, c)
}

// A field has a bunch of points and a derivative operator for those points
type dEqField struct {
	Pts        []dEqPt
	DtOperator func(t float64, spc r2.Vec) r2.Vec
	t          float64
}

func (f *dEqField) Update(t float64) {
	const (
		speed = 0.1
	)
	f.t = t

	for i := range f.Pts {
		f.Pts[i].Vel = f.DtOperator(t, f.Pts[i].Pos).Scale(speed * global.DT)
		mv := f.Pts[i].Update(t)
		if !mv {
			f.Pts[i].Pos = randV()
			f.Pts[i].AccPos = f.Pts[i].AccPos[:0]
		}
	}
}

func (f *dEqField) Render(img *image.RGBA, gc *draw2dimg.GraphicContext) {
	c1, _ := colorful.Hex("#fdffcc")
	c2, _ := colorful.Hex("#242a42")

	bnds := img.Bounds()
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			spc := coordsys.ImgToUnit(r2.Vec{X: float64(x), Y: float64(y)})
			grad := f.DtOperator(f.t, spc)
			sph := jr2.ToSpherical(grad)

			img.Set(x, y, c1.BlendHsv(c2, sph.Y/jmath.Tau))
		}
	}

	for _, pt := range f.Pts {
		pt.Render(img, gc)
	}
}

type DiffEq struct {
	Field dEqField
}

func NewDiffEq() *DiffEq {
	const (
		ptCount = 100
	)

	pts := make([]dEqPt, ptCount)

	for i := range pts {
		pts[i] = dEqPt{
			Pos: randV(),
		}
	}
	return &DiffEq{
		Field: dEqField{
			Pts: pts,
			DtOperator: func(t float64, spc r2.Vec) r2.Vec {
				z := complex(spc.X, spc.Y)

				w := 4 * (z + 1/z)

				return jcmplx.ToVec(w)
			},
		},
	}
}

func randV() r2.Vec {
	return r2.Vec{X: 2.0*rand.Float64() - 1.0, Y: 2.0*rand.Float64() - 1.0}
}

func (d *DiffEq) Init() {
	jpalette.Palette = jpalette.Debug
}

func (d *DiffEq) Frame(t float64) image.Image {
	img, gc := jimage.New()

	d.Update(t)
	d.Render(img, gc)

	return img
}

func (d *DiffEq) Update(t float64) {
	d.Field.Update(t)
}

func (d *DiffEq) Render(img *image.RGBA, gc *draw2dimg.GraphicContext) {
	d.Field.Render(img, gc)
}
