package ilibgo

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// A minimal but valid BDF font defining a single 8x8 glyph for 'A'.
const fixtureBDF = `STARTFONT 2.1
FONT -test-Fixture-Medium-R-Normal--8-80-75-75-P-50-ISO8859-1
SIZE 8 75 75
FONTBOUNDINGBOX 8 8 0 0
STARTPROPERTIES 5
PIXEL_SIZE 8
FONT_ASCENT 7
FONT_DESCENT 1
SLANT "R"
SPACING "P"
ENDPROPERTIES
CHARS 1
STARTCHAR A
ENCODING 65
SWIDTH 500 0
DWIDTH 8 0
BBX 8 8 0 0
BITMAP
18
24
42
42
7E
42
42
00
ENDCHAR
ENDFONT
`

func TestLoadFontFromDataBundled(t *testing.T) {
	f := mustFont(t)
	if f.font.slant != "R" { // regression guard: SLANT must populate slant, not weight
		t.Errorf("slant = %q, want \"R\"", f.font.slant)
	}
	if f.font.pixelSize != 34 {
		t.Errorf("pixelSize = %d, want 34", f.font.pixelSize)
	}
	if f.font.fontAscent != 28 {
		t.Errorf("fontAscent = %d, want 28", f.font.fontAscent)
	}
}

func TestLoadFontFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fixture.bdf")
	if err := os.WriteFile(path, []byte(fixtureBDF), 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := LoadFontFromFile(path, "fixture")
	if err != nil {
		t.Fatalf("LoadFontFromFile: %v", err)
	}
	if f.font.pixelSize != 8 {
		t.Errorf("pixelSize = %d, want 8", f.font.pixelSize)
	}
	if f.font.slant != "R" {
		t.Errorf("slant = %q, want \"R\"", f.font.slant)
	}

	a := FontBDFGetRune(f.font, 'A')
	if a == nil {
		t.Fatal("FontBDFGetRune('A') = nil")
	}
	// Regression guard for the DWIDTH fix: actualWidth comes from DWIDTH and
	// height must remain the BBX value (not be clobbered).
	if a.actualWidth != 8 {
		t.Errorf("'A' actualWidth = %d, want 8", a.actualWidth)
	}
	if a.height != 8 {
		t.Errorf("'A' height = %d, want 8 (must not be clobbered by DWIDTH)", a.height)
	}
	set := false
	for _, on := range a.data {
		if on {
			set = true
			break
		}
	}
	if !set {
		t.Error("'A' bitmap has no set pixels; BITMAP parsing failed")
	}
}

func TestLoadFontFromBytes(t *testing.T) {
	f, err := LoadFontFromBytes("fixture", []byte(fixtureBDF))
	if err != nil {
		t.Fatalf("LoadFontFromBytes: %v", err)
	}
	if f.font.pixelSize != 8 {
		t.Errorf("pixelSize = %d, want 8", f.font.pixelSize)
	}
	if a := FontBDFGetRune(f.font, 'A'); a == nil || a.actualWidth != 8 {
		t.Errorf("glyph 'A' not parsed correctly: %v", a)
	}
}

func TestLoadFontFromFileMissing(t *testing.T) {
	if _, err := LoadFontFromFile(filepath.Join(t.TempDir(), "nope.bdf"), "x"); err == nil {
		t.Error("LoadFontFromFile on missing path expected error")
	}
}

func TestLoadFontFromDataError(t *testing.T) {
	// ENDCHAR with no preceding BITMAP/STARTCHAR is a parse error.
	if _, err := LoadFontFromData("bad", []string{"ENDCHAR"}); err == nil {
		t.Error("LoadFontFromData([ENDCHAR]) expected error")
	}
}

func TestFontBDFGetChar(t *testing.T) {
	f := mustFont(t)
	if c := FontBDFGetChar(f.font, "A"); c == nil || c.actualWidth == 0 {
		t.Error("FontBDFGetChar(A) should return a glyph with width")
	}
	// Unknown multi-rune name falls back to the space glyph (never nil).
	if c := FontBDFGetChar(f.font, "no-such-glyph-name"); c == nil {
		t.Error("FontBDFGetChar(unknown) should fall back to space, not nil")
	}
}

func TestFontDescentParsed(t *testing.T) {
	// FONT_DESCENT used to mis-parse (leading space) and silently stay 0.
	f := mustFont(t)
	if f.font.fontDescent != 7 {
		t.Errorf("helvR24 fontDescent = %d, want 7", f.font.fontDescent)
	}

	ff, err := LoadFontFromData("fixture", strings.Split(fixtureBDF, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	if ff.font.fontDescent != 1 {
		t.Errorf("fixture fontDescent = %d, want 1", ff.font.fontDescent)
	}
}

func TestEncodingRouting(t *testing.T) {
	build := func(enc int) string {
		return strings.ReplaceAll(`STARTPROPERTIES 1
PIXEL_SIZE 8
ENDPROPERTIES
STARTCHAR glyph
ENCODING @
DWIDTH 8 0
BBX 8 8 0 0
BITMAP
18
24
42
7E
42
42
42
00
ENDCHAR
`, "@", strconv.Itoa(enc))
	}

	t.Run("latin1 code lands in chars array", func(t *testing.T) {
		f, err := LoadFontFromData("l1", strings.Split(build(200), "\n"))
		if err != nil {
			t.Fatal(err)
		}
		c := FontBDFGetRune(f.font, rune(200))
		if c == nil || c.actualWidth != 8 {
			t.Errorf("Latin-1 glyph (code 200) should be retrievable via chars array, got %v", c)
		}
	})

	t.Run("code >= 256 goes to otherChars without panic", func(t *testing.T) {
		for _, enc := range []int{500, 70000} {
			f, err := LoadFontFromData("hi", strings.Split(build(enc), "\n"))
			if err != nil {
				t.Fatalf("ENCODING %d: %v", enc, err)
			}
			if len(f.font.otherChars) != 1 {
				t.Errorf("ENCODING %d: otherChars = %d, want 1", enc, len(f.font.otherChars))
			}
		}
	})
}

func FuzzLoadFontFromData(f *testing.F) {
	f.Add(fixtureBDF)
	f.Add("STARTCHAR A\nENCODING 65\nBBX 1 1 0 0\nBITMAP\n8\nENDCHAR\n")
	f.Add("")
	f.Fuzz(func(t *testing.T, s string) {
		// Must never panic regardless of input.
		_, _ = LoadFontFromData("fuzz", strings.Split(s, "\n"))
	})
}
