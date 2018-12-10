package main

import (
	"os"

	"github.com/tdewolff/canvas"
)

func main() {
	fonts := canvas.NewFonts()
	fonts.Add("DejaVuSerif", canvas.Regular, "DejaVuSerif.ttf")

	svgFile, err := os.Create("example.svg")
	if err != nil {
		panic(err)
	}
	defer svgFile.Close()

	svg := canvas.NewSVG(svgFile, fonts)
	defer svg.Close()

	plot := New()
	plot.Add(NewLine([]float64{0, 10, 20, 30, 40, 50, 60}, []float64{15, 25, 40, 30, 10, 5, 5}))
	plot.Draw(svg, 80.0, 50.0)
}
