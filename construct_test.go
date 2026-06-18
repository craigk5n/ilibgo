package ilibgo

import (
	"image/color"
	"testing"
)

// Color must satisfy the standard image/color.Color interface.
var _ color.Color = Color{}

func TestPtConstructor(t *testing.T) {
	p := Pt(3, 4)
	if p.X != 3 || p.Y != 4 {
		t.Errorf("Pt(3,4) = %+v, want {X:3 Y:4}", p)
	}
	// Equivalent to a struct literal with exported fields.
	if p != (Point{X: 3, Y: 4}) {
		t.Error("Pt and struct literal should produce equal points")
	}
}

func TestNewColorAndRGBA(t *testing.T) {
	c := NewColor(255, 0, 0, 255)
	r, g, b, a := c.RGBA()
	// color.RGBA scales 8-bit components to 16-bit (x * 0x101).
	if r != 0xffff || g != 0 || b != 0 || a != 0xffff {
		t.Errorf("NewColor(255,0,0,255).RGBA() = (%d,%d,%d,%d), want (65535,0,0,65535)", r, g, b, a)
	}

	// Alpha is honored (unlike AllocColor, which is always opaque).
	semi := NewColor(0, 0, 0, 128)
	if _, _, _, sa := semi.RGBA(); sa != uint32(128)|uint32(128)<<8 {
		t.Errorf("NewColor alpha not preserved: got %d", sa)
	}
}

// TestExternalConstructionFlow exercises the path that was previously
// impossible from outside the package: build points and a color from scratch
// and render a filled polygon with them.
func TestExternalConstructionFlow(t *testing.T) {
	img := newWhiteImage(t, 30, 30)
	gc := CreateGraphicsContext()
	SetForeground(&gc, NewColor(10, 200, 30, 255))

	triangle := []Point{Pt(5, 2), Pt(25, 2), Pt(15, 25)}
	if err := FillPolygon(img, gc, triangle); err != nil {
		t.Fatal(err)
	}
	if !isSet(img, 15, 10) {
		t.Error("polygon built from Pt() should fill its interior")
	}
}
