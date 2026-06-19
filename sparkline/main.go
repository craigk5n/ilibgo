// sparkline renders a small inline trend graph (a "sparkline") from a list of
// numbers to a PNG. It's a minimal example of the drawing API: just a scaled
// polyline, with optional area fill and min/max/last markers.
//
// Numbers come from the command-line arguments, or from stdin (separated by
// spaces, commas, or newlines) when no arguments are given.
//
// Usage:
//
//	sparkline [options] n1 n2 n3 ...
//
//	  -out string     output PNG file (default "sparkline.png")
//	  -h int          image height in pixels (default 32)
//	  -step int       horizontal pixels between points (default 6)
//	  -color string   line color name (default "steelblue")
//	  -fill           shade the area under the line
//	  -markers        dot the min (red), max (green) and last (black) points
//
// Examples:
//
//	sparkline 3 5 2 8 4 9 6 7
//	echo "1,4,2,8,5,9,3" | sparkline -fill -markers
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/craigk5n/ilibgo"
)

const margin = 3

func main() {
	out := flag.String("out", "sparkline.png", "output PNG file")
	height := flag.Int("h", 32, "image height in pixels")
	step := flag.Int("step", 6, "horizontal pixels between points")
	colorName := flag.String("color", "steelblue", "line color name")
	fill := flag.Bool("fill", false, "shade the area under the line")
	markers := flag.Bool("markers", false, "dot the min (red), max (green) and last (black) points")
	flag.Parse()

	values, err := readNumbers(flag.Args())
	if err != nil {
		fail(err)
	}
	if len(values) < 2 {
		fail(fmt.Errorf("need at least 2 numbers, got %d", len(values)))
	}
	if *height < 8 || *step < 1 {
		fail(fmt.Errorf("height must be >= 8 and step >= 1"))
	}

	img, err := render(values, *height, *step, *colorName, *fill, *markers)
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
	fmt.Printf("wrote %s (%dx%d, %d points)\n", *out, ilibgo.ImageWidth(img), *height, len(values))
}

func render(values []float64, height, step int, colorName string, fill, markers bool) (*ilibgo.Image, error) {
	lineColor, err := ilibgo.AllocNamedColor(colorName)
	if err != nil {
		return nil, err
	}

	width := (len(values)-1)*step + 2*margin
	white, _ := ilibgo.AllocNamedColor("white")
	img := ilibgo.CreateImageWithBackground(width, height, white)
	gc := ilibgo.CreateGraphicsContext()

	min, max := values[0], values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	// Map a value to a y pixel (top = high value).
	plotH := height - 2*margin
	yOf := func(v float64) int {
		if max == min {
			return height / 2
		}
		return margin + int(float64(plotH)*(max-v)/(max-min))
	}
	xOf := func(i int) int { return margin + i*step }

	// Optional area fill under the line.
	if fill {
		fillColor, _ := ilibgo.AllocColor(205, 222, 240)
		ilibgo.SetForeground(&gc, fillColor)
		pts := make([]ilibgo.Point, 0, len(values)+2)
		for i, v := range values {
			pts = append(pts, ilibgo.Pt(xOf(i), yOf(v)))
		}
		pts = append(pts, ilibgo.Pt(xOf(len(values)-1), height-margin))
		pts = append(pts, ilibgo.Pt(xOf(0), height-margin))
		img.FillPolygon(gc, pts)
	}

	// The trend line.
	ilibgo.SetForeground(&gc, lineColor)
	for i := 1; i < len(values); i++ {
		img.DrawLine(gc, xOf(i-1), yOf(values[i-1]), xOf(i), yOf(values[i]))
	}

	// Optional markers.
	if markers {
		minI, maxI := 0, 0
		for i, v := range values {
			if v < values[minI] {
				minI = i
			}
			if v > values[maxI] {
				maxI = i
			}
		}
		dot(img, &gc, "red", xOf(minI), yOf(values[minI]))
		dot(img, &gc, "green", xOf(maxI), yOf(values[maxI]))
		dot(img, &gc, "black", xOf(len(values)-1), yOf(values[len(values)-1]))
	}

	return img, nil
}

func dot(img *ilibgo.Image, gc *ilibgo.GraphicsContext, colorName string, x, y int) {
	c, _ := ilibgo.AllocNamedColor(colorName)
	ilibgo.SetForeground(gc, c)
	img.FillCircle(*gc, x, y, 2)
}

func readNumbers(args []string) ([]float64, error) {
	fields := args
	if len(fields) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		fields = strings.FieldsFunc(string(data), func(r rune) bool {
			return r == ',' || r == ' ' || r == '\n' || r == '\r' || r == '\t'
		})
	}
	values := make([]float64, 0, len(fields))
	for _, f := range fields {
		v, err := strconv.ParseFloat(strings.TrimSpace(f), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", f, err)
		}
		values = append(values, v)
	}
	return values, nil
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "sparkline: %v\n", err)
	os.Exit(1)
}
