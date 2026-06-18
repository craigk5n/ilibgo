package ilibgo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateImageDimensions(t *testing.T) {
	img := CreateImage(13, 27)
	if ImageWidth(img) != 13 || ImageHeight(img) != 27 {
		t.Errorf("CreateImage(13,27) dims = %dx%d, want 13x27",
			ImageWidth(img), ImageHeight(img))
	}
}

func TestCreateImageWithBackground(t *testing.T) {
	red, _ := AllocNamedColor("red")
	img := CreateImageWithBackground(5, 5, red)
	c := GetPoint(img, 2, 2)
	if c.color.R != 255 || c.color.G != 0 || c.color.B != 0 {
		t.Errorf("background pixel = %v, want red", c.color)
	}
}

// paintSample fills an image with a handful of distinct, palette-safe colors.
func paintSample(t *testing.T) *Image {
	t.Helper()
	img := newWhiteImage(t, 8, 8)
	gc := CreateGraphicsContext()
	colors := []string{"red", "green", "blue", "black", "yellow", "cyan"}
	for i, name := range colors {
		c, err := AllocNamedColor(name)
		if err != nil {
			t.Fatalf("AllocNamedColor(%q): %v", name, err)
		}
		SetForeground(&gc, c)
		SetPoint(img, gc, i, i)
	}
	return img
}

func writeRead(t *testing.T, img *Image, format ImageFormat, ext string) *Image {
	t.Helper()
	path := filepath.Join(t.TempDir(), "img."+ext)

	out, err := os.Create(path)
	if err != nil {
		t.Fatalf("create %s: %v", path, err)
	}
	if err := WriteImageFile(out, img, format); err != nil {
		out.Close()
		t.Fatalf("WriteImageFile(%s): %v", ext, err)
	}
	if err := out.Close(); err != nil {
		t.Fatalf("close after write: %v", err)
	}

	in, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer in.Close()
	got, err := ReadImageFile(in)
	if err != nil {
		t.Fatalf("ReadImageFile(%s): %v", ext, err)
	}
	return got
}

func TestRoundTripLossless(t *testing.T) {
	cases := []struct {
		format ImageFormat
		ext    string
	}{
		{FormatPNG, "png"},
		{FormatBMP, "bmp"},
		{FormatTIFF, "tiff"},
		{FormatPPM, "ppm"},
	}
	for _, tc := range cases {
		t.Run(tc.ext, func(t *testing.T) {
			src := paintSample(t)
			got := writeRead(t, src, tc.format, tc.ext)
			if ImageWidth(got) != ImageWidth(src) || ImageHeight(got) != ImageHeight(src) {
				t.Fatalf("%s dims = %dx%d, want %dx%d", tc.ext,
					ImageWidth(got), ImageHeight(got), ImageWidth(src), ImageHeight(src))
			}
			for y := 0; y < ImageHeight(src); y++ {
				for x := 0; x < ImageWidth(src); x++ {
					s := GetPoint(src, x, y).color
					g := GetPoint(got, x, y).color
					if s.R != g.R || s.G != g.G || s.B != g.B {
						t.Fatalf("%s pixel (%d,%d) = %v, want %v", tc.ext, x, y, g, s)
					}
				}
			}
		})
	}
}

func TestRoundTripLossy(t *testing.T) {
	// JPEG and GIF need not preserve every pixel; verify dimensions survive.
	cases := []struct {
		format ImageFormat
		ext    string
	}{
		{FormatJPEG, "jpg"},
		{FormatGIF, "gif"},
	}
	for _, tc := range cases {
		t.Run(tc.ext, func(t *testing.T) {
			src := paintSample(t)
			got := writeRead(t, src, tc.format, tc.ext)
			if ImageWidth(got) != ImageWidth(src) || ImageHeight(got) != ImageHeight(src) {
				t.Errorf("%s dims = %dx%d, want %dx%d", tc.ext,
					ImageWidth(got), ImageHeight(got), ImageWidth(src), ImageHeight(src))
			}
		})
	}
}

func TestWriteUnsupportedFormats(t *testing.T) {
	for _, format := range []ImageFormat{FormatPGM, FormatPBM, FormatXPM} {
		img := CreateImage(4, 4)
		path := filepath.Join(t.TempDir(), "x")
		f, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		err = WriteImageFile(f, img, format)
		f.Close()
		if err == nil {
			t.Errorf("WriteImageFile(format=%d) expected unsupported error, got nil", format)
		}
	}
}

func TestWriteImageFilePropagatesEncoderError(t *testing.T) {
	// Open the file read-only so the encoder's writes fail; WriteImageFile
	// must surface that error rather than swallow it.
	path := filepath.Join(t.TempDir(), "ro.png")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := WriteImageFile(f, CreateImage(4, 4), FormatPNG); err == nil {
		t.Error("WriteImageFile to read-only file returned nil, want error")
	}
}

func TestReadImageFileError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "garbage.png")
	if err := os.WriteFile(path, []byte("not an image"), 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := ReadImageFile(f); err == nil {
		t.Error("ReadImageFile on garbage returned nil error, want decode error")
	}
}
