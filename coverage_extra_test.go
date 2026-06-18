package ilibgo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLineHelpers(t *testing.T) {
	vert := lineType{x1: 3, y1: 0, x2: 3, y2: 9}
	setLineSlope(&vert)
	if vert.slope != 0 {
		t.Errorf("vertical slope = %v, want 0 (sentinel)", vert.slope)
	}

	diag := lineType{x1: 0, y1: 0, x2: 10, y2: 10}
	setLineSlope(&diag)
	if got := getIntersectionXValue(diag, 5); got != 5 {
		t.Errorf("getIntersectionXValue(diag, 5) = %d, want 5", got)
	}
	if got := getIntersectionXValue(vert, 4); got != 3 {
		t.Errorf("getIntersectionXValue(vertical, 4) = %d, want 3", got)
	}
	horiz := lineType{x1: 1, y1: 5, x2: 9, y2: 5}
	if got := getIntersectionXValue(horiz, 5); got != 1 {
		t.Errorf("getIntersectionXValue(horizontal) = %d, want x1=1", got)
	}

	if !lineIncludesYValue(lineType{y1: 2, y2: 8}, 5) {
		t.Error("lineIncludesYValue should include 5 in [2,8]")
	}
	if !lineIncludesYValue(lineType{y1: 8, y2: 2}, 5) {
		t.Error("lineIncludesYValue should handle reversed endpoints")
	}
	if lineIncludesYValue(lineType{y1: 2, y2: 8}, 20) {
		t.Error("lineIncludesYValue should exclude 20 from [2,8]")
	}
}

func TestDrawLineSteepAndUpward(t *testing.T) {
	cases := []struct {
		name           string
		x1, y1, x2, y2 int
	}{
		{"steep-positive", 0, 0, 2, 9},
		{"steep-negative", 0, 9, 2, 0},
		{"vertical-upward", 5, 9, 5, 0}, // exercises the y-swap branch
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			img := newWhiteImage(t, 10, 10)
			if err := DrawLine(img, redGC(t), tc.x1, tc.y1, tc.x2, tc.y2); err != nil {
				t.Fatal(err)
			}
			if countSet(img) == 0 {
				t.Errorf("%s drew nothing", tc.name)
			}
		})
	}
}

func TestFillPolygonRectangle(t *testing.T) {
	// A rectangle has vertical and horizontal edges, exercising both special
	// cases in the polygon scanline intersection logic.
	rect := []Point{{x: 2, y: 2}, {x: 12, y: 2}, {x: 12, y: 12}, {x: 2, y: 12}}
	img := newWhiteImage(t, 16, 16)
	if err := FillPolygon(img, redGC(t), rect); err != nil {
		t.Fatal(err)
	}
	if !isSet(img, 7, 7) {
		t.Error("rectangle polygon interior (7,7) should be filled")
	}
}

func TestFillArcFullCircle(t *testing.T) {
	// a1=360, a2=0 makes the span >= 359.9 after the y-flip, hitting the
	// "draw a full circle, no center point" branch.
	img := newWhiteImage(t, 60, 60)
	if err := IFillArc(img, redGC(t), 30, 30, 20, 20, 360, 0); err != nil {
		t.Fatal(err)
	}
	if countSet(img) == 0 {
		t.Error("full-circle IFillArc drew nothing")
	}
}

// fontWithOtherChar defines a glyph stored in the non-ASCII otherChars list
// (ENCODING -1 leaves its multi-rune STARTCHAR name intact).
const fontWithOtherChar = `STARTFONT 2.1
FONTBOUNDINGBOX 8 8 0 0
STARTPROPERTIES 3
PIXEL_SIZE 8
FONT_ASCENT 7
SLANT "R"
ENDPROPERTIES
CHARS 2
STARTCHAR A
ENCODING 65
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
STARTCHAR Aacute
ENCODING -1
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
ENDFONT
`

func TestFontBDFGetCharOtherChars(t *testing.T) {
	f, err := LoadFontFromData("other", strings.Split(fontWithOtherChar, "\n"))
	if err != nil {
		t.Fatalf("LoadFontFromData: %v", err)
	}
	if len(f.font.otherChars) == 0 {
		t.Fatal("expected a glyph in otherChars")
	}
	if c := FontBDFGetChar(f.font, "Aacute"); c == nil || c.name != "Aacute" {
		t.Errorf("FontBDFGetChar(Aacute) = %v, want the otherChars glyph", c)
	}
}

func TestLoadFontFromFileParseError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "broken.bdf")
	if err := os.WriteFile(path, []byte("ENDCHAR\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFontFromFile(path, "broken"); err == nil {
		t.Error("LoadFontFromFile on malformed BDF expected error")
	}
}

func TestEtchedShadowBackgrounds(t *testing.T) {
	// Drives makeTopAndBottomShadow through its mid-range and dark branches.
	for _, bg := range []struct {
		name    string
		r, g, b uint8
	}{
		{"midgray", 100, 100, 100},
		{"neardark", 30, 30, 30},
		{"bright", 250, 250, 250},
	} {
		t.Run(bg.name, func(t *testing.T) {
			gc, _ := textGC(t)
			SetTextStyle(&gc, TextEtchedIn)
			c, _ := AllocColor(bg.r, bg.g, bg.b)
			SetBackground(&gc, c)
			img := newWhiteImage(t, 120, 60)
			DrawString(img, gc, 10, 45, "Xy")
			if countSet(img) == 0 {
				t.Errorf("etched text on %s background drew nothing", bg.name)
			}
		})
	}
}

func TestRotatedWhitespace(t *testing.T) {
	gc, _ := textGC(t)
	for _, dir := range []TextDirection{TextLeftToRight, TextTopToBottom, TextBottomToTop} {
		img := newWhiteImage(t, 160, 160)
		DrawStringRotated(img, gc, 80, 80, "A\nB\tC", dir)
		if countSet(img) == 0 {
			t.Errorf("rotated whitespace text (dir=%d) drew nothing", dir)
		}
	}
	img := newWhiteImage(t, 200, 200)
	DrawStringRotatedAngle(img, gc, 60, 120, "A\nB\tC", 20.0)
	if countSet(img) == 0 {
		t.Error("angled whitespace text drew nothing")
	}
}

func TestTextDimensionsEscapeAndTab(t *testing.T) {
	gc, f := textGC(t)
	// Escape and tab runes must be handled without error.
	if _, _, err := TextDimensions(gc, f, "A\033bad;B\tC"); err != nil {
		t.Errorf("TextDimensions with escape/tab: %v", err)
	}
}
