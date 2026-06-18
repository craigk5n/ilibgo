package ilibgo

import "testing"

func TestFormatStringToType(t *testing.T) {
	tests := []struct {
		in      string
		want    ImageFormat
		wantErr bool
	}{
		{"gif", FormatGIF, false},
		{"ppm", FormatPPM, false},
		{"pgm", FormatPGM, false},
		{"pbm", FormatPBM, false},
		{"xpm", FormatXPM, false},
		{"xbm", FormatXBM, false},
		{"png", FormatPNG, false},
		{"PNG", FormatPNG, false},  // case-insensitive
		{".png", FormatPNG, false}, // leading dot stripped
		{"jpeg", FormatJPEG, false},
		{"jpg", FormatJPEG, false},
		{"bmp", FormatBMP, false},
		{"tiff", FormatTIFF, false},
		{"webp", FormatPPM, true}, // unsupported
		{"", FormatPPM, true},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := FormatStringToType(tc.in)
			if tc.wantErr != (err != nil) {
				t.Fatalf("FormatStringToType(%q) err = %v, wantErr = %v", tc.in, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("FormatStringToType(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestFileType(t *testing.T) {
	got, err := FileType("photo.JPG")
	if err != nil || got != FormatJPEG {
		t.Errorf("FileType(photo.JPG) = %v, %v; want FormatJPEG, nil", got, err)
	}
	if _, err := FileType("noextension"); err == nil {
		t.Error("FileType(noextension) expected error, got nil")
	}
}

func TestIsSupportedFormat(t *testing.T) {
	if !IsSupportedFormat("png") {
		t.Error("IsSupportedFormat(png) = false, want true")
	}
	if IsSupportedFormat("nope") {
		t.Error("IsSupportedFormat(nope) = true, want false")
	}
}
