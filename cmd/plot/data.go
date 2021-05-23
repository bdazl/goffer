package main

import (
	"github.com/cnkei/gospline"
	"gonum.org/v1/plot/plotter"
)

func genSpline(x, y plotter.Values) plotter.XYs {
	var (
		//Y = []float64{2, 1, 1, 1, 1, 1}
		//Y = []float64{1, 1, 2, 2, half, half}
		//X = []float64{0.0, 0.1, 0.4, 0.41, 0.7, 1.0}

		X = []float64(x)
		Y = []float64(y)

		spline = gospline.NewCubicSpline(X, Y)

		yf = spline.Range(0.0, 1.0, 1/2000.)
		xf = uniform(len(yf))
	)

	ox := plotter.Values(xf)
	oy := plotter.Values(yf)
	return createXYs(ox, oy)
}
