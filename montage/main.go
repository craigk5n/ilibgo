// montage composes several images into a single labeled grid, scaling each to
// fit its cell while preserving aspect ratio. It is a configurable
// generalization of the thumbnails tool.
//
// Usage:
//
//	montage [options] image1 image2 ...
//
//	  -out string     output PNG file (default "montage.png")
//	  -cols int       columns in the grid (default: ~square)
//	  -cw int         cell width in pixels (default 160)
//	  -ch int         cell height in pixels (default 120)
//	  -pad int        padding between cells (default 10)
//	  -bg string      background color name (default "white")
//	  -label          draw each image's filename under its cell (default true)
//	  -border         draw a border around each cell (default true)
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/craigk5n/ilibgo"
	adobe "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)

type tile struct {
	name string
	img  *ilibgo.Image
}

func main() {
	out := flag.String("out", "montage.png", "output PNG file")
	cols := flag.Int("cols", 0, "columns in the grid (0 = ~square)")
	cw := flag.Int("cw", 160, "cell width in pixels")
	ch := flag.Int("ch", 120, "cell height in pixels")
	pad := flag.Int("pad", 10, "padding between cells")
	bg := flag.String("bg", "white", "background color name")
	label := flag.Bool("label", true, "draw each image's filename under its cell")
	border := flag.Bool("border", true, "draw a border around each cell")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "montage: no input images")
		flag.Usage()
		os.Exit(2)
	}
	if *cw < 8 || *ch < 8 || *pad < 0 {
		fail(fmt.Errorf("cw/ch must be >= 8 and pad >= 0"))
	}

	// Load inputs, skipping (with a warning) any that fail to decode.
	var tiles []tile
	for _, path := range flag.Args() {
		img, err := readImage(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "montage: skipping %s: %v\n", path, err)
			continue
		}
		tiles = append(tiles, tile{name: filepath.Base(path), img: img})
	}
	if len(tiles) == 0 {
		fail(fmt.Errorf("no images could be loaded"))
	}

	ncols := *cols
	if ncols < 1 {
		ncols = int(math.Ceil(math.Sqrt(float64(len(tiles)))))
	}

	img, err := compose(tiles, ncols, *cw, *ch, *pad, *bg, *label, *border)
	if err != nil {
		fail(err)
	}

	f, err := os.Create(*out)
	if err != nil {
		fail(fmt.Errorf("create %s: %w", *out, err))
	}
	defer f.Close()
	if err := ilibgo.WriteImageFile(f, img, ilibgo.FormatPNG); err != nil {
		fail(fmt.Errorf("write %s: %w", *out, err))
	}
	fmt.Printf("wrote %s (%dx%d, %d images, %d cols)\n",
		*out, ilibgo.ImageWidth(img), ilibgo.ImageHeight(img), len(tiles), ncols)
}

func compose(tiles []tile, cols, cw, ch, pad int, bgName string, label, border bool) (*ilibgo.Image, error) {
	labelFont, err := ilibgo.LoadFontFromData("helvR10", adobe.Font_helvR10())
	if err != nil {
		return nil, err
	}
	labelH := 0
	if label {
		labelH, _ = ilibgo.GetFontSize(labelFont)
		labelH += 4
	}

	rows := (len(tiles) + cols - 1) / cols
	cellTotalH := ch + labelH
	width := pad + cols*(cw+pad)
	height := pad + rows*(cellTotalH+pad)

	bg, err := ilibgo.AllocNamedColor(bgName)
	if err != nil {
		return nil, err
	}
	img := ilibgo.CreateImageWithBackground(width, height, bg)
	gc := ilibgo.CreateGraphicsContext()
	gray, _ := ilibgo.AllocNamedColor("gray")
	black, _ := ilibgo.AllocNamedColor("black")

	for i, t := range tiles {
		col, row := i%cols, i/cols
		cellX := pad + col*(cw+pad)
		cellY := pad + row*(cellTotalH+pad)

		// Scale the source to fit the cell, preserving aspect ratio.
		sw, sh := ilibgo.ImageWidth(t.img), ilibgo.ImageHeight(t.img)
		scale := math.Min(float64(cw)/float64(sw), float64(ch)/float64(sh))
		dw, dh := int(float64(sw)*scale), int(float64(sh)*scale)
		if dw < 1 {
			dw = 1
		}
		if dh < 1 {
			dh = 1
		}
		offX, offY := (cw-dw)/2, (ch-dh)/2
		img.CopyImageScaledQuality(t.img, 0, 0, sw, sh, cellX+offX, cellY+offY, dw, dh, ilibgo.ScaleCatmullRom)

		if border {
			ilibgo.SetForeground(&gc, gray)
			img.DrawRectangle(gc, cellX, cellY, cw, ch)
		}
		if label {
			ilibgo.SetFont(&gc, labelFont)
			ilibgo.SetForeground(&gc, black)
			name := elide(gc, labelFont, t.name, cw)
			lw, _, _ := ilibgo.TextDimensions(gc, labelFont, name)
			img.DrawString(gc, cellX+(cw-lw)/2, cellY+ch+labelH-2, name)
		}
	}

	return img, nil
}

// elide trims a label with a trailing ellipsis until it fits within maxWidth.
func elide(gc ilibgo.GraphicsContext, font *ilibgo.Font, s string, maxWidth int) string {
	if w, _, _ := ilibgo.TextDimensions(gc, font, s); w <= maxWidth {
		return s
	}
	for len(s) > 1 {
		s = s[:len(s)-1]
		if w, _, _ := ilibgo.TextDimensions(gc, font, s+"..."); w <= maxWidth {
			return s + "..."
		}
	}
	return s
}

func readImage(path string) (*ilibgo.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ilibgo.ReadImageFile(f)
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "montage: %v\n", err)
	os.Exit(1)
}
