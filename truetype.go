package ilibgo

import (
	"image"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// LoadTrueTypeFromBytes parses a TrueType or OpenType font from raw bytes and
// builds a Font rasterized at the given point size and resolution. Because
// TrueType is scalable, the size is baked into the returned Font; load again at
// a different size to change it. A dpi of 0 defaults to 72 (so points == pixels).
//
// Unlike BDF fonts, a TrueType font renders anti-aliased and supports the
// etched/shadowed text styles, but only horizontal text (no arbitrary-angle
// rotation).
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
