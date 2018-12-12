package main

import (
	"fmt"
	"os"

	"github.com/tdewolff/canvas"
)

func main() {
	svgFile, err := os.Create("example.svg")
	if err != nil {
		panic(err)
	}
	defer svgFile.Close()

	svg := canvas.NewSVG(svgFile)
	svg.AddFontFile("DejaVuSerif", canvas.Regular, "DejaVuSerif.ttf")
	defer svg.Close()

	font, err := svg.Font("DejaVuSerif")
	if err != nil {
		panic(err)
	}

	face := font.Face(12)
	fmt.Println(face.LineHeight(), face.Ascent(), face.Descent())

	plot := New("")
	plot.Add(NewLine([]float64{0, 10, 20, 30, 40, 50, 60}, []float64{15, 25, 40, 30, 10, 5, 5}))
	plot.Draw(svg, font, 80.0, 50.0)
}
