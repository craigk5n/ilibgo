// barcode renders a Code 39 barcode for the given text. Code 39 encodes the
// digits 0-9, uppercase A-Z, and the symbols - . space $ / + %, framed by a
// '*' start/stop character. It is a pure-geometry showcase of FillRectangle
// (no font required for the bars themselves).
//
// Usage:
//
//	barcode [options] TEXT
//
//	  -unit    narrow element width in pixels (wide elements are 3x)
//	  -height  bar height in pixels
//	  -out     output file (extension selects the format)
//
// Example:
//
//	barcode -out code.png HELLO-123
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/craigk5n/ilibgo"
)

// code39 maps each supported character to its 9-element width pattern. Each
// element alternates bar, space, bar, ... starting with a bar; 'w' is a wide
// element and 'n' a narrow one. Every pattern has exactly three wide elements.
var code39 = map[rune]string{
	'0': "nnnwwnwnn", '1': "wnnwnnnnw", '2': "nnwwnnnnw", '3': "wnwwnnnnn",
	'4': "nnnwwnnnw", '5': "wnnwwnnnn", '6': "nnwwwnnnn", '7': "nnnwnnwnw",
	'8': "wnnwnnwnn", '9': "nnwwnnwnn",
	'A': "wnnnnwnnw", 'B': "nnwnnwnnw", 'C': "wnwnnwnnn", 'D': "nnnnwwnnw",
	'E': "wnnnwwnnn", 'F': "nnwnwwnnn", 'G': "nnnnnwwnw", 'H': "wnnnnwwnn",
	'I': "nnwnnwwnn", 'J': "nnnnwwwnn", 'K': "wnnnnnnww", 'L': "nnwnnnnww",
	'M': "wnwnnnnwn", 'N': "nnnnwnnww", 'O': "wnnnwnnwn", 'P': "nnwnwnnwn",
	'Q': "nnnnnnwww", 'R': "wnnnnnwwn", 'S': "nnwnnnwwn", 'T': "nnnnwnwwn",
	'U': "wwnnnnnnw", 'V': "nwwnnnnnw", 'W': "wwwnnnnnn", 'X': "nwnnwnnnw",
	'Y': "wwnnwnnnn", 'Z': "nwwnwnnnn",
	'-': "nwnnnnwnw", '.': "wwnnnnwnn", ' ': "nwwnnnwnn", '$': "nwnwnwnnn",
	'/': "nwnwnnnwn", '+': "nwnnnwnwn", '%': "nnnwnwnwn",
	'*': "nwnnwnwnn", // start/stop, not used in data
}

func main() {
	unit := flag.Int("unit", 3, "narrow element width in pixels")
	height := flag.Int("height", 100, "bar height in pixels")
	out := flag.String("out", "barcode.png", "output file")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 || *unit < 1 || *height < 1 {
		usage()
		os.Exit(2)
	}
	text := strings.ToUpper(args[0])

	if err := validate(text); err != nil {
		fail(err)
	}
	img, err := render(text, *unit, *height)
	if err != nil {
		fail(err)
	}

	format, err := ilibgo.FileType(*out)
	if err != nil {
		fail(err)
	}
	if err := ilibgo.SaveImageFile(*out, img, format); err != nil {
		fail(err)
	}
	fmt.Printf("wrote %s (%q, %dx%d)\n", *out, text, ilibgo.ImageWidth(img), ilibgo.ImageHeight(img))
}

func validate(text string) error {
	if text == "" {
		return fmt.Errorf("empty text")
	}
	for _, r := range text {
		if r == '*' {
			return fmt.Errorf("'*' is reserved as the start/stop character")
		}
		if _, ok := code39[r]; !ok {
			return fmt.Errorf("character %q is not encodable in Code 39", r)
		}
	}
	return nil
}

func render(text string, unit, height int) (*ilibgo.Image, error) {
	const quiet = 10 // quiet-zone margin in narrow units

	// Framed message: *TEXT*.
	framed := "*" + text + "*"

	// Each character is 9 elements: 6 narrow (1u) + 3 wide (3u) = 15u, plus a
	// 1u inter-character gap after all but the last character.
	const charUnits = 6*1 + 3*3
	totalUnits := len(framed)*(charUnits+1) - 1 // drop the trailing gap
	width := (totalUnits + 2*quiet) * unit

	white, _ := ilibgo.AllocNamedColor("white")
	black, _ := ilibgo.AllocNamedColor("black")
	img := ilibgo.CreateImageWithBackground(width, height, white)

	gc := ilibgo.CreateGraphicsContext()
	ilibgo.SetForeground(&gc, black)

	x := quiet * unit
	for i, r := range framed {
		pattern := code39[r]
		for e, w := range pattern {
			ew := unit
			if w == 'w' {
				ew = 3 * unit
			}
			if e%2 == 0 { // bar (black); odd elements are spaces (left white)
				img.FillRectangle(gc, x, 0, ew, height)
			}
			x += ew
		}
		if i < len(framed)-1 {
			x += unit // inter-character gap (narrow space)
		}
	}
	return img, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: barcode [-unit N] [-height N] [-out file] TEXT")
	flag.PrintDefaults()
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "barcode: %v\n", err)
	os.Exit(1)
}
