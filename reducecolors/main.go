// Command reducecolors renders a before/after demonstration of ReduceColors
// (median-cut color quantization): a smooth multi-hue gradient next to the same
// image reduced to a small palette. It is the source of the "reducecolors"
// gallery image in the README.
//
//	go run ./reducecolors -outfile reducecolors.png -colors 8
package main

import (
	"flag"
	"log"
	"math"

	"github.com/craigk5n/ilibgo"
	font "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)

func mustColor(name string) ilibgo.Color {
	c, err := ilibgo.AllocNamedColor(name)
	if err != nil {
		log.Fatalf("color %q: %v", name, err)
	}
	return c
}

// gradient fills img with a smooth, many-colored field so quantization produces
// visible posterization bands.
func gradient(img *ilibgo.Image) {
	w, h := ilibgo.ImageWidth(img), ilibgo.ImageHeight(img)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := 128 + 127*math.Sin(float64(x)*0.045)
			g := 128 + 127*math.Sin(float64(y)*0.045)
			b := 128 + 127*math.Sin(float64(x+y)*0.045)
			img.SetPixel(x, y, clamp(r), clamp(g), clamp(b))
		}
	}
}

func clamp(v float64) int {
	i := int(v + 0.5)
	if i < 0 {
		return 0
	}
	if i > 255 {
		return 255
	}
	return i
}

// distinctColors counts the unique RGB triples in img.
func distinctColors(img *ilibgo.Image) int {
	seen := map[uint32]struct{}{}
	w, h := ilibgo.ImageWidth(img), ilibgo.ImageHeight(img)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := img.GetPixel(x, y)
			seen[uint32(r)<<16|uint32(g)<<8|uint32(b)] = struct{}{}
		}
	}
	return len(seen)
}

func main() {
	outfile := flag.String("outfile", "reducecolors.png", "output PNG file")
	colors := flag.Int("colors", 8, "number of colors to reduce to")
	flag.Parse()

	const panelW, panelH, gap, margin = 280, 280, 24, 16

	// Build the source gradient.
	src := ilibgo.CreateImage(panelW, panelH)
	gradient(src)

	// Reduce a copy so the original is preserved (also demonstrates
	// DuplicateImage).
	reduced := src.DuplicateImage()
	if err := reduced.ReduceColors(*colors); err != nil {
		log.Fatalf("ReduceColors: %v", err)
	}

	srcN := distinctColors(src)
	redN := distinctColors(reduced)

	// Compose both panels onto a labeled canvas.
	canvasW := margin*2 + panelW*2 + gap
	canvasH := margin + panelH + 36
	canvas := ilibgo.CreateImageWithBackground(canvasW, canvasH, mustColor("white"))

	gc := ilibgo.CreateGraphicsContext()
	canvas.CopyImage(src, gc, 0, 0, panelW, panelH, margin, margin)
	canvas.CopyImage(reduced, gc, 0, 0, panelW, panelH, margin+panelW+gap, margin)

	// Thin frames around each panel.
	ilibgo.SetForeground(&gc, mustColor("dimgray"))
	canvas.DrawRectangle(gc, margin, margin, panelW, panelH)
	canvas.DrawRectangle(gc, margin+panelW+gap, margin, panelW, panelH)

	// Labels.
	f, err := ilibgo.LoadFontFromData("helvB12", font.Font_helvB12())
	if err != nil {
		log.Fatalf("font: %v", err)
	}
	ilibgo.SetForeground(&gc, mustColor("black"))
	ilibgo.SetFont(&gc, f)
	canvas.DrawString(gc, margin+8, margin+panelH+24,
		labelf("original", srcN))
	canvas.DrawString(gc, margin+panelW+gap+8, margin+panelH+24,
		labelf("ReduceColors", redN))

	if err := ilibgo.SaveImageFile(*outfile, canvas, ilibgo.FormatPNG); err != nil {
		log.Fatalf("save %s: %v", *outfile, err)
	}
	log.Printf("wrote %s (%d -> %d colors)", *outfile, srcN, redN)
}

func labelf(name string, n int) string {
	plural := "colors"
	if n == 1 {
		plural = "color"
	}
	return name + " — " + itoa(n) + " " + plural
}

// itoa avoids pulling in strconv for one call.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
