package ilibgo

import "testing"

func TestAllocColor(t *testing.T) {
	c, err := AllocColor(10, 20, 30)
	if err != nil {
		t.Fatalf("AllocColor: unexpected error %v", err)
	}
	if c.color.R != 10 || c.color.G != 20 || c.color.B != 30 || c.color.A != 255 {
		t.Errorf("AllocColor(10,20,30) = %v, want {10 20 30 255}", c.color)
	}
}

func TestAllocNamedColor(t *testing.T) {
	tests := []struct {
		name    string
		wantR   uint8
		wantG   uint8
		wantB   uint8
		wantErr bool
	}{
		{"red", 255, 0, 0, false},
		{"green", 0, 255, 0, false}, // X11 green is pure green
		{"blue", 0, 0, 255, false},
		{"white", 255, 255, 255, false},
		{"black", 0, 0, 0, false},
		{"definitelynotacolor", 255, 255, 255, true}, // falls back to white + error
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, err := AllocNamedColor(tc.name)
			if tc.wantErr && err == nil {
				t.Fatalf("AllocNamedColor(%q): expected error, got nil", tc.name)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("AllocNamedColor(%q): unexpected error %v", tc.name, err)
			}
			if c.color.R != tc.wantR || c.color.G != tc.wantG || c.color.B != tc.wantB {
				t.Errorf("AllocNamedColor(%q) = {%d %d %d}, want {%d %d %d}",
					tc.name, c.color.R, c.color.G, c.color.B, tc.wantR, tc.wantG, tc.wantB)
			}
		})
	}
}

func TestColorsMatch(t *testing.T) {
	red1, _ := AllocColor(255, 0, 0)
	red2 := NewColor(255, 0, 0, 128) // same RGB, different alpha
	blue, _ := AllocColor(0, 0, 255)

	if !colorsMatch(red1, red2) {
		t.Error("colorsMatch should ignore alpha and report equal RGB as matching")
	}
	if colorsMatch(red1, blue) {
		t.Error("colorsMatch reported different colors as matching")
	}
}
