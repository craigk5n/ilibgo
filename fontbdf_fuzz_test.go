package ilibgo

import (
	"strings"
	"testing"
)

// FuzzLoadFontFromBytes feeds arbitrary bytes to the BDF parser to ensure it
// never panics on malformed, truncated, or hostile input. Errors are fine;
// crashes are not. This guards the hand-rolled field parsing and bitmap
// unpacking (IDEAS §6.3, §3.2).
func FuzzLoadFontFromBytes(f *testing.F) {
	// Seed with snippets that exercise the keyword and bitmap paths, plus
	// deliberately broken lines (short keywords, missing values, bad hex).
	seeds := []string{
		"",
		"STARTCHAR",
		"STARTCHAR A\nENCODING 65\n",
		"BBX 4 6 0 0\nBITMAP\nF0\n90\nF0\nENDCHAR\n",
		"STARTCHAR A\nENCODING 65\nBBX 8 8 0 0\nBITMAP\nFF\nGG\n00\nENDCHAR\n",
		"FONT_ASCENT\nFONT_DESCENT 3\nSLANT\nBBX -1 -1 0 0\n",
		"ENDCHAR\n",
		"BBX 2 2 0 0\nBITMAP\nFFFFFFFF\nENDCHAR\n",
		"PIXEL_SIZE notanumber\nENCODING 999999999999999999999\n",
	}
	for _, s := range seeds {
		f.Add([]byte(s))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		// Must not panic. The result may be an error or a (possibly empty) font.
		font, err := LoadFontFromBytes("fuzz", data)
		if err == nil && font != nil {
			// Exercise the lookup paths too — they must also be panic-safe.
			_ = FontBDFGetChar(font.font, "A")
			_ = FontBDFGetChar(font.font, "nonexistent-glyph")
			if !strings.Contains(string(data), "\x00") {
				_ = font.font // touch field to avoid unused warnings
			}
		}
	})
}
