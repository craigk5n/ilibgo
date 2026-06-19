package ilibgo

import "testing"

// Benchmarks anchoring the performance-sensitive paths (IDEAS §8): fills,
// nearest-neighbor vs high-quality scaling, flood fill, and per-pixel drawing.

func benchImage(w, h int) *Image {
	white := NewColor(255, 255, 255, 255)
	return CreateImageWithBackground(w, h, white)
}

func benchRedGC() GraphicsContext {
	gc := CreateGraphicsContext()
	SetForeground(&gc, NewColor(255, 0, 0, 255))
	return gc
}

func BenchmarkFillRectangle(b *testing.B) {
	img := benchImage(512, 512)
	gc := benchRedGC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		img.FillRectangle(gc, 16, 16, 480, 480)
	}
}

func BenchmarkCopyImageScaledNearest(b *testing.B) {
	src := benchImage(256, 256)
	src.FillRectangle(benchRedGC(), 0, 0, 256, 256)
	dst := benchImage(512, 512)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst.CopyImageScaled(src, 0, 0, 256, 256, 0, 0, 512, 512)
	}
}

func BenchmarkCopyImageScaledCatmullRom(b *testing.B) {
	src := benchImage(256, 256)
	src.FillRectangle(benchRedGC(), 0, 0, 256, 256)
	dst := benchImage(512, 512)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst.CopyImageScaledQuality(src, 0, 0, 256, 256, 0, 0, 512, 512, ScaleCatmullRom)
	}
}

func BenchmarkFloodFill(b *testing.B) {
	gc := benchRedGC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		img := benchImage(128, 128)
		b.StartTimer()
		img.FloodFill(gc, 64, 64)
	}
}

func BenchmarkSetPointFill(b *testing.B) {
	img := benchImage(256, 256)
	gc := benchRedGC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for y := 0; y < 256; y++ {
			for x := 0; x < 256; x++ {
				img.SetPoint(gc, x, y)
			}
		}
	}
}
