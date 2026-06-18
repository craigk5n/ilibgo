package ilibgo

import "testing"

// TestImageMethodAPI exercises the canonical method-style API directly (the
// free functions in deprecated.go forward to these).
func TestImageMethodAPI(t *testing.T) {
	shape := []struct {
		name string
		draw func(*Image, GraphicsContext)
	}{
		{"SetPoint/DrawPoint", func(im *Image, gc GraphicsContext) { im.SetPoint(gc, 1, 1); im.DrawPoint(gc, 2, 2) }},
		{"DrawLine", func(im *Image, gc GraphicsContext) { im.DrawLine(gc, 2, 2, 30, 30) }},
		{"DrawRectangle", func(im *Image, gc GraphicsContext) { im.DrawRectangle(gc, 2, 2, 20, 20) }},
		{"FillRectangle", func(im *Image, gc GraphicsContext) { im.FillRectangle(gc, 2, 2, 20, 20) }},
		{"DrawArc", func(im *Image, gc GraphicsContext) { im.DrawArc(gc, 20, 20, 10, 10, 0, 90) }},
		{"FillArc", func(im *Image, gc GraphicsContext) { im.FillArc(gc, 20, 20, 12, 12, 0, 180) }},
		{"DrawEnclosedArc", func(im *Image, gc GraphicsContext) { im.DrawEnclosedArc(gc, 20, 20, 10, 10, 0, 90) }},
		{"DrawCircle", func(im *Image, gc GraphicsContext) { im.DrawCircle(gc, 20, 20, 10) }},
		{"FillCircle", func(im *Image, gc GraphicsContext) { im.FillCircle(gc, 20, 20, 10) }},
		{"DrawEllipse", func(im *Image, gc GraphicsContext) { im.DrawEllipse(gc, 20, 20, 14, 8) }},
		{"DrawPolygon", func(im *Image, gc GraphicsContext) { im.DrawPolygon(gc, []Point{Pt(2, 2), Pt(30, 2), Pt(15, 30)}) }},
		{"FillPolygon", func(im *Image, gc GraphicsContext) { im.FillPolygon(gc, []Point{Pt(2, 2), Pt(30, 2), Pt(15, 30)}) }},
	}
	for _, s := range shape {
		t.Run(s.name, func(t *testing.T) {
			img := newWhiteImage(t, 40, 40)
			s.draw(img, redGC(t))
			if countSet(img) == 0 {
				t.Errorf("%s method drew nothing", s.name)
			}
		})
	}

	t.Run("GetPoint", func(t *testing.T) {
		img := newWhiteImage(t, 5, 5)
		img.SetPoint(redGC(t), 2, 2)
		if c := img.GetPoint(2, 2); c.color.R != 255 || c.color.G != 0 {
			t.Errorf("GetPoint method = %v, want red", c.color)
		}
	})

	t.Run("FloodFill", func(t *testing.T) {
		img := newWhiteImage(t, 10, 10)
		if err := img.FloodFill(redGC(t), 5, 5); err != nil {
			t.Fatal(err)
		}
		if countSet(img) != 100 {
			t.Errorf("FloodFill method filled %d/100", countSet(img))
		}
	})

	t.Run("CopyImage/CopyImageScaled", func(t *testing.T) {
		src := newWhiteImage(t, 8, 8)
		src.FillRectangle(redGC(t), 0, 0, 8, 8)
		dst := newWhiteImage(t, 20, 20)
		if err := dst.CopyImage(src, CreateGraphicsContext(), 0, 0, 4, 4, 1, 1); err != nil {
			t.Fatal(err)
		}
		if err := dst.CopyImageScaled(src, 0, 0, 8, 8, 8, 8, 8, 8); err != nil {
			t.Fatal(err)
		}
		if !isSet(dst, 2, 2) || !isSet(dst, 10, 10) {
			t.Error("CopyImage/CopyImageScaled methods did not draw")
		}
	})

	t.Run("text methods", func(t *testing.T) {
		gc, _ := textGC(t)
		img := newWhiteImage(t, 160, 120)
		img.DrawString(gc, 10, 40, "Hi")
		img.DrawStringRotated(gc, 10, 80, "Hi", TextBottomToTop)
		img.DrawStringRotatedAngle(gc, 10, 110, "Hi", 20.0)
		if countSet(img) == 0 {
			t.Error("text methods drew nothing")
		}
	})
}
