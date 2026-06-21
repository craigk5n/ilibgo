package ilibgo

import "testing"

// Tests for the features added to reach parity with the C ilib library:
// explicit pixel access, AllocColorAlpha, curves, FillEllipse, ArcProperties,
// image metadata/duplication, color reduction, blend modes, and anti-aliasing.

func newWhite(t *testing.T, w, h int) *Image {
	t.Helper()
	white, _ := AllocNamedColor("white")
	return CreateImageWithBackground(w, h, white)
}

func TestAllocColorAlpha(t *testing.T) {
	c, err := AllocColorAlpha(10, 20, 30, 40)
	if err != nil {
		t.Fatalf("AllocColorAlpha: %v", err)
	}
	if c.color.R != 10 || c.color.G != 20 || c.color.B != 30 || c.color.A != 40 {
		t.Errorf("got %+v, want {10,20,30,40}", c.color)
	}
}

func TestPixelAccess(t *testing.T) {
	img := newWhite(t, 4, 4)
	if err := img.SetPixel(1, 2, 11, 22, 33); err != nil {
		t.Fatalf("SetPixel: %v", err)
	}
	r, g, b, err := img.GetPixel(1, 2)
	if err != nil || r != 11 || g != 22 || b != 33 {
		t.Errorf("GetPixel = (%d,%d,%d) err=%v, want (11,22,33)", r, g, b, err)
	}
	if err := img.SetPixelAlpha(0, 0, 1, 2, 3, 4); err != nil {
		t.Fatalf("SetPixelAlpha: %v", err)
	}
	r, g, b, a, err := img.GetPixelAlpha(0, 0)
	if err != nil || r != 1 || g != 2 || b != 3 || a != 4 {
		t.Errorf("GetPixelAlpha = (%d,%d,%d,%d) err=%v, want (1,2,3,4)", r, g, b, a, err)
	}
}

func TestPixelAccessErrors(t *testing.T) {
	img := newWhite(t, 4, 4)
	if err := img.SetPixel(-1, 0, 0, 0, 0); err == nil {
		t.Error("SetPixel out of bounds should error")
	}
	if err := img.SetPixel(0, 0, 256, 0, 0); err == nil {
		t.Error("SetPixel out-of-range component should error")
	}
	if _, _, _, err := img.GetPixel(0, 99); err == nil {
		t.Error("GetPixel out of bounds should error")
	}
}

func TestDuplicateImage(t *testing.T) {
	img := newWhite(t, 4, 4)
	img.SetComment("hello")
	red, _ := AllocNamedColor("red")
	img.SetTransparent(red)
	img.SetPixel(1, 1, 9, 9, 9)

	dup := img.DuplicateImage()
	if dup.width != 4 || dup.height != 4 {
		t.Fatalf("dup dims = %dx%d", dup.width, dup.height)
	}
	if dup.GetComment() != "hello" {
		t.Errorf("dup comment = %q", dup.GetComment())
	}
	if _, ok := dup.GetTransparent(); !ok {
		t.Error("dup should have transparent color")
	}
	r, g, b, _ := dup.GetPixel(1, 1)
	if r != 9 || g != 9 || b != 9 {
		t.Errorf("dup pixel = (%d,%d,%d)", r, g, b)
	}
	// Mutating the duplicate must not touch the original (deep copy).
	dup.SetPixel(1, 1, 0, 0, 0)
	r, _, _, _ = img.GetPixel(1, 1)
	if r != 9 {
		t.Error("mutating dup changed original")
	}
}

func TestCommentAndTransparent(t *testing.T) {
	img := newWhite(t, 2, 2)
	if img.GetComment() != "" {
		t.Error("new image should have empty comment")
	}
	if _, ok := img.GetTransparent(); ok {
		t.Error("new image should have no transparent color")
	}
	img.SetComment("note")
	if img.GetComment() != "note" {
		t.Errorf("comment = %q", img.GetComment())
	}
	blue, _ := AllocNamedColor("blue")
	img.SetTransparent(blue)
	c, ok := img.GetTransparent()
	if !ok || c.color.B != 255 {
		t.Errorf("transparent = %+v ok=%v", c.color, ok)
	}
}

func TestArcProperties(t *testing.T) {
	// 0..360 sweep about (100,100) radius 50. The midpoint angle is 180
	// (transformed), which lands at x = 100 + 50*cos(180-deg-transformed)...
	// Just sanity check the start point at angle 0.
	p := ArcProperties(100, 100, 50, 50, 0, 360)
	// At a1=0 -> transformed 360 -> cos=1, sin=0 -> (150,100).
	if p.A1X != 150 || p.A1Y != 100 {
		t.Errorf("A1 = (%d,%d), want (150,100)", p.A1X, p.A1Y)
	}
}

func TestDrawBezierValidation(t *testing.T) {
	img := newWhite(t, 20, 20)
	gc := CreateGraphicsContext()
	if err := img.DrawBezier(gc, []Point{Pt(0, 0), Pt(1, 1)}); err == nil {
		t.Error("DrawBezier with 2 points should error")
	}
	pts := []Point{Pt(0, 0), Pt(5, 10), Pt(15, 10), Pt(19, 0)}
	if err := img.DrawBezier(gc, pts); err != nil {
		t.Fatalf("DrawBezier: %v", err)
	}
}

func TestDrawSplineValidation(t *testing.T) {
	img := newWhite(t, 20, 20)
	gc := CreateGraphicsContext()
	if err := img.DrawSpline(gc, []Point{Pt(0, 0)}); err == nil {
		t.Error("DrawSpline with 1 point should error")
	}
	if err := img.DrawSpline(gc, []Point{Pt(0, 0), Pt(10, 10), Pt(19, 0)}); err != nil {
		t.Fatalf("DrawSpline: %v", err)
	}
}

// drawsSomething reports whether any pixel differs from the white background.
func drawsSomething(img *Image) bool {
	pix := img.data.Pix
	for i := 0; i+3 < len(pix); i += 4 {
		if pix[i] != 255 || pix[i+1] != 255 || pix[i+2] != 255 {
			return true
		}
	}
	return false
}

func TestCurvesDrawPixels(t *testing.T) {
	img := newWhite(t, 40, 40)
	gc := CreateGraphicsContext()
	black, _ := AllocNamedColor("black")
	SetForeground(&gc, black)
	img.DrawBezier(gc, []Point{Pt(2, 20), Pt(12, 2), Pt(28, 38), Pt(38, 20)})
	if !drawsSomething(img) {
		t.Error("DrawBezier drew nothing")
	}
}

func TestFillEllipse(t *testing.T) {
	img := newWhite(t, 40, 40)
	gc := CreateGraphicsContext()
	black, _ := AllocNamedColor("black")
	SetForeground(&gc, black)
	img.FillEllipse(gc, 20, 20, 15, 8)
	r, _, _, _ := img.GetPixel(20, 20)
	if r != 0 {
		t.Errorf("center of filled ellipse should be black, got R=%d", r)
	}
}

func TestReduceColors(t *testing.T) {
	// Build a 16x16 image with a smooth gradient (many distinct colors).
	img := CreateImage(16, 16)
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			img.SetPixel(x, y, x*16, y*16, (x+y)*8)
		}
	}
	if img.qWithinLimit(8) {
		t.Fatal("gradient should exceed 8 colors")
	}
	if err := img.ReduceColors(8); err != nil {
		t.Fatalf("ReduceColors: %v", err)
	}
	if !img.qWithinLimit(8) {
		t.Error("after ReduceColors(8) image should have <= 8 colors")
	}
}

func TestReduceColorsFewColorsUnchanged(t *testing.T) {
	img := CreateImage(8, 8)
	red, _ := AllocNamedColor("red")
	gc := CreateGraphicsContext()
	SetForeground(&gc, red)
	img.FillRectangle(gc, 0, 0, 8, 8)
	before := make([]byte, len(img.data.Pix))
	copy(before, img.data.Pix)
	if err := img.ReduceColors(4); err != nil {
		t.Fatalf("ReduceColors: %v", err)
	}
	for i := range before {
		if before[i] != img.data.Pix[i] {
			t.Fatal("solid image changed by ReduceColors")
		}
	}
}

func TestBlendOver(t *testing.T) {
	img := newWhite(t, 2, 2)
	gc := CreateGraphicsContext()
	// 50% red over white.
	half, _ := AllocColorAlpha(255, 0, 0, 128)
	SetForeground(&gc, half)
	SetBlendMode(&gc, BlendOver)
	img.SetPoint(gc, 0, 0)
	r, g, b, a := func() (int, int, int, int) {
		r, g, b, a, _ := img.GetPixelAlpha(0, 0)
		return r, g, b, a
	}()
	if r != 255 || a != 255 {
		t.Errorf("blended pixel R=%d A=%d, want R=255 A=255", r, a)
	}
	if g < 110 || g > 140 || b < 110 || b > 140 {
		t.Errorf("blended pixel G=%d B=%d, want ~127 (50%% red over white)", g, b)
	}
}

func TestBlendReplaceDefault(t *testing.T) {
	img := newWhite(t, 2, 2)
	gc := CreateGraphicsContext()
	half, _ := AllocColorAlpha(255, 0, 0, 128)
	SetForeground(&gc, half)
	// Default blend mode is replace: the pixel becomes the raw foreground.
	img.SetPoint(gc, 0, 0)
	r, g, b, a, _ := img.GetPixelAlpha(0, 0)
	if r != 255 || g != 0 || b != 0 || a != 128 {
		t.Errorf("replace pixel = (%d,%d,%d,%d), want (255,0,0,128)", r, g, b, a)
	}
}

// hasIntermediateGray reports whether any pixel is a non-pure gray (a sign of
// anti-aliased coverage of black over white).
func hasIntermediateGray(img *Image) bool {
	pix := img.data.Pix
	for i := 0; i+3 < len(pix); i += 4 {
		v := pix[i]
		if pix[i+1] == v && pix[i+2] == v && v > 0 && v < 255 {
			return true
		}
	}
	return false
}

func TestAntiAliasedLine(t *testing.T) {
	black, _ := AllocNamedColor("black")

	// Non-AA diagonal: only pure black/white pixels.
	plain := newWhite(t, 32, 32)
	gc := CreateGraphicsContext()
	SetForeground(&gc, black)
	plain.DrawLine(gc, 2, 4, 30, 26)
	if hasIntermediateGray(plain) {
		t.Error("non-AA line produced gray pixels")
	}

	// AA diagonal: should produce intermediate gray coverage.
	aa := newWhite(t, 32, 32)
	gca := CreateGraphicsContext()
	SetForeground(&gca, black)
	SetAntiAlias(&gca, true)
	aa.DrawLine(gca, 2, 4, 30, 26)
	if !hasIntermediateGray(aa) {
		t.Error("AA line produced no gray pixels")
	}
}

func TestAntiAliasedFillCircle(t *testing.T) {
	black, _ := AllocNamedColor("black")
	img := newWhite(t, 40, 40)
	gc := CreateGraphicsContext()
	SetForeground(&gc, black)
	SetAntiAlias(&gc, true)
	img.FillCircle(gc, 20, 20, 15)
	// Center is fully covered.
	if r, _, _, _ := img.GetPixel(20, 20); r != 0 {
		t.Errorf("AA filled circle center R=%d, want 0", r)
	}
	// Edge should have partially-covered gray pixels.
	if !hasIntermediateGray(img) {
		t.Error("AA fill circle produced no edge gray pixels")
	}
}

func TestAntiAliasedCircleOutline(t *testing.T) {
	black, _ := AllocNamedColor("black")
	img := newWhite(t, 40, 40)
	gc := CreateGraphicsContext()
	SetForeground(&gc, black)
	SetAntiAlias(&gc, true)
	img.DrawCircle(gc, 20, 20, 15)
	if !drawsSomething(img) {
		t.Error("AA circle outline drew nothing")
	}
	if !hasIntermediateGray(img) {
		t.Error("AA circle outline produced no gray pixels")
	}
}

func TestSetters(t *testing.T) {
	gc := CreateGraphicsContext()
	if gc.blendMode != BlendReplace {
		t.Error("default blend mode should be BlendReplace")
	}
	SetBlendMode(&gc, BlendOver)
	if gc.blendMode != BlendOver {
		t.Error("SetBlendMode failed")
	}
	SetAntiAlias(&gc, true)
	if !gc.antiAlias {
		t.Error("SetAntiAlias failed")
	}
	SetAntiAliasedFont(&gc, nil)
	if !gc.antiAliasedFont {
		t.Error("SetAntiAliasedFont should set the flag")
	}
	SetFont(&gc, nil)
	if gc.antiAliasedFont {
		t.Error("SetFont should clear the antiAliasedFont flag")
	}
}
