package ilibgo

import "testing"

func redGC(t *testing.T) GraphicsContext {
	t.Helper()
	gc := CreateGraphicsContext()
	red, _ := AllocNamedColor("red")
	SetForeground(&gc, red)
	return gc
}

func TestSetGetPoint(t *testing.T) {
	img := newWhiteImage(t, 10, 10)
	gc := redGC(t)
	if err := DrawPoint(img, gc, 3, 4); err != nil {
		t.Fatalf("DrawPoint: %v", err)
	}
	c := GetPoint(img, 3, 4)
	if c.color.R != 255 || c.color.G != 0 || c.color.B != 0 {
		t.Errorf("GetPoint(3,4) = %v, want red", c.color)
	}
	if !isWhite(img, 0, 0) {
		t.Error("unrelated pixel should remain white")
	}
}

func TestDrawLineDirections(t *testing.T) {
	cases := []struct {
		name           string
		x1, y1, x2, y2 int
		wx, wy         int // a pixel known to lie on the drawn line
	}{
		{"horizontal", 0, 5, 9, 5, 5, 5},
		{"vertical", 5, 0, 5, 9, 5, 5},
		{"diagonal", 0, 0, 9, 9, 5, 5},
		{"reverse-diagonal", 9, 0, 0, 9, 4, 5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			img := newWhiteImage(t, 10, 10)
			gc := redGC(t)
			if err := DrawLine(img, gc, tc.x1, tc.y1, tc.x2, tc.y2); err != nil {
				t.Fatalf("DrawLine: %v", err)
			}
			if !isSet(img, tc.wx, tc.wy) {
				t.Errorf("%s: expected (%d,%d) to be drawn", tc.name, tc.wx, tc.wy)
			}
		})
	}
}

func TestDrawLineWidths(t *testing.T) {
	for _, w := range []int{0, 1, 2, 3} {
		img := newWhiteImage(t, 20, 20)
		gc := redGC(t)
		SetLineWidth(&gc, w)
		DrawLine(img, gc, 2, 10, 17, 10)
		if countSet(img) == 0 {
			t.Errorf("lineWidth=%d drew no pixels", w)
		}
	}
}

func TestDrawLineDashed(t *testing.T) {
	solid := newWhiteImage(t, 40, 10)
	DrawLine(solid, redGC(t), 0, 5, 39, 5)

	dashed := newWhiteImage(t, 40, 10)
	gc := redGC(t)
	SetLineStyle(&gc, LineOnOffDash)
	DrawLine(dashed, gc, 0, 5, 39, 5)

	if countSet(dashed) == 0 {
		t.Fatal("dashed line drew nothing")
	}
	if countSet(dashed) >= countSet(solid) {
		t.Errorf("dashed line (%d px) should set fewer pixels than solid (%d px)",
			countSet(dashed), countSet(solid))
	}

	// Also exercise the dashed sloped path.
	diag := newWhiteImage(t, 40, 40)
	DrawLine(diag, gc, 0, 0, 39, 39)
	if countSet(diag) == 0 {
		t.Error("dashed diagonal drew nothing")
	}
}

func TestFillRectangle(t *testing.T) {
	img := newWhiteImage(t, 10, 10)
	DrawLineThenFill := redGC(t)
	if err := FillRectangle(img, DrawLineThenFill, 2, 2, 4, 4); err != nil {
		t.Fatalf("FillRectangle: %v", err)
	}
	if !isSet(img, 3, 3) {
		t.Error("interior pixel (3,3) should be filled")
	}
	if !isWhite(img, 0, 0) {
		t.Error("(0,0) should remain white")
	}
	if !isWhite(img, 8, 8) {
		t.Error("(8,8) outside rect should remain white")
	}
}

func TestFillRectangleClipping(t *testing.T) {
	img := newWhiteImage(t, 10, 10)
	gc := redGC(t)
	// Rectangle partly past the bottom-right edge must clip without panic.
	if err := FillRectangle(img, gc, 8, 8, 5, 5); err != nil {
		t.Fatalf("FillRectangle clip: %v", err)
	}
	if !isSet(img, 9, 9) {
		t.Error("(9,9) should be filled")
	}

	// Negative origin must clip without panic.
	img2 := newWhiteImage(t, 10, 10)
	if err := FillRectangle(img2, gc, -2, -2, 4, 4); err != nil {
		t.Fatalf("FillRectangle negative: %v", err)
	}
	if !isSet(img2, 0, 0) {
		t.Error("(0,0) should be filled when rect starts off-image")
	}
}

func TestDrawRectangle(t *testing.T) {
	img := newWhiteImage(t, 10, 10)
	gc := redGC(t)
	if err := DrawRectangle(img, gc, 1, 1, 5, 5); err != nil {
		t.Fatalf("DrawRectangle: %v", err)
	}
	if !isSet(img, 3, 1) {
		t.Error("top border pixel (3,1) should be drawn")
	}
	if !isWhite(img, 3, 3) {
		t.Error("interior (3,3) should be white for an unfilled rectangle")
	}
}

func TestArcsAndCircles(t *testing.T) {
	t.Run("DrawCircle", func(t *testing.T) {
		img := newWhiteImage(t, 40, 40)
		if err := DrawCircle(img, redGC(t), 20, 20, 10); err != nil {
			t.Fatal(err)
		}
		if countSet(img) == 0 {
			t.Error("DrawCircle drew nothing")
		}
		if !isWhite(img, 20, 20) {
			t.Error("outline circle should not fill its center")
		}
	})
	t.Run("FillCircle", func(t *testing.T) {
		img := newWhiteImage(t, 40, 40)
		if err := FillCircle(img, redGC(t), 20, 20, 10); err != nil {
			t.Fatal(err)
		}
		if !isSet(img, 20, 20) {
			t.Error("filled circle should set its center pixel")
		}
	})
	t.Run("DrawArc", func(t *testing.T) {
		img := newWhiteImage(t, 40, 40)
		if err := DrawArc(img, redGC(t), 20, 20, 10, 10, 0, 90); err != nil {
			t.Fatal(err)
		}
		if countSet(img) == 0 {
			t.Error("DrawArc drew nothing")
		}
	})
	t.Run("FillArc", func(t *testing.T) {
		img := newWhiteImage(t, 40, 40)
		if err := FillArc(img, redGC(t), 20, 20, 15, 15, 0, 180); err != nil {
			t.Fatal(err)
		}
		if countSet(img) == 0 {
			t.Error("FillArc drew nothing")
		}
	})
	t.Run("EnclosedArc", func(t *testing.T) {
		img := newWhiteImage(t, 60, 60)
		if err := DrawEnclosedArc(img, redGC(t), 30, 30, 20, 20, 0, 90); err != nil {
			t.Fatal(err)
		}
		if countSet(img) == 0 {
			t.Error("DrawEnclosedArc drew nothing")
		}
	})
	t.Run("Ellipse", func(t *testing.T) {
		img := newWhiteImage(t, 60, 40)
		if err := DrawEllipse(img, redGC(t), 30, 20, 20, 10); err != nil {
			t.Fatal(err)
		}
		if countSet(img) == 0 {
			t.Error("DrawEllipse drew nothing")
		}
	})
}

// TestDeprecatedArcAliases confirms the legacy I-prefixed names still forward
// to the renamed functions.
func TestDeprecatedArcAliases(t *testing.T) {
	checks := []struct {
		name string
		draw func(*Image, GraphicsContext) error
	}{
		{"IDrawArc", func(img *Image, gc GraphicsContext) error { return IDrawArc(img, gc, 20, 20, 10, 10, 0, 90) }},
		{"IFillArc", func(img *Image, gc GraphicsContext) error { return IFillArc(img, gc, 20, 20, 15, 15, 0, 180) }},
		{"IDrawEnclosedArc", func(img *Image, gc GraphicsContext) error { return IDrawEnclosedArc(img, gc, 20, 20, 10, 10, 0, 90) }},
		{"IDrawCircle", func(img *Image, gc GraphicsContext) error { return IDrawCircle(img, gc, 20, 20, 10) }},
		{"IFillCircle", func(img *Image, gc GraphicsContext) error { return IFillCircle(img, gc, 20, 20, 10) }},
		{"IDrawEllipse", func(img *Image, gc GraphicsContext) error { return IDrawEllipse(img, gc, 20, 20, 15, 8) }},
	}
	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			img := newWhiteImage(t, 40, 40)
			if err := c.draw(img, redGC(t)); err != nil {
				t.Fatalf("%s: %v", c.name, err)
			}
			if countSet(img) == 0 {
				t.Errorf("%s alias drew nothing", c.name)
			}
		})
	}
}

func TestPolygons(t *testing.T) {
	triangle := []Point{{X: 5, Y: 2}, {X: 25, Y: 2}, {X: 15, Y: 25}}

	outline := newWhiteImage(t, 30, 30)
	if err := DrawPolygon(outline, redGC(t), triangle); err != nil {
		t.Fatalf("DrawPolygon: %v", err)
	}
	if countSet(outline) == 0 {
		t.Error("DrawPolygon drew nothing")
	}

	filled := newWhiteImage(t, 30, 30)
	if err := FillPolygon(filled, redGC(t), triangle); err != nil {
		t.Fatalf("FillPolygon: %v", err)
	}
	if countSet(filled) <= countSet(outline) {
		t.Errorf("filled polygon (%d px) should set more pixels than outline (%d px)",
			countSet(filled), countSet(outline))
	}
	// A point well inside the triangle should be filled.
	if !isSet(filled, 15, 10) {
		t.Error("interior point (15,10) of filled triangle should be set")
	}
}

func TestFloodFill(t *testing.T) {
	t.Run("whole canvas", func(t *testing.T) {
		img := newWhiteImage(t, 10, 10)
		if err := FloodFill(img, redGC(t), 5, 5); err != nil {
			t.Fatal(err)
		}
		if countSet(img) != 100 {
			t.Errorf("flood fill of white canvas set %d/100 pixels", countSet(img))
		}
	})

	t.Run("bounded by a border", func(t *testing.T) {
		img := newWhiteImage(t, 20, 20)
		black, _ := AllocNamedColor("black")
		bgc := CreateGraphicsContext()
		SetForeground(&bgc, black)
		DrawRectangle(img, bgc, 4, 4, 10, 10) // closed border

		before := countSet(img)
		blue, _ := AllocNamedColor("blue")
		fgc := CreateGraphicsContext()
		SetForeground(&fgc, blue)
		if err := FloodFill(img, fgc, 9, 9); err != nil { // seed inside border
			t.Fatal(err)
		}

		// The fill should add pixels but must not escape to the corners.
		if countSet(img) <= before {
			t.Error("flood fill added no pixels inside the border")
		}
		if !isWhite(img, 0, 0) {
			t.Error("flood fill escaped the border to (0,0)")
		}
	})

	t.Run("seed already fill color is a no-op", func(t *testing.T) {
		img := newWhiteImage(t, 10, 10)
		white, _ := AllocNamedColor("white")
		gc := CreateGraphicsContext()
		SetForeground(&gc, white)
		if err := FloodFill(img, gc, 5, 5); err != nil {
			t.Fatal(err)
		}
		if countSet(img) != 0 {
			t.Error("filling white over white should change nothing")
		}
	})

	t.Run("out-of-bounds seed", func(t *testing.T) {
		img := newWhiteImage(t, 10, 10)
		if err := FloodFill(img, redGC(t), -1, -1); err != nil {
			t.Fatal(err)
		}
		if err := FloodFill(img, redGC(t), 100, 100); err != nil {
			t.Fatal(err)
		}
		if countSet(img) != 0 {
			t.Error("out-of-bounds seed should not draw")
		}
	})
}

func TestCopyImage(t *testing.T) {
	src := newWhiteImage(t, 10, 10)
	red, _ := AllocNamedColor("red")
	sgc := CreateGraphicsContext()
	SetForeground(&sgc, red)
	FillRectangle(src, sgc, 0, 0, 10, 10) // all red

	dst := newWhiteImage(t, 10, 10)
	if err := CopyImage(src, dst, CreateGraphicsContext(), 0, 0, 5, 5, 2, 2); err != nil {
		t.Fatalf("CopyImage: %v", err)
	}
	if !isSet(dst, 3, 3) {
		t.Error("copied region pixel (3,3) should be red")
	}
	if !isWhite(dst, 0, 0) {
		t.Error("(0,0) outside copy target should remain white")
	}
}

func TestCopyImageScaled(t *testing.T) {
	src := newWhiteImage(t, 4, 4)
	red, _ := AllocNamedColor("red")
	sgc := CreateGraphicsContext()
	SetForeground(&sgc, red)
	FillRectangle(src, sgc, 0, 0, 4, 4) // all red

	dst := newWhiteImage(t, 16, 16)
	if err := CopyImageScaled(src, dst, 0, 0, 4, 4, 2, 2, 8, 8); err != nil {
		t.Fatalf("CopyImageScaled: %v", err)
	}
	if !isSet(dst, 5, 5) {
		t.Error("scaled copy should set pixel (5,5)")
	}
	if !isWhite(dst, 0, 0) {
		t.Error("(0,0) outside scaled target should remain white")
	}
}
