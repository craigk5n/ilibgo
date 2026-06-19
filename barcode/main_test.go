package main

import "testing"

// TestCode39PatternsWellFormed checks that every pattern is 9 elements long
// with exactly three wide elements — the defining invariant of Code 39. This
// catches transcription errors in the table.
func TestCode39PatternsWellFormed(t *testing.T) {
	for r, pattern := range code39 {
		if len(pattern) != 9 {
			t.Errorf("%q: pattern length %d, want 9", r, len(pattern))
		}
		wide := 0
		for _, c := range pattern {
			switch c {
			case 'w':
				wide++
			case 'n':
			default:
				t.Errorf("%q: invalid element %q (want 'n' or 'w')", r, c)
			}
		}
		if wide != 3 {
			t.Errorf("%q: %d wide elements, want 3", r, wide)
		}
	}
}

func TestValidate(t *testing.T) {
	if err := validate("HELLO-123"); err != nil {
		t.Errorf("validate(HELLO-123) = %v, want nil", err)
	}
	if err := validate(""); err == nil {
		t.Error("validate(\"\") = nil, want error")
	}
	if err := validate("abc*"); err == nil {
		t.Error("validate with '*' = nil, want error")
	}
	if err := validate("héllo"); err == nil {
		t.Error("validate with non-encodable rune = nil, want error")
	}
}

func TestRenderProducesImage(t *testing.T) {
	img, err := render("*AB-12*", 2, 40)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if img == nil {
		t.Fatal("render returned nil image")
	}
}
