package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	title = "jacob plåttar"
	half  = 1.0 / 2.0
)

const (
	// font
	fntL = 20
)

var (
	titleL  font.Length = fntL
	legendL font.Length = fntL
	labelL  font.Length = fntL
)

var (
	unitScale4 = []plot.Tick{
		{Value: 0, Label: "0"},
		{Value: 0.25, Label: "1/4"},
		{Value: 0.5, Label: "0.5"},
		{Value: 0.75, Label: "3/4"},
		{Value: 1, Label: "1"},
	}
)

func main() {
	var (
		w    = 1000
		h    = 1000
		bins = 20

		typ = "plot"
		// Example usage
		// x   = "0 0.5 1"
		// y   = "1 2 1"
		// empty string means default values
		x = ""
		y = ""
	)
	flag.IntVar(&w, "w", w, "width")
	flag.IntVar(&h, "h", h, "height")
	flag.StringVar(&typ, "type", typ, "type of plot")
	flag.IntVar(&bins, "bins", bins, "(hist) number of bins")
	flag.StringVar(&x, "x", x, "X data")
	flag.StringVar(&y, "y", y, "Y data")
	flag.Parse()

	if y == "" {
		y = "-1 1 0 1 -1"
	}
	Y := parseValues(y)

	// If x is empty X is empty
	X := parseValues(x)
	if len(X) == 0 {
		X = uniform(len(Y))
	}

	plt := createPlot(typ, bins, X, Y)

	outpath := fmt.Sprintf("out/%v.png", typ)
	fmt.Printf("output to: %v\n", outpath)

	err := plt.Save(vg.Length(w), vg.Length(h), outpath)
	panicOn(err)
}

func parseValues(s string) plotter.Values {
	vstrs := strings.Split(s, " ")
	if len(vstrs) <= 1 {
		return plotter.Values{}
	}

	out := make(plotter.Values, len(vstrs))
	for i, v := range vstrs {
		f, err := strconv.ParseFloat(v, 64)
		panicOn(err)

		out[i] = f
	}

	return out
}

func createPlot(
	typ string,
	bins int,
	x, y plotter.Values) *plot.Plot {

	p := plot.New()

	p.Title.Text = title

	p.Title.TextStyle.Font.Size = titleL
	p.Legend.TextStyle.Font.Size = legendL

	p.X.Label.TextStyle.Font.Size = labelL
	p.X.Tick.Label.Font.Size = labelL
	p.Y.Label.TextStyle.Font.Size = labelL
	p.Y.Tick.Label.Font.Size = labelL

	p.X.Tick.Marker = plot.ConstantTicks(unitScale4)
	p.Y.Tick.Marker = plot.ConstantTicks(unitScale4)

	pts := genSpline(x, y)

	line, err := plotter.NewLine(pts)
	panicOn(err)

	//scatter, err := plotter.NewScatter(pts)
	//panicOn(err)

	p.Add(line) //, scatter)

	/*switch typ {
	case "bar":
		bar := plotter.NewScatter(
	case "hist", "histogram":
		hist, err := plotter.NewHist(y, bins)
		panicOn(err)

		p.Add(hist)
	}*/
	return p
}

func createXYs(X, Y plotter.Values) plotter.XYs {
	var (
		lenx, leny = len(X), len(Y)
	)
	if lenx != leny {
		panic("not same count")
	}

	var pts plotter.XYs
	for i, x := range X {
		y := Y[i]

		pts = append(pts, plotter.XY{X: x, Y: y})
	}
	return pts
}

func plotToFloat(v plotter.Values) []float64 {
	return []float64(v)
}

func uniform(lenx int) []float64 {
	var (
		xmax = float64(lenx - 1)
		out  = make([]float64, lenx)
	)

	for i := 0; i < lenx; i++ {
		t := float64(i) / xmax
		out[i] = t
	}
	return out
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
