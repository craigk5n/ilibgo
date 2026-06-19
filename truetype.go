package ilibgo

import (
	"image"
	"math"
	"os"
	"strings"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/f64"
	"golang.org/x/image/math/fixed"
)

// LoadTrueTypeFromBytes parses a TrueType or OpenType font from raw bytes and
// builds a Font rasterized at the given point size and resolution. Because
// TrueType is scalable, the size is baked into the returned Font; load again at
// a different size to change it. A dpi of 0 defaults to 72 (so points == pixels).
//
// Unlike BDF fonts, a TrueType font renders anti-aliased. It supports the
// etched/shadowed text styles and arbitrary-angle rotation via
// DrawStringRotatedAngle (the glyphs are rasterized horizontally and then
// affine-transformed).
func LoadTrueTypeFromBytes(data []byte, name string, points float64, dpi float64) (*Font, error) {
	if dpi <= 0 {
		dpi = 72
	}
	sf, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}
	face, err := opentype.NewFace(sf, &opentype.FaceOptions{
		Size:    points,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	m := face.Metrics()
	height := (m.Ascent + m.Descent).Ceil()
	return &Font{name: name, face: face, height: height}, nil
}

// LoadTrueTypeFromFile loads a TrueType/OpenType font file (see
// LoadTrueTypeFromBytes).
func LoadTrueTypeFromFile(path string, name string, points float64, dpi float64) (*Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadTrueTypeFromBytes(data, name, points, dpi)
}

// drawTrueTypeGlyphs renders text at the baseline (x, y) in the graphics
// context's foreground color, anti-aliased, one line per "\n".
func (img *Image) drawTrueTypeGlyphs(gc GraphicsContext, x int, y int, text string) {
	if gc.font == nil || gc.font.face == nil {
		return
	}
	src := image.NewUniform(gc.foreground.color)
	lineH := gc.font.height
	for i, line := range strings.Split(text, "\n") {
		d := &font.Drawer{
			Dst:  img.data,
			Src:  src,
			Face: gc.font.face,
			Dot:  fixed.P(x, y+i*lineH),
		}
		d.DrawString(line)
	}
}

// drawTrueTypeGlyphsRotated renders text rotated by angle degrees about the
// baseline-left origin (x, y). It first rasterizes the text horizontally onto a
// transparent scratch image, then affine-transforms that image onto the
// destination with bilinear sampling so the anti-aliased glyphs rotate
// smoothly. The rotation matrix matches the BDF DrawStringRotatedAngle
// convention.
func (img *Image) drawTrueTypeGlyphsRotated(gc GraphicsContext, x int, y int, text string, angle float64) {
	if gc.font == nil || gc.font.face == nil {
		return
	}

	w, h, _ := gc.font.measureTrueType(text)
	if w <= 0 || h <= 0 {
		return
	}
	ascent := gc.font.face.Metrics().Ascent.Ceil()
	lineH := gc.font.height

	// Rasterize the text horizontally onto a transparent scratch image, with
	// the first line's baseline at y = ascent.
	scratch := image.NewRGBA(image.Rect(0, 0, w, h))
	src := image.NewUniform(gc.foreground.color)
	for i, line := range strings.Split(text, "\n") {
		d := &font.Drawer{
			Dst:  scratch,
			Src:  src,
			Face: gc.font.face,
			Dot:  fixed.P(0, ascent+i*lineH),
		}
		d.DrawString(line)
	}

	// Rotate about the source baseline-left origin (0, ascent) and translate so
	// that origin lands at (x, y) in the destination. Same orientation as the
	// BDF renderer: dstX =  cos*sx + sin*sy + c, dstY = -sin*sx + cos*sy + f.
	rad := angle * math.Pi / 180.0
	cos, sin := math.Cos(rad), math.Sin(rad)
	s2d := f64.Aff3{
		cos, sin, float64(x) - sin*float64(ascent),
		-sin, cos, float64(y) - cos*float64(ascent),
	}
	xdraw.BiLinear.Transform(img.data, s2d, scratch, scratch.Bounds(), xdraw.Over, nil)
}

// measureTrueType returns the pixel width and height of text in a TrueType
// font, counting one line height per "\n"-separated line.
func (f *Font) measureTrueType(text string) (width int, height int, err error) {
	lines := strings.Split(text, "\n")
	maxW := 0
	for _, line := range lines {
		w := font.MeasureString(f.face, line).Ceil()
		if w > maxW {
			maxW = w
		}
	}
	return maxW, f.height * len(lines), nil
}
