// mandelbrot renders the Mandelbrot set to an image, coloring each pixel by
// how quickly it escapes. It is a per-pixel showcase of SetPoint/NewColor and
// a natural CPU benchmark target.
//
// Usage:
//
//	mandelbrot [options]
//
//	  -w, -h     output dimensions
//	  -iter      maximum escape iterations
//	  -cx, -cy   center of the view in the complex plane
//	  -scale     half-height of the view (smaller zooms in)
//	  -out       output file (extension selects the format)
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/craigk5n/ilibgo"
)

func main() {
	width := flag.Int("w", 800, "image width")
	height := flag.Int("h", 600, "image height")
	maxIter := flag.Int("iter", 200, "maximum escape iterations")
	cx := flag.Float64("cx", -0.5, "view center real part")
	cy := flag.Float64("cy", 0.0, "view center imaginary part")
	scale := flag.Float64("scale", 1.2, "view half-height (smaller zooms in)")
	out := flag.String("out", "mandelbrot.png", "output file")
	flag.Parse()

	if *width < 1 || *height < 1 || *maxIter < 1 || *scale <= 0 {
		fmt.Fprintln(os.Stderr, "mandelbrot: width, height, iter must be >= 1 and scale > 0")
		os.Exit(2)
	}

	img := render(*width, *height, *maxIter, *cx, *cy, *scale)

	format, err := ilibgo.FileType(*out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mandelbrot: %v\n", err)
		os.Exit(1)
	}
	if err := ilibgo.SaveImageFile(*out, img, format); err != nil {
		fmt.Fprintf(os.Stderr, "mandelbrot: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s (%dx%d, %d iterations)\n", *out, *width, *height, *maxIter)
}

func render(width, height, maxIter int, cx, cy, scale float64) *ilibgo.Image {
	img := ilibgo.CreateImage(width, height)
	gc := ilibgo.CreateGraphicsContext()

	aspect := float64(width) / float64(height)
	for py := 0; py < height; py++ {
		// Map pixel rows to the imaginary axis.
		y0 := cy + (float64(py)/float64(height-1)*2-1)*scale
		for px := 0; px < width; px++ {
			x0 := cx + (float64(px)/float64(width-1)*2-1)*scale*aspect
			n := escape(x0, y0, maxIter)
			ilibgo.SetForeground(&gc, color(n, maxIter))
			img.SetPoint(gc, px, py)
		}
	}
	return img
}

// escape returns the iteration count at which z leaves the escape radius, or
// maxIter if the point is (assumed) in the set.
func escape(x0, y0 float64, maxIter int) int {
	var x, y float64
	for n := 0; n < maxIter; n++ {
		x2, y2 := x*x, y*y
		if x2+y2 > 4 {
			return n
		}
		y = 2*x*y + y0
		x = x2 - y2 + x0
	}
	return maxIter
}

// color maps an escape count to an RGB color. Points in the set are black;
// escaping points get a smooth blue-gold gradient.
func color(n, maxIter int) ilibgo.Color {
	if n >= maxIter {
		return ilibgo.NewColor(0, 0, 0, 255)
	}
	t := float64(n) / float64(maxIter)
	r := uint8(9 * (1 - t) * t * t * t * 255)
	g := uint8(15 * (1 - t) * (1 - t) * t * t * 255)
	b := uint8(8.5 * (1 - t) * (1 - t) * (1 - t) * t * 255)
	// Guarantee escaping points are never pure black (which marks the set).
	if r == 0 && g == 0 && b == 0 {
		b = uint8(math.Max(1, t*255))
	}
	return ilibgo.NewColor(r, g, b, 255)
}
