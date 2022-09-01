// Sample app for ilibgo package.  This app shows how to
// load and use a BDF font, draw lines, draw arcs, draw rectangles, etc.
// Usage:
//	sample [-width N] [-height N] [-text inputtext] [-out file]
package main

import (
	"flag"
	"fmt"
	"os"

	ilibgo "github.com/craigk5n/ilibgo"
	font "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)


const DefaultFontName string = "helvR24"
const DefaultFontName2 string = "helvR08"

func main() {
	width := 600
	height := 250
	text := ""
	outfile := "out.png"
	infile := ""
	copyright := "@Copyright 2022 Craig Knudsen"
	url := "https://github.com/craigk5n/ilibgo"

	flag.IntVar(&width, "width", width, "output image width")
	flag.IntVar(&height, "height", height, "output image height")
	flag.StringVar(&text, "text", "", "sample text to embed")
	flag.StringVar(&outfile, "outfile", "out.png", "output filename")
	flag.StringVar(&infile, "infile", "", "input filename")
	flag.Parse()

	if len(text) == 0 {
		text = fmt.Sprintf("Ilib v%s (%s)\n%s", ilibgo.IlibVersion, ilibgo.IlibVersionDate, url)
	}

	var img *ilibgo.Image
	if len(infile) > 0 {
		var err error
		fp, err := os.Open(infile)
		if err != nil {
			fmt.Printf("error loading input file %s: %v", infile, err)
			os.Exit(1)
		}
		img, err = ilibgo.ReadImageFile(fp)
		if err != nil {
			fmt.Printf("error loading input file %s: %v", infile, err)
			os.Exit(1)
		}
	} else {
		white, _ := ilibgo.AllocNamedColor("white")
		img = ilibgo.CreateImageWithBackground(width, height, white)
	}

	gc := ilibgo.CreateGraphicsContext()
	background, _ := ilibgo.AllocNamedColor("gray")
	ilibgo.SetBackground(&gc, background)
	topshadow, _ := ilibgo.AllocNamedColor("lightgrey")
	bottomshadow, _ := ilibgo.AllocNamedColor("darkgrey")
	textcolor, _ := ilibgo.AllocNamedColor("orange")

	// draw top shadow rectangle *
	ilibgo.SetForeground(&gc, topshadow)
	ilibgo.FillRectangle(img, gc, 0, 0, width, height)

	// draw bottom shadow rectangle
	ilibgo.SetForeground(&gc, bottomshadow)
	ilibgo.FillRectangle(img, gc, 2, 2, width-2, height-2)

	// draw background rectangle
	ilibgo.SetForeground(&gc, background)
	ilibgo.FillRectangle(img, gc, 2, 2, width-4, height-4)

	// Now the fun part: draw some text
	smallfont, _ := ilibgo.LoadFontFromData("timR12", font.Font_timR12())
	largefont, _ := ilibgo.LoadFontFromData("helvR24", font.Font_helvR24())

	ilibgo.SetFont(&gc, largefont)
	textWidth, textHeight, _ := ilibgo.TextDimensions(gc, largefont, text)
	textHeight += (textHeight / 2) // extra padding

	// draw arc
	boxHeight := 100
	arcLen := boxHeight / 2
	x := (width - textWidth) / 2
	//y := ((height - textHeight) / 2) + fontHeight
	y := height / 2
	ilibgo.SetForeground(&gc, topshadow)
	ilibgo.IFillArc(img, gc, x-2, y-2, 20, arcLen, 90, 270)
	ilibgo.IFillArc(img, gc, width-x-2, y-2, 20, arcLen, -90, 90)
	ilibgo.FillRectangle(img, gc, x-2, y-arcLen-2, width-2*x, boxHeight)
	ilibgo.SetForeground(&gc, bottomshadow)
	ilibgo.IFillArc(img, gc, x+2, y+2, 20, arcLen, 90, 270)
	ilibgo.IFillArc(img, gc, width-x+2, y+2, 20, arcLen, -90, 90)
	ilibgo.FillRectangle(img, gc, x+2, y-arcLen+2, width-2*x, boxHeight)
	ilibgo.SetForeground(&gc, background)
	ilibgo.IFillArc(img, gc, x, y, 20, arcLen, 90, 270)
	ilibgo.IFillArc(img, gc, width-x, y, 20, arcLen, -90, 90)
	ilibgo.FillRectangle(img, gc, x, y-arcLen, width-2*x+2, boxHeight)

	// draw text
	ilibgo.SetForeground(&gc, textcolor)
	ilibgo.SetTextStyle(&gc, ilibgo.TextShadowed)
	ilibgo.DrawString(img, gc, x, y, text)

	// draw "SAMPLE" from top to bottom on the left side
	ilibgo.SetTextStyle(&gc, ilibgo.TextEtchedIn)
	sampleWidth, _, _ := ilibgo.TextDimensions(gc, largefont, "SAMPLE")
	ilibgo.DrawStringRotated(img, gc, 8, (height-sampleWidth)/2+1,
		"SAMPLE", ilibgo.TextTopToBottom)

	// draw "SAMPLE" from bottom to top on the right side
	ilibgo.DrawStringRotated(img, gc,
		width-6, (height+sampleWidth)/2-1,
		"SAMPLE", ilibgo.TextBottomToTop)

	// draw copyright
	smallTextWidth, smallTextHeight, _ := ilibgo.TextDimensions(gc, smallfont, copyright)
	x = (width - smallTextWidth) / 2 // centered
	y = height - smallTextHeight     // bottom of image
	ilibgo.SetTextStyle(&gc, ilibgo.TextShadowed)
	ilibgo.SetFont(&gc, smallfont)
	ilibgo.SetForeground(&gc, textcolor)
	ilibgo.DrawString(img, gc, x, y, copyright)

	// write output image file
	fp, err := os.Create(outfile)
	if err != nil {
		fmt.Printf("error opening file %s: %v\n", outfile, err)
	}
	defer fp.Close()
	ilibgo.WriteImageFile(fp, img, ilibgo.FormatPNG)
}
