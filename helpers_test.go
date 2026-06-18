package ilibgo

import (
	"testing"

	font "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)

// newWhiteImage returns a w x h image filled white.
func newWhiteImage(t *testing.T, w, h int) *Image {
	t.Helper()
	white, err := AllocNamedColor("white")
	if err != nil {
		t.Fatalf("AllocNamedColor(white): %v", err)
	}
	return CreateImageWithBackground(w, h, white)
}

// isWhite reports whether the pixel at (x, y) is pure white.
func isWhite(img *Image, x, y int) bool {
	c := GetPoint(img, x, y)
	return c.color.R == 255 && c.color.G == 255 && c.color.B == 255
}

// isSet reports whether the pixel at (x, y) differs from white (i.e. drawn).
func isSet(img *Image, x, y int) bool {
	return !isWhite(img, x, y)
}

// countSet counts the non-white pixels in the image.
func countSet(img *Image) int {
	n := 0
	for y := 0; y < img.height; y++ {
		for x := 0; x < img.width; x++ {
			if isSet(img, x, y) {
				n++
			}
		}
	}
	return n
}

// mustFont loads the bundled helvR24 BDF font or fails the test.
func mustFont(t *testing.T) *Font {
	t.Helper()
	f, err := LoadFontFromData("helvR24", font.Font_helvR24())
	if err != nil {
		t.Fatalf("LoadFontFromData(helvR24): %v", err)
	}
	return f
}
