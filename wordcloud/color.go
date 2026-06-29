package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/craigk5n/ilibgo"
)

// parseColor accepts either an X11 color name (e.g. "tomato", "black") or a hex
// triple ("#rrggbb" / "rrggbb"). Names are resolved via the library's embedded
// rgb.txt map.
func parseColor(s string) (ilibgo.Color, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return ilibgo.Color{}, fmt.Errorf("empty color")
	}
	if strings.HasPrefix(s, "#") || isHex6(s) {
		return parseHex(s)
	}
	c, err := ilibgo.AllocNamedColor(s)
	if err != nil {
		return ilibgo.Color{}, fmt.Errorf("unknown color %q: %w", s, err)
	}
	return c, nil
}

// isHex6 reports whether s is six hex digits (a bare "rrggbb").
func isHex6(s string) bool {
	if len(s) != 6 {
		return false
	}
	if _, err := strconv.ParseUint(s, 16, 32); err != nil {
		return false
	}
	return true
}

// parseHex parses "#rrggbb" or "rrggbb" into an opaque Color.
func parseHex(s string) (ilibgo.Color, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return ilibgo.Color{}, fmt.Errorf("hex color must be 6 digits, got %q", s)
	}
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return ilibgo.Color{}, fmt.Errorf("invalid hex color %q: %w", s, err)
	}
	r := uint8(v >> 16)
	g := uint8(v >> 8)
	b := uint8(v)
	return ilibgo.AllocColor(r, g, b)
}
