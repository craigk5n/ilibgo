// Create a PNG file that displays the font for ASCII values up to 255.
// Usage:
//    displayfont [options] -infile fontfile -outfile output.png
//
//    where options are:
//    -hex	display ascii char numbers in hex
//    -dec	display ascii char numbers in decimal (default)
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/craigk5n/ilibgo"
)

//  12-Aug-2022 Craig Knudsen craig@k5n.us
//		Converted from C to go
//  19-Jul-1999	Added -png option
//		Craig Knudsen	cknudsen@radix.net
//  12-Apr-1999	Created
//		Craig Knudsen	cknudsen@radix.net
func main() {
	outfile := ""
	infile := ""
	useHex := false
	imageFormat := ilibgo.FormatPNG

	flag.StringVar(&outfile, "outfile", "", "output PNG filename")
	flag.StringVar(&infile, "infile", "", "input BDF font filename")
	flag.BoolVar(&useHex, "hex", false, "Display values as hex instead of integer")
	flag.Parse()

	if len(infile) == 0 {
		fmt.Printf("You must provide an input file using -infile\n")
		os.Exit(1)
	}

	base := filepath.Base(infile)
	base = strings.Replace(base, ".bdf", "", 1)
	font, err := ilibgo.LoadFontFromFile(infile, base)
	if err != nil {
		fmt.Printf("error loading font: %v\n", err)
		os.Exit(1)
	}
	if len(outfile) == 0 {
		outfile = fmt.Sprintf("%s.png", base)
	}
	// load the small font used for labeling
	smallfont, err := ilibgo.LoadFontFromData("smallfont", IFont_courR10())
	if err != nil {
		fmt.Printf("error loading small font: %v\n", err)
		os.Exit(1)
	}

	// create output
	gc := ilibgo.CreateGraphicsContext()
	_, fontH, _ := ilibgo.TextDimensions(gc, font, "X")
	cellWidth := 40
	cellHeight := 40
	if fontH > 24 {
		cellHeight = 40 + (fontH - 24)
	}
	width := 16*cellWidth + 10
	height := 16*cellHeight + 10
	black, _ := ilibgo.AllocNamedColor("black")
	image := ilibgo.CreateImageWithBackground(width, height, black)
	white, _ := ilibgo.AllocNamedColor("white")
	grey, _ := ilibgo.AllocNamedColor("grey")
	navy, _ := ilibgo.AllocNamedColor("navy")
	ilibgo.SetForeground(&gc, white)
	ilibgo.FillRectangle(image, gc, 0, 0, width, height)
	ilibgo.SetForeground(&gc, black)

	x := 5
	y := 0
	for loop := 0; loop < 256; loop++ {
		if loop%16 == 0 && loop > 0 {
			y += cellHeight
			x = 5
		} else if loop > 0 {
			x += cellWidth
		}
		ilibgo.DrawRectangle(image, gc, x, y, cellWidth, cellHeight)
		var temp string
		if useHex {
			temp = fmt.Sprintf("%02X", loop)
		} else {
			temp = fmt.Sprintf("%d", loop)
		}
		ilibgo.SetFont(&gc, smallfont)
		w, h, _ := ilibgo.TextDimensions(gc, smallfont, temp)
		subx := x + (cellWidth-w)/2
		suby := y + h + 2
		ilibgo.SetForeground(&gc, navy)
		ilibgo.FillRectangle(image, gc, x+2, y+2, cellWidth-3, h)
		ilibgo.SetForeground(&gc, white)
		ilibgo.DrawString(image, gc, subx, suby-1, temp)

		ilibgo.SetFont(&gc, font)
		ilibgo.SetForeground(&gc, black)
		temp = fmt.Sprintf("%c", loop)
		w, _, _ = ilibgo.TextDimensions(gc, font, temp)
		subx = x + (cellWidth-w)/2
		suby = y + cellHeight - 6
		// draw a baseline
		ilibgo.SetForeground(&gc, grey)
		ilibgo.DrawLine(image, gc, x+1, suby, x+cellWidth-1, suby)
		// draw the letter
		ilibgo.SetForeground(&gc, black)
		ilibgo.DrawString(image, gc, subx, suby, temp)
	}

	// Write PNG output file.
	fp, err := os.Create(outfile)
	if err != nil {
		fmt.Printf("error opening file %s: %v\n", outfile, err)
	}

	ilibgo.WriteImageFile(fp, image, imageFormat)
	fp.Close()
}
