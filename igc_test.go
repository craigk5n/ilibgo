package ilibgo

import "testing"

func TestCreateGraphicsContextDefaults(t *testing.T) {
	gc := CreateGraphicsContext()
	if gc.foreground.color.R != 0 || gc.foreground.color.G != 0 || gc.foreground.color.B != 0 {
		t.Errorf("default foreground = %v, want black", gc.foreground.color)
	}
	if gc.background.color.R != 255 || gc.background.color.G != 255 || gc.background.color.B != 255 {
		t.Errorf("default background = %v, want white", gc.background.color)
	}
	if gc.lineWidth != 1 {
		t.Errorf("default lineWidth = %d, want 1", gc.lineWidth)
	}
	if gc.lineStyle != LineSolid {
		t.Errorf("default lineStyle = %d, want LineSolid", gc.lineStyle)
	}
	if gc.textStyle != TextNormal {
		t.Errorf("default textStyle = %d, want TextNormal", gc.textStyle)
	}
}

func TestGraphicsContextSetters(t *testing.T) {
	gc := CreateGraphicsContext()

	red, _ := AllocNamedColor("red")
	blue, _ := AllocNamedColor("blue")
	SetForeground(&gc, red)
	SetBackground(&gc, blue)
	if gc.foreground.color.R != 255 || gc.background.color.B != 255 {
		t.Error("SetForeground/SetBackground did not apply")
	}

	SetLineStyle(&gc, LineOnOffDash)
	if gc.lineStyle != LineOnOffDash {
		t.Error("SetLineStyle did not apply")
	}
	SetTextStyle(&gc, TextShadowed)
	if gc.textStyle != TextShadowed {
		t.Error("SetTextStyle did not apply")
	}

	f := mustFont(t)
	SetFont(&gc, f)
	if gc.font != f {
		t.Error("SetFont did not apply")
	}
}

func TestSetLineWidthClamps(t *testing.T) {
	gc := CreateGraphicsContext()
	SetLineWidth(&gc, 2)
	if gc.lineWidth != 2 {
		t.Errorf("SetLineWidth(2) = %d, want 2", gc.lineWidth)
	}
	SetLineWidth(&gc, 99) // clamped to 3
	if gc.lineWidth != 3 {
		t.Errorf("SetLineWidth(99) = %d, want 3 (clamped)", gc.lineWidth)
	}
}

func TestGetFontSize(t *testing.T) {
	if _, err := GetFontSize(nil); err == nil {
		t.Error("GetFontSize(nil) expected error")
	}
	if _, err := GetFontSize(&Font{}); err == nil {
		t.Error("GetFontSize(empty font) expected error")
	}
	size, err := GetFontSize(mustFont(t))
	if err != nil {
		t.Fatalf("GetFontSize: %v", err)
	}
	if size != 34 {
		t.Errorf("GetFontSize(helvR24) = %d, want 34", size)
	}
}
