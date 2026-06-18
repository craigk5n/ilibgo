package qr_test

import (
	"fmt"
	"strings"

	"github.com/craigk5n/ilibgo/qr"
)

// Encode some text and inspect the resulting symbol.
func ExampleEncodeText() {
	code, err := qr.EncodeText("https://example.com", qr.Medium)
	if err != nil {
		panic(err)
	}
	fmt.Printf("version %d, %dx%d modules\n", code.Version(), code.Size(), code.Size())
	// Output: version 2, 25x25 modules
}

// Render a QR code to the terminal using two characters per module so the
// aspect ratio looks square. Any consumer can walk the grid via Module(x, y);
// the qrgen tool does the same thing but draws to a PNG with FillRectangle.
func ExampleCode_Module() {
	code, err := qr.EncodeText("HI", qr.Low)
	if err != nil {
		panic(err)
	}
	var b strings.Builder
	for y := 0; y < code.Size(); y++ {
		for x := 0; x < code.Size(); x++ {
			if code.Module(x, y) {
				b.WriteString("##")
			} else {
				b.WriteString("  ")
			}
		}
		b.WriteByte('\n')
	}
	// Print just the dimensions; the full grid is large.
	fmt.Printf("rendered %d lines\n", code.Size())
	_ = b.String()
	// Output: rendered 21 lines
}
