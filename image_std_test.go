package ilibgo

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"testing"
)

func TestImageStdInterfaceMethods(t *testing.T) {
	img := CreateImage(8, 6)

	if got := img.Bounds(); got != image.Rect(0, 0, 8, 6) {
		t.Errorf("Bounds() = %v, want (0,0)-(8,6)", got)
	}
	if img.ColorModel() != color.RGBAModel {
		t.Error("ColorModel() should be color.RGBAModel")
	}

	// Set via draw.Image, read back via image.Image.
	img.Set(2, 3, color.RGBA{10, 20, 30, 255})
	r, g, b, a := img.At(2, 3).RGBA()
	if r>>8 != 10 || g>>8 != 20 || b>>8 != 30 || a>>8 != 255 {
		t.Errorf("At(2,3) 8-bit = (%d,%d,%d,%d), want (10,20,30,255)", r>>8, g>>8, b>>8, a>>8)
	}
}

func TestImageInteropWithStdlib(t *testing.T) {
	img := CreateImage(10, 10)

	// *Image works as a draw.Image destination.
	red := color.RGBA{255, 0, 0, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: red}, image.Point{}, draw.Src)
	if c, ok := img.At(5, 5).(color.RGBA); !ok || c != red {
		t.Errorf("after draw.Draw, At(5,5) = %v, want %v", img.At(5, 5), red)
	}

	// *Image works as an image.Image source for the standard encoders.
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode(*Image): %v", err)
	}
	decoded, err := png.Decode(&buf)
	if err != nil {
		t.Fatalf("png.Decode: %v", err)
	}
	if decoded.Bounds() != img.Bounds() {
		t.Errorf("decoded bounds %v != original %v", decoded.Bounds(), img.Bounds())
	}
}
