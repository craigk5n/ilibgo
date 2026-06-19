// truetype renders text using a scalable, anti-aliased TrueType/OpenType font.
//
// With no -font, it uses the built-in Go font, so it runs out of the box.
//
// Usage:
//
//	truetype [options]
//
//	  -font string    path to a .ttf/.otf font (default: built-in Go font)
//	  -size float     font size in points (default 48)
//	  -text string    text to render (default "Hello, ilibgo!")
//	  -color string   text color name (default "black")
//	  -style string   normal, shadowed, etched-in, or etched-out (default "normal")
//	  -angle float    rotation angle in degrees, counterclockwise (default 0)
//	  -out string     output PNG file (default "truetype.png")
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/craigk5n/ilibgo"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	fontPath := flag.String("font", "", "path to a .ttf/.otf font (default: built-in Go font)")
	size := flag.Float64("size", 48, "font size in points")
	text := flag.String("text", "Hello, ilibgo!", "text to render")
	colorName := flag.String("color", "black", "text color name")
	style := flag.String("style", "normal", "normal, shadowed, etched-in, or etched-out")
	angle := flag.Float64("angle", 0, "rotation angle in degrees (counterclockwise)")
	out := flag.String("out", "truetype.png", "output PNG file")
	flag.Parse()

	font, err := loadFont(*fontPath, *size)
	if err != nil {
		fail(err)
	}
	ts, err := parseStyle(*style)
	if err != nil {
		fail(err)
	}
	fg, err := ilibgo.AllocNamedColor(*colorName)
	if err != nil {
		fail(err)
	}

	gc := ilibgo.CreateGraphicsContext()
	ilibgo.SetFont(&gc, font)
	ilibgo.SetForeground(&gc, fg)
	ilibgo.SetTextStyle(&gc, ts)

	w, h, _ := ilibgo.TextDimensions(gc, font, *text)
	margin := int(*size / 3)

	var img *ilibgo.Image
	if *angle == 0 {
		img = ilibgo.CreateImageWithBackground(w+2*margin, h+2*margin, mustColor("white"))
		// y is the text baseline; sit it within the top margin plus the ascent.
		baseline := margin + h*4/5
		img.DrawString(gc, margin, baseline, *text)
	} else {
		// Rotated text fans out from a pivot; use a square canvas big enough to
		// hold the rotation and pivot about its center.
		diag := int(math.Hypot(float64(w), float64(h)))
		side := 2*diag + 2*margin
		img = ilibgo.CreateImageWithBackground(side, side, mustColor("white"))
		img.DrawStringRotatedAngle(gc, side/2, side/2, *text, *angle)
	}

	f, err := os.Create(*out)
	if err != nil {
		fail(fmt.Errorf("create %s: %w", *out, err))
	}
	defer f.Close()
	if err := ilibgo.WriteImageFile(f, img, ilibgo.FormatPNG); err != nil {
		fail(fmt.Errorf("write %s: %w", *out, err))
	}
	fmt.Printf("wrote %s (%dx%d)\n", *out, ilibgo.ImageWidth(img), ilibgo.ImageHeight(img))
}

func loadFont(path string, size float64) (*ilibgo.Font, error) {
	if path != "" {
		return ilibgo.LoadTrueTypeFromFile(path, "user", size, 72)
	}
	return ilibgo.LoadTrueTypeFromBytes(goregular.TTF, "goregular", size, 72)
}

func parseStyle(s string) (ilibgo.TextStyle, error) {
	switch s {
	case "normal":
		return ilibgo.TextNormal, nil
	case "shadowed":
		return ilibgo.TextShadowed, nil
	case "etched-in":
		return ilibgo.TextEtchedIn, nil
	case "etched-out":
		return ilibgo.TextEtchedOut, nil
	}
	return ilibgo.TextNormal, fmt.Errorf("invalid style %q", s)
}

func mustColor(name string) ilibgo.Color {
	c, _ := ilibgo.AllocNamedColor(name)
	return c
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "truetype: %v\n", err)
	os.Exit(1)
}
