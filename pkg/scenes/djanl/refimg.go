package djanl

import (
	"image"
	"math/rand"
)

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

func (dj *Djanl) randRefImg() image.Image {
	r0, r1 := rand.Int(), rand.Int()

	ref := dj.refImgs[r0%len(dj.refImgs)]
	return ref.cutouts[r1%len(ref.cutouts)]
}

func (dj *Djanl) refImgCount() int {
	return len(dj.refImgs) * cutoutCnt
}
