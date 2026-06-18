package ilibgo

import "testing"

func textGC(t *testing.T) (GraphicsContext, *Font) {
	t.Helper()
	f := mustFont(t)
	gc := CreateGraphicsContext()
	SetFont(&gc, f)
	black, _ := AllocNamedColor("black")
	SetForeground(&gc, black)
	return gc, f
}

func TestTextDimensions(t *testing.T) {
	gc, f := textGC(t)

	w, h, err := TextDimensions(gc, f, "A")
	if err != nil {
		t.Fatalf("TextDimensions: %v", err)
	}
	if w <= 0 {
		t.Errorf("width of \"A\" = %d, want > 0", w)
	}
	if h != 34 {
		t.Errorf("height of \"A\" = %d, want 34 (pixelSize)", h)
	}

	wA, _ := TextWidth(gc, f, "A")
	wAB, _ := TextWidth(gc, f, "AB")
	if wAB <= wA {
		t.Errorf("width(AB)=%d should exceed width(A)=%d", wAB, wA)
	}

	hMulti, _ := TextHeight(gc, f, "A\nB")
	if hMulti <= h {
		t.Errorf("multiline height %d should exceed single-line %d", hMulti, h)
	}

	// Tab handling must not panic.
	if _, _, err := TextDimensions(gc, f, "A\tB"); err != nil {
		t.Errorf("TextDimensions with tab: %v", err)
	}
}

func TestTextDimensionsNilFont(t *testing.T) {
	gc := CreateGraphicsContext()
	if _, _, err := TextDimensions(gc, nil, "x"); err == nil {
		t.Error("TextDimensions(nil font) expected error")
	}
}

func TestDrawString(t *testing.T) {
	gc, _ := textGC(t)
	img := newWhiteImage(t, 200, 60)
	DrawString(img, gc, 10, 45, "Hello")
	if countSet(img) == 0 {
		t.Error("DrawString drew nothing")
	}
}

func TestDrawStringRotatedDirections(t *testing.T) {
	gc, _ := textGC(t)
	dirs := map[string]TextDirection{
		"left-to-right": TextLeftToRight,
		"top-to-bottom": TextTopToBottom,
		"bottom-to-top": TextBottomToTop,
	}
	for name, dir := range dirs {
		t.Run(name, func(t *testing.T) {
			img := newWhiteImage(t, 120, 120)
			DrawStringRotated(img, gc, 60, 60, "Hi", dir)
			if countSet(img) == 0 {
				t.Errorf("DrawStringRotated(%s) drew nothing", name)
			}
		})
	}
}

func TestDrawStringRotatedAngle(t *testing.T) {
	gc, _ := textGC(t)
	img := newWhiteImage(t, 200, 200)
	// Regression guard for the empty-substring fix: this must render pixels.
	DrawStringRotatedAngle(img, gc, 40, 120, "Angle", 30.0)
	if countSet(img) == 0 {
		t.Error("DrawStringRotatedAngle drew nothing")
	}
}

func TestTextStyles(t *testing.T) {
	styles := map[string]TextStyle{
		"normal":     TextNormal,
		"etched-in":  TextEtchedIn,
		"etched-out": TextEtchedOut,
		"shadowed":   TextShadowed,
	}
	for name, style := range styles {
		t.Run(name, func(t *testing.T) {
			gc, _ := textGC(t)
			SetTextStyle(&gc, style)
			img := newWhiteImage(t, 200, 60)
			DrawString(img, gc, 10, 45, "Style")
			if countSet(img) == 0 {
				t.Errorf("text style %s produced no pixels", name)
			}
		})
	}
}
