package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"sort"

	"github.com/craigk5n/ilibgo"
)

// ascentRatio approximates the fraction of a line's height that sits above the
// baseline. The library's Font.Ascent() reports the full line height for
// TrueType faces rather than the true baseline ascent, so we estimate it to
// position glyphs vertically centered within their layout box. Typical Latin
// faces fall in the 0.75-0.85 range.
const ascentRatio = 0.8

// boxPadding is the gap (in pixels) kept around each word's measured box so
// neighbours never touch.
const boxPadding = 3

// placement is a word resolved to a concrete position, color, and font.
type placement struct {
	text  string
	x, y  int // baseline-left draw origin
	color ilibgo.Color
	font  *ilibgo.Font
}

// fontLoader returns a *Font at a given point size, caching by rounded size so
// repeated sizes are only rasterized once.
type fontLoader func(size float64) (*ilibgo.Font, error)

// layout places words into a width x height canvas using Archimedean-spiral
// packing: the largest words are placed first near the center, and each word
// spirals outward until it finds a position whose bounding box overlaps no
// previously placed word and stays within the canvas. Words that cannot be
// placed are returned in skipped (callers may warn about them).
func layout(words []sizedWord, cfg Config, load fontLoader, rng *rand.Rand) (placed []placement, skipped []string, err error) {
	// Largest first so big words claim the center.
	sorted := make([]sizedWord, len(words))
	copy(sorted, words)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].size > sorted[j].size })

	canvas := image.Rect(0, 0, cfg.Width, cfg.Height)
	cx, cy := float64(cfg.Width)/2, float64(cfg.Height)/2
	var boxes []image.Rectangle

	for _, w := range sorted {
		font, ferr := load(w.size)
		if ferr != nil {
			return nil, nil, ferr
		}
		gc := ilibgo.CreateGraphicsContext()
		ilibgo.SetFont(&gc, font)
		tw, th, merr := ilibgo.TextDimensions(gc, font, w.text)
		if merr != nil {
			return nil, nil, fmt.Errorf("measure %q: %w", w.text, merr)
		}
		boxW, boxH := tw+2*boxPadding, th+2*boxPadding

		topLeft, ok := spiralFind(boxW, boxH, cx, cy, canvas, boxes)
		if !ok {
			skipped = append(skipped, w.text)
			continue
		}
		boxes = append(boxes, image.Rect(topLeft.X, topLeft.Y, topLeft.X+boxW, topLeft.Y+boxH))

		col, cerr := resolveColor(w.color, cfg.Palette, rng)
		if cerr != nil {
			return nil, nil, cerr
		}
		// Baseline-left origin: inset by the padding, then drop by the ascent.
		placed = append(placed, placement{
			text:  w.text,
			x:     topLeft.X + boxPadding,
			y:     topLeft.Y + boxPadding + int(math.Round(ascentRatio*float64(th))),
			color: col,
			font:  font,
		})
	}
	return placed, skipped, nil
}

// spiralFind walks an Archimedean spiral outward from (cx, cy) looking for a
// top-left position where a boxW x boxH rectangle fits inside canvas and
// overlaps none of the existing boxes. The first candidate (t=0) is the center.
func spiralFind(boxW, boxH int, cx, cy float64, canvas image.Rectangle, boxes []image.Rectangle) (image.Point, bool) {
	// radius grows ~radiusStep pixels per radian; the cap is the canvas diagonal
	// so we always exhaust reachable positions before giving up.
	const radiusStep = 4.0
	const angleStep = 0.35
	maxR := math.Hypot(float64(canvas.Dx()), float64(canvas.Dy()))

	for t := 0.0; radiusStep*t <= maxR; t += angleStep {
		r := radiusStep * t
		px := cx + r*math.Cos(t)
		py := cy + r*math.Sin(t)
		// Center the box on the spiral point.
		x := int(math.Round(px)) - boxW/2
		y := int(math.Round(py)) - boxH/2
		candidate := image.Rect(x, y, x+boxW, y+boxH)
		if !candidate.In(canvas) {
			continue
		}
		if overlapsAny(candidate, boxes) {
			continue
		}
		return image.Point{X: x, Y: y}, true
	}
	return image.Point{}, false
}

// overlapsAny reports whether r intersects any rectangle in boxes.
func overlapsAny(r image.Rectangle, boxes []image.Rectangle) bool {
	for _, b := range boxes {
		if r.Overlaps(b) {
			return true
		}
	}
	return false
}

// resolveColor picks the word's color: an explicit per-word override (name or
// hex) when set, otherwise a random pick from the palette.
func resolveColor(override string, palette []string, rng *rand.Rand) (ilibgo.Color, error) {
	if override != "" {
		return parseColor(override)
	}
	return parseColor(palette[rng.Intn(len(palette))])
}

// render draws the background and every placed word onto a fresh image.
func render(cfg Config, placed []placement) (*ilibgo.Image, error) {
	bg, err := parseColor(cfg.Background)
	if err != nil {
		return nil, fmt.Errorf("background color: %w", err)
	}
	img := ilibgo.CreateImageWithBackground(cfg.Width, cfg.Height, bg)
	for _, p := range placed {
		gc := ilibgo.CreateGraphicsContext()
		ilibgo.SetFont(&gc, p.font)
		ilibgo.SetForeground(&gc, p.color)
		img.DrawString(gc, p.x, p.y, p.text)
	}
	return img, nil
}
