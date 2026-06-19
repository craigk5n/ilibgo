package ilibgo

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
)

func ttFont(t *testing.T, points float64) *Font {
	t.Helper()
	f, err := LoadTrueTypeFromBytes(goregular.TTF, "goregular", points, 72)
	if err != nil {
		t.Fatalf("LoadTrueTypeFromBytes: %v", err)
	}
	return f
}

func TestTrueTypeLoadAndSize(t *testing.T) {
	f := ttFont(t, 24)
	if !f.isTrueType() {
		t.Fatal("expected a TrueType-backed font")
	}
	h, err := GetFontSize(f)
	if err != nil {
		t.Fatalf("GetFontSize: %v", err)
	}
	// A 24pt font at 72dpi should be roughly 24-34px tall (ascent+descent).
	if h < 20 || h > 50 {
		t.Errorf("font height = %d, want ~24-34", h)
	}
}

func TestTrueTypeMeasure(t *testing.T) {
	f := ttFont(t, 24)
	gc := CreateGraphicsContext()
	SetFont(&gc, f)

	wi, h, err := TextDimensions(gc, f, "i")
	if err != nil {
		t.Fatal(err)
	}
	wW, _, _ := TextDimensions(gc, f, "WWWW")
	if wW <= wi {
		t.Errorf("width(WWWW)=%d should exceed width(i)=%d", wW, wi)
	}
	if size, _ := GetFontSize(f); h != size {
		t.Errorf("single-line height = %d, want font height %d", h, size)
	}
	if _, hm, _ := TextDimensions(gc, f, "a\nb\nc"); hm <= h {
		t.Errorf("3-line height %d should exceed single-line %d", hm, h)
	}
}

func TestTrueTypeDrawAntiAliased(t *testing.T) {
	f := ttFont(t, 28)
	gc := CreateGraphicsContext()
	SetFont(&gc, f)
	black, _ := AllocNamedColor("black")
	SetForeground(&gc, black)

	img := newWhiteImage(t, 220, 60)
	img.DrawString(gc, 10, 42, "Hello")
	if countSet(img) == 0 {
		t.Fatal("DrawString drew nothing")
	}
	// TrueType rendering is anti-aliased: there should be gray edge pixels
	// (neither pure white nor pure black), unlike the 1-bit BDF path.
	gray := false
	for y := 0; y < 60 && !gray; y++ {
		for x := 0; x < 220; x++ {
			c := GetPoint(img, x, y).color
			if c.R > 0 && c.R < 255 {
				gray = true
				break
			}
		}
	}
	if !gray {
		t.Error("expected anti-aliased (gray) edge pixels from TrueType rendering")
	}
}

func TestTrueTypeStylesAndAngle(t *testing.T) {
	f := ttFont(t, 24)
	for _, style := range []TextStyle{TextNormal, TextEtchedIn, TextEtchedOut, TextShadowed} {
		gc := CreateGraphicsContext()
		SetFont(&gc, f)
		SetForeground(&gc, mustColor(t, "black"))
		SetTextStyle(&gc, style)
		img := newWhiteImage(t, 200, 60)
		img.DrawString(gc, 10, 42, "Style")
		if countSet(img) == 0 {
			t.Errorf("style %d drew nothing", style)
		}
	}

	// DrawStringRotatedAngle falls back to horizontal for TrueType (no panic).
	gc := CreateGraphicsContext()
	SetFont(&gc, f)
	SetForeground(&gc, mustColor(t, "black"))
	img := newWhiteImage(t, 200, 100)
	img.DrawStringRotatedAngle(gc, 10, 50, "Angle", 30)
	if countSet(img) == 0 {
		t.Error("DrawStringRotatedAngle (TrueType) drew nothing")
	}
}

func TestTrueTypeLoadFromFileAndDefaultDPI(t *testing.T) {
	// Default DPI (0 -> 72).
	if _, err := LoadTrueTypeFromBytes(goregular.TTF, "go", 18, 0); err != nil {
		t.Fatalf("LoadTrueTypeFromBytes with default dpi: %v", err)
	}

	// Round-trip through a file.
	path := filepath.Join(t.TempDir(), "go.ttf")
	if err := os.WriteFile(path, goregular.TTF, 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := LoadTrueTypeFromFile(path, "go", 18, 72)
	if err != nil {
		t.Fatalf("LoadTrueTypeFromFile: %v", err)
	}
	if !f.isTrueType() {
		t.Error("file-loaded font should be TrueType")
	}
}

func TestTrueTypeLoadError(t *testing.T) {
	if _, err := LoadTrueTypeFromBytes([]byte("not a font"), "bad", 12, 72); err == nil {
		t.Error("expected error parsing invalid font data")
	}
	if _, err := LoadTrueTypeFromFile("/no/such/font.ttf", "missing", 12, 72); err == nil {
		t.Error("expected error for missing file")
	}
}

func mustColor(t *testing.T, name string) Color {
	t.Helper()
	c, err := AllocNamedColor(name)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
