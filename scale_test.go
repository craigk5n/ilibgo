package ilibgo

import "testing"

// TestCopyImageScaledQuality verifies each resampling filter upscales a source
// onto the destination and produces non-background output.
func TestCopyImageScaledQuality(t *testing.T) {
	qualities := []struct {
		name string
		q    ScaleQuality
	}{
		{"nearest", ScaleNearestNeighbor},
		{"approx-bilinear", ScaleApproxBiLinear},
		{"bilinear", ScaleBiLinear},
		{"catmull-rom", ScaleCatmullRom},
	}
	for _, tc := range qualities {
		t.Run(tc.name, func(t *testing.T) {
			// 2x2 source: one red pixel, rest white.
			src := newWhiteImage(t, 2, 2)
			src.SetPoint(redGC(t), 0, 0)

			dst := newWhiteImage(t, 20, 20)
			if err := dst.CopyImageScaledQuality(src, 0, 0, 2, 2, 0, 0, 20, 20, tc.q); err != nil {
				t.Fatalf("CopyImageScaledQuality(%s): %v", tc.name, err)
			}
			// The upper-left region scaled from the red pixel must be non-white.
			if !isSet(dst, 2, 2) {
				t.Errorf("%s: expected scaled red region at (2,2), got white", tc.name)
			}
		})
	}
}

// TestCopyImageScaledQualityDownscale checks high-quality downscaling runs and
// fills the destination.
func TestCopyImageScaledQualityDownscale(t *testing.T) {
	src := newWhiteImage(t, 40, 40)
	// Fill the source with red so any sampled output is non-white.
	gc := redGC(t)
	src.FillRectangle(gc, 0, 0, 40, 40)

	dst := newWhiteImage(t, 10, 10)
	if err := dst.CopyImageScaledQuality(src, 0, 0, 40, 40, 0, 0, 10, 10, ScaleCatmullRom); err != nil {
		t.Fatalf("downscale: %v", err)
	}
	if countSet(dst) == 0 {
		t.Error("downscaled image is entirely background")
	}
}
