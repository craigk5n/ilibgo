package ilibgo

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	font "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)

// renderGoldenScene draws a fixed scene using only deterministic, integer-only
// operations (solid fills, horizontal/vertical lines, and BDF bitmap text — no
// anti-aliasing or floating-point rasterization). Its pixel output is therefore
// stable across platforms and Go versions, which makes a hash-based golden test
// reliable.
func renderGoldenScene(t *testing.T) *Image {
	t.Helper()
	white, _ := AllocNamedColor("white")
	red, _ := AllocNamedColor("red")
	blue, _ := AllocNamedColor("blue")
	black, _ := AllocNamedColor("black")

	img := CreateImageWithBackground(160, 48, white)
	gc := CreateGraphicsContext()

	SetForeground(&gc, red)
	img.FillRectangle(gc, 4, 4, 40, 40)

	SetForeground(&gc, blue)
	img.DrawLine(gc, 0, 0, 159, 0)  // horizontal
	img.DrawLine(gc, 0, 0, 0, 47)   // vertical
	img.DrawLine(gc, 60, 4, 60, 44) // vertical

	f, err := LoadFontFromData("helvB12", font.Font_helvB12())
	if err != nil {
		t.Fatalf("LoadFontFromData: %v", err)
	}
	SetForeground(&gc, black)
	SetFont(&gc, f)
	img.DrawString(gc, 70, 28, "Golden")

	return img
}

// TestGoldenScene guards against unintended visual regressions in the
// deterministic drawing primitives. If a change legitimately alters the
// rendering, update goldenSceneHash to the new value reported by this test.
func TestGoldenScene(t *testing.T) {
	const goldenSceneHash = "b69a266765604e660b33c6c82da4fedf1213ef20e28e3baa7d2dbb740296c8a1"

	img := renderGoldenScene(t)
	sum := sha256.Sum256(img.data.Pix)
	got := hex.EncodeToString(sum[:])

	if goldenSceneHash == "REPLACE_ME" {
		t.Fatalf("golden hash not set; computed hash is %s", got)
	}
	if got != goldenSceneHash {
		t.Errorf("golden scene hash changed:\n got  %s\n want %s\n"+
			"If this change is intentional, update goldenSceneHash.", got, goldenSceneHash)
	}
}
