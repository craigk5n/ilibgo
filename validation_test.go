package ilibgo

import "testing"

// These tests exercise the input-validation error paths added so that the
// drawing primitives return meaningful errors instead of silently doing
// nothing (or panicking, in the polygon case).

func TestArcNegativeRadius(t *testing.T) {
	img := newWhiteImage(t, 20, 20)
	gc := redGC(t)
	cases := []struct {
		name string
		fn   func() error
	}{
		{"DrawArc", func() error { return img.DrawArc(gc, 10, 10, -1, 5, 0, 360) }},
		{"DrawArc r2", func() error { return img.DrawArc(gc, 10, 10, 5, -1, 0, 360) }},
		{"DrawEnclosedArc", func() error { return img.DrawEnclosedArc(gc, 10, 10, -1, 5, 0, 90) }},
		{"FillArc", func() error { return img.FillArc(gc, 10, 10, 5, -1, 0, 90) }},
		{"DrawCircle", func() error { return img.DrawCircle(gc, 10, 10, -3) }},
		{"FillCircle", func() error { return img.FillCircle(gc, 10, 10, -3) }},
		{"DrawEllipse", func() error { return img.DrawEllipse(gc, 10, 10, -3, 4) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.fn(); err == nil {
				t.Errorf("%s with negative radius: want error, got nil", tc.name)
			}
		})
	}
}

func TestArcValidRadius(t *testing.T) {
	img := newWhiteImage(t, 40, 40)
	gc := redGC(t)
	if err := img.FillCircle(gc, 20, 20, 8); err != nil {
		t.Fatalf("FillCircle valid radius: %v", err)
	}
	if !isSet(img, 20, 20) {
		t.Error("filled circle should color its center")
	}
}

func TestPolygonTooFewPoints(t *testing.T) {
	img := newWhiteImage(t, 20, 20)
	gc := redGC(t)
	if err := img.DrawPolygon(gc, []Point{Pt(1, 1)}); err == nil {
		t.Error("DrawPolygon with 1 point: want error, got nil")
	}
	if err := img.FillPolygon(gc, []Point{Pt(1, 1), Pt(5, 5)}); err == nil {
		t.Error("FillPolygon with 2 points: want error, got nil")
	}
	// Empty slice must not panic.
	if err := img.FillPolygon(gc, nil); err == nil {
		t.Error("FillPolygon with no points: want error, got nil")
	}
}

func TestAllocNamedColorError(t *testing.T) {
	if _, err := AllocNamedColor("definitely-not-a-color"); err == nil {
		t.Error("AllocNamedColor with unknown name: want error, got nil")
	}
}
