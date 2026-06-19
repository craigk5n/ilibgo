package ilibgo

import (
	"testing"

	font "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)

func TestBdfFontMetadata(t *testing.T) {
	f, err := LoadFontFromData("helvB12", font.Font_helvB12())
	if err != nil {
		t.Fatalf("LoadFontFromData: %v", err)
	}
	if f.Name() != "helvB12" {
		t.Errorf("Name() = %q, want helvB12", f.Name())
	}
	if f.IsTrueType() {
		t.Error("IsTrueType() = true for a BDF font")
	}
	if f.Foundry() == "" {
		t.Error("Foundry() is empty")
	}
	if f.Family() == "" {
		t.Error("Family() is empty")
	}
	if f.Slant() == "" {
		t.Error("Slant() is empty")
	}
	if f.Weight() == "" {
		t.Error("Weight() is empty")
	}
	if f.PixelSize() <= 0 {
		t.Errorf("PixelSize() = %d, want > 0", f.PixelSize())
	}
	if f.Ascent() <= 0 {
		t.Errorf("Ascent() = %d, want > 0", f.Ascent())
	}
	if f.GlyphCount() <= 0 {
		t.Errorf("GlyphCount() = %d, want > 0", f.GlyphCount())
	}
	// FaceName / Descent / Proportional should at least not panic.
	_ = f.FaceName()
	_ = f.Descent()
	_ = f.Proportional()
}

func TestFontMetadataNilSafe(t *testing.T) {
	var f *Font
	if f.Name() != "" || f.Foundry() != "" || f.Family() != "" ||
		f.FaceName() != "" || f.Slant() != "" || f.Weight() != "" {
		t.Error("nil Font string accessors should return empty")
	}
	if f.IsTrueType() || f.Proportional() {
		t.Error("nil Font bool accessors should return false")
	}
	if f.PixelSize() != 0 || f.Ascent() != 0 || f.Descent() != 0 || f.GlyphCount() != 0 {
		t.Error("nil Font int accessors should return 0")
	}
}
