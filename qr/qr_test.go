package qr

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestVersionSelectionAndSize(t *testing.T) {
	tests := []struct {
		text    string
		ecl     Ecc
		wantVer int
	}{
		{"HELLO", Low, 1},
		{"https://github.com/craigk5n/ilibgo", Medium, 3},
		{"https://github.com/craigk5n/ilibgo", High, 4},
	}
	for _, tc := range tests {
		c, err := EncodeText(tc.text, tc.ecl)
		if err != nil {
			t.Fatalf("EncodeText(%q): %v", tc.text, err)
		}
		if c.Version() != tc.wantVer {
			t.Errorf("EncodeText(%q, %d) version = %d, want %d", tc.text, tc.ecl, c.Version(), tc.wantVer)
		}
		if want := 17 + 4*c.Version(); c.Size() != want {
			t.Errorf("size = %d, want %d for version %d", c.Size(), want, c.Version())
		}
	}
}

func TestStructuralPatterns(t *testing.T) {
	c, err := EncodeText("structure", Medium)
	if err != nil {
		t.Fatal(err)
	}
	s := c.Size()

	// Finder pattern centers (dark) in three corners.
	for _, p := range [][2]int{{3, 3}, {s - 4, 3}, {3, s - 4}} {
		if !c.Module(p[0], p[1]) {
			t.Errorf("finder center (%d,%d) should be dark", p[0], p[1])
		}
	}
	// Always-dark module.
	if !c.Module(8, s-8) {
		t.Error("module (8, size-8) must always be dark")
	}
	// Timing patterns alternate along row/column 6.
	for i := 8; i < s-8; i++ {
		want := i%2 == 0
		if c.Module(i, 6) != want || c.Module(6, i) != want {
			t.Errorf("timing pattern wrong at index %d", i)
			break
		}
	}
}

func TestReedSolomon(t *testing.T) {
	if reedSolomonMultiply(0, 5) != 0 {
		t.Error("RS multiply by 0 should be 0")
	}
	if reedSolomonMultiply(1, 5) != 5 || reedSolomonMultiply(5, 1) != 5 {
		t.Error("RS multiply by 1 should be identity")
	}
	// Generator polynomial for 2 EC codewords is x^2 + 3x + 2.
	if got := reedSolomonDivisor(2); len(got) != 2 || got[0] != 3 || got[1] != 2 {
		t.Errorf("reedSolomonDivisor(2) = %v, want [3 2]", got)
	}
}

// TestGoldenFingerprint locks the module layout for known inputs. The expected
// hashes were captured from output independently verified to scan with a real
// (ZXing-port) QR decoder, so a change here flags a regression in the encoder.
func TestGoldenFingerprint(t *testing.T) {
	tests := []struct {
		text   string
		ecl    Ecc
		prefix string
	}{
		{"https://github.com/craigk5n/ilibgo", Medium, "e5363e029811eef6"},
		{"HELLO", Low, "586b654cba5fa196"},
	}
	for _, tc := range tests {
		c, err := EncodeText(tc.text, tc.ecl)
		if err != nil {
			t.Fatal(err)
		}
		h := sha256.New()
		for y := 0; y < c.Size(); y++ {
			for x := 0; x < c.Size(); x++ {
				b := byte(0)
				if c.Module(x, y) {
					b = 1
				}
				h.Write([]byte{b})
			}
		}
		got := hex.EncodeToString(h.Sum(nil))[:16]
		if got != tc.prefix {
			t.Errorf("fingerprint for %q = %s, want %s", tc.text, got, tc.prefix)
		}
	}
}

func TestEncodeErrors(t *testing.T) {
	if _, err := Encode(make([]byte, 4000), High); err == nil {
		t.Error("expected error for data exceeding max capacity")
	}
	if _, err := Encode([]byte("x"), Ecc(99)); err == nil {
		t.Error("expected error for invalid ecc level")
	}
}

func TestEccLevelsAllEncode(t *testing.T) {
	for _, ecl := range []Ecc{Low, Medium, Quartile, High} {
		c, err := EncodeText("the quick brown fox jumps over the lazy dog", ecl)
		if err != nil {
			t.Fatalf("ecc %d: %v", ecl, err)
		}
		if c.Size() <= 0 {
			t.Errorf("ecc %d: bad size %d", ecl, c.Size())
		}
	}
}
