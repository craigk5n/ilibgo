package ilibgo

import (
	"bytes"
	"path/filepath"
	"testing"
)

// TestEncodeDecodeRoundTrip exercises the streaming io.Writer/io.Reader API:
// an image encoded to an in-memory buffer must decode back to the same size.
func TestEncodeDecodeRoundTrip(t *testing.T) {
	formats := []struct {
		name   string
		format ImageFormat
	}{
		{"png", FormatPNG},
		{"jpeg", FormatJPEG},
		{"gif", FormatGIF},
		{"bmp", FormatBMP},
		{"tiff", FormatTIFF},
		{"ppm", FormatPPM},
	}
	src := newWhiteImage(t, 12, 8)
	for _, f := range formats {
		t.Run(f.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := Encode(&buf, src, f.format); err != nil {
				t.Fatalf("Encode(%s): %v", f.name, err)
			}
			if buf.Len() == 0 {
				t.Fatalf("Encode(%s) produced no bytes", f.name)
			}
			got, err := Decode(&buf)
			if err != nil {
				t.Fatalf("Decode(%s): %v", f.name, err)
			}
			if ImageWidth(got) != 12 || ImageHeight(got) != 8 {
				t.Errorf("Decode(%s) size = %dx%d, want 12x8", f.name, ImageWidth(got), ImageHeight(got))
			}
		})
	}
}

func TestEncodeUnsupportedFormats(t *testing.T) {
	src := newWhiteImage(t, 4, 4)
	for _, f := range []ImageFormat{FormatPGM, FormatPBM, FormatXPM} {
		var buf bytes.Buffer
		if err := Encode(&buf, src, f); err == nil {
			t.Errorf("Encode(format %d): want error, got nil", f)
		}
	}
}

// TestSaveLoadImageFile covers the path-based convenience wrappers.
func TestSaveLoadImageFile(t *testing.T) {
	src := newWhiteImage(t, 10, 6)
	path := filepath.Join(t.TempDir(), "round.png")
	if err := SaveImageFile(path, src, FormatPNG); err != nil {
		t.Fatalf("SaveImageFile: %v", err)
	}
	got, err := LoadImageFile(path)
	if err != nil {
		t.Fatalf("LoadImageFile: %v", err)
	}
	if ImageWidth(got) != 10 || ImageHeight(got) != 6 {
		t.Errorf("LoadImageFile size = %dx%d, want 10x6", ImageWidth(got), ImageHeight(got))
	}
}

func TestLoadImageFileMissing(t *testing.T) {
	if _, err := LoadImageFile(filepath.Join(t.TempDir(), "nope.png")); err == nil {
		t.Error("LoadImageFile on missing file: want error, got nil")
	}
}

func TestSaveImageFileBadPath(t *testing.T) {
	src := newWhiteImage(t, 4, 4)
	// A path inside a non-existent directory cannot be created.
	bad := filepath.Join(t.TempDir(), "no-such-dir", "x.png")
	if err := SaveImageFile(bad, src, FormatPNG); err == nil {
		t.Error("SaveImageFile to bad path: want error, got nil")
	}
}
