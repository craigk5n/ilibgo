// chart renders a simple bar chart to a PNG image.
//
// Data points are given as "label=value" arguments, or as "label,value" lines
// on stdin when no arguments are supplied.
//
// Usage:
//
//	chart [options] [label=value ...]
//
//	  -out string     output PNG file (default "chart.png")
//	  -title string   chart title
//	  -w int          image width (default 800)
//	  -h int          image height (default 400)
//
// Examples:
//
//	chart -title "Sales" Jan=12 Feb=25 Mar=18 Apr=30
//	printf 'A,3\nB,7\nC,2\n' | chart -title Letters
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/craigk5n/ilibgo"
	font "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)

type point struct {
	label string
	value float64
}

func main() {
	out := flag.String("out", "chart.png", "output PNG file")
	title := flag.String("title", "", "chart title")
	width := flag.Int("w", 800, "image width")
	height := flag.Int("h", 400, "image height")
	flag.Parse()

	points, err := readPoints(flag.Args())
	if err != nil {
		fail(err)
	}
	if len(points) == 0 {
		fail(fmt.Errorf("no data points; pass label=value args or label,value lines on stdin"))
	}

	img, err := render(points, *title, *width, *height)
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
	fmt.Printf("wrote %s (%d bars)\n", *out, len(points))
}

func readPoints(args []string) ([]point, error) {
	if len(args) > 0 {
		return parsePairs(args, "=")
	}
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return parsePairs(lines, ",")
}

func parsePairs(items []string, sep string) ([]point, error) {
	points := make([]point, 0, len(items))
	for _, item := range items {
		label, valStr, ok := strings.Cut(item, sep)
		if !ok {
			return nil, fmt.Errorf("invalid data point %q (want label%svalue)", item, sep)
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(valStr), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value in %q: %w", item, err)
		}
		points = append(points, point{label: strings.TrimSpace(label), value: v})
	}
	return points, nil
}

func render(points []point, title string, width, height int) (*ilibgo.Image, error) {
	white, _ := ilibgo.AllocNamedColor("white")
	img := ilibgo.CreateImageWithBackground(width, height, white)

	gc := ilibgo.CreateGraphicsContext()
	labelFont, err := ilibgo.LoadFontFromData("helvR12", font.Font_helvR12())
	if err != nil {
		return nil, err
	}
	ilibgo.SetFont(&gc, labelFont)

	black, _ := ilibgo.AllocNamedColor("black")
	barColor, _ := ilibgo.AllocNamedColor("steelblue")

	// Plot area.
	const marginL, marginR, marginB = 50, 20, 40
	marginT := 30
	if title != "" {
		marginT = 50
	}
	plotL, plotR := marginL, width-marginR
	plotT, plotB := marginT, height-marginB
	plotW, plotH := plotR-plotL, plotB-plotT

	// Max value (avoid divide-by-zero).
	maxVal := 0.0
	for _, p := range points {
		if p.value > maxVal {
			maxVal = p.value
		}
	}
	if maxVal <= 0 {
		maxVal = 1
	}

	// Axes.
	ilibgo.SetForeground(&gc, black)
	img.DrawLine(gc, plotL, plotT, plotL, plotB)
	img.DrawLine(gc, plotL, plotB, plotR, plotB)

	// Title (centered).
	if title != "" {
		ilibgo.SetFont(&gc, labelFont)
		tw, _, _ := ilibgo.TextDimensions(gc, labelFont, title)
		img.DrawString(gc, (width-tw)/2, marginT-15, title)
	}

	// Bars.
	slot := plotW / len(points)
	barW := slot * 6 / 10
	if barW < 1 {
		barW = 1
	}
	for i, p := range points {
		barH := int(p.value / maxVal * float64(plotH))
		if barH < 0 {
			barH = 0
		}
		x := plotL + i*slot + (slot-barW)/2
		y := plotB - barH

		ilibgo.SetForeground(&gc, barColor)
		img.FillRectangle(gc, x, y, barW, barH)

		ilibgo.SetForeground(&gc, black)
		valStr := strconv.FormatFloat(p.value, 'g', -1, 64)
		vw, _, _ := ilibgo.TextDimensions(gc, labelFont, valStr)
		img.DrawString(gc, x+(barW-vw)/2, y-3, valStr)

		lw, _, _ := ilibgo.TextDimensions(gc, labelFont, p.label)
		img.DrawString(gc, x+(barW-lw)/2, plotB+14, p.label)
	}

	return img, nil
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "chart: %v\n", err)
	os.Exit(1)
}
