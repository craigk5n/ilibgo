// watermark overlays a tiled, semi-transparent text watermark onto an image.
//
// The watermark text is rendered to a 1-bit coverage mask and then
// alpha-blended over the source image, so it tints the picture rather than
// painting over it.
//
// Usage:
//
//	watermark -in photo.jpg [options]
//
//	  -in string       input image (required)
//	  -out string      output PNG file (default "watermark.png")
//	  -text string     watermark text (default "CONFIDENTIAL")
//	  -color string    watermark color name (default "red")
//	  -opacity float    blend strength 0..1 (default 0.30)
//	  -angle float      text angle in degrees (default 30)
//	  -tile             repeat the text across the whole image (default true)
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/craigk5n/ilibgo"
	font "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)

func main() {
	in := flag.String("in", "", "input image (required)")
	out := flag.String("out", "watermark.png", "output PNG file")
	text := flag.String("text", "CONFIDENTIAL", "watermark text")
	colorName := flag.String("color", "red", "watermark color name")
	opacity := flag.Float64("opacity", 0.30, "blend strength 0..1")
	angle := flag.Float64("angle", 30, "text angle in degrees")
	tile := flag.Bool("tile", true, "repeat the text across the whole image")
	flag.Parse()

	if *in == "" {
		fmt.Fprintln(os.Stderr, "watermark: -in is required")
		flag.Usage()
		os.Exit(2)
	}
	if *opacity < 0 || *opacity > 1 {
		fail(fmt.Errorf("opacity must be between 0 and 1"))
	}
	wmColor, err := ilibgo.AllocNamedColor(*colorName)
	if err != nil {
		fail(err)
	}

	base, err := readImage(*in)
	if err != nil {
		fail(err)
	}
	w, h := ilibgo.ImageWidth(base), ilibgo.ImageHeight(base)

	mask, err := buildMask(w, h, *text, *angle, *tile)
	if err != nil {
		fail(err)
	}
	blend(base, mask, wmColor, *opacity)

	f, err := os.Create(*out)
	if err != nil {
		fail(fmt.Errorf("create %s: %w", *out, err))
	}
	defer f.Close()
	if err := ilibgo.WriteImageFile(f, base, ilibgo.FormatPNG); err != nil {
		fail(fmt.Errorf("write %s: %w", *out, err))
	}
	fmt.Printf("wrote %s (%dx%d)\n", *out, w, h)
}

// buildMask renders the watermark text (black on white) onto a scratch image;
// any non-white pixel is "covered" by the watermark.
func buildMask(w, h int, text string, angle float64, tile bool) (*ilibgo.Image, error) {
	white, _ := ilibgo.AllocNamedColor("white")
	mask := ilibgo.CreateImageWithBackground(w, h, white)

	gc := ilibgo.CreateGraphicsContext()
	f, err := ilibgo.LoadFontFromData("helvB18", font.Font_helvB18())
	if err != nil {
		return nil, err
	}
	ilibgo.SetFont(&gc, f)
	black, _ := ilibgo.AllocNamedColor("black")
	ilibgo.SetForeground(&gc, black)

	tw, th, _ := ilibgo.TextDimensions(gc, f, text)
	if tw <= 0 || th <= 0 {
		return nil, fmt.Errorf("empty watermark text")
	}

	if !tile {
		mask.DrawStringRotatedAngle(gc, (w-tw)/2, (h+th)/2, text, angle)
		return mask, nil
	}

	stepX := tw + 60
	stepY := th + 50
	row := 0
	for y := th; y < h+stepY; y += stepY {
		offset := 0
		if row%2 == 1 {
			offset = stepX / 2 // brick-stagger alternate rows
		}
		for x := -tw + offset; x < w; x += stepX {
			mask.DrawStringRotatedAngle(gc, x, y, text, angle)
		}
		row++
	}
	return mask, nil
}

// blend tints every watermark-covered pixel of base toward wmColor by opacity.
func blend(base, mask *ilibgo.Image, wmColor ilibgo.Color, opacity float64) {
	w, h := ilibgo.ImageWidth(base), ilibgo.ImageHeight(base)
	wr, wg, wb, _ := wmColor.RGBA() // 16-bit, opaque
	gc := ilibgo.CreateGraphicsContext()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			mr, _, _, _ := mask.GetPoint(x, y).RGBA()
			if mr >= 0x8000 {
				continue // light mask pixel: not covered by the watermark
			}
			br, bg, bb, _ := base.GetPoint(x, y).RGBA()
			out := ilibgo.NewColor(
				mix(br, wr, opacity),
				mix(bg, wg, opacity),
				mix(bb, wb, opacity),
				255,
			)
			ilibgo.SetForeground(&gc, out)
			base.SetPoint(gc, x, y)
		}
	}
}

// mix blends two 16-bit channel values by t and returns an 8-bit result.
func mix(dst, src uint32, t float64) uint8 {
	d := float64(dst>>8) * (1 - t)
	s := float64(src>>8) * t
	return uint8(d + s)
}

func readImage(path string) (*ilibgo.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	img, err := ilibgo.ReadImageFile(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return img, nil
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "watermark: %v\n", err)
	os.Exit(1)
}
