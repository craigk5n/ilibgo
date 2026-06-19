// fontsheet renders a specimen catalog of the bundled BDF fonts: a sample line
// drawn in each font, stacked and labeled. It loads the fonts straight from the
// embedded accessors (no external .bdf needed), so it doubles as visual QA for
// the bundled font set.
//
// Usage:
//
//	fontsheet [options]
//
//	  -out string    output PNG file (default "fontsheet.png")
//	  -text string   sample text drawn in each font
//	  -w int         image width (default 760)
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/craigk5n/ilibgo"
	adobe "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
	utopia "github.com/craigk5n/ilibgo/fonts/adobe_utopia_100dpi"
	lucida "github.com/craigk5n/ilibgo/fonts/bh_lucidatypewriter_100dpi"
)

type specimen struct {
	name  string
	lines func() []string
}

// A curated cross-section of the bundled fonts (families, styles, sizes).
var catalog = []specimen{
	{"adobe courR12", adobe.Font_courR12},
	{"adobe courB12", adobe.Font_courB12},
	{"adobe courO12", adobe.Font_courO12},
	{"adobe helvR12", adobe.Font_helvR12},
	{"adobe helvB12", adobe.Font_helvB12},
	{"adobe helvO12", adobe.Font_helvO12},
	{"adobe timR12", adobe.Font_timR12},
	{"adobe timB12", adobe.Font_timB12},
	{"adobe timI12", adobe.Font_timI12},
	{"adobe ncenR12", adobe.Font_ncenR12},
	{"adobe symb12", adobe.Font_symb12},
	{"utopia UTRG 12", utopia.Font_UTRG__12},
	{"utopia UTB 12", utopia.Font_UTB___12},
	{"utopia UTI 12", utopia.Font_UTI___12},
	{"lucida lutRS12", lucida.Font_lutRS12},
	{"lucida lutBS12", lucida.Font_lutBS12},
	{"adobe helvR18", adobe.Font_helvR18},
	{"adobe helvR24", adobe.Font_helvR24},
}

const (
	leftMargin  = 12
	textColumn  = 170
	topMargin   = 16
	rowPadding  = 10
	labelFontPx = 11 // helvR10 is ~11px tall
)

func main() {
	out := flag.String("out", "fontsheet.png", "output PNG file")
	text := flag.String("text", "Quick Brown Fox 0123 !?@#$%", "sample text drawn in each font")
	width := flag.Int("w", 760, "image width")
	flag.Parse()

	if err := run(*out, *text, *width); err != nil {
		fmt.Fprintf(os.Stderr, "fontsheet: %v\n", err)
		os.Exit(1)
	}
}

func run(out, text string, width int) error {
	labelFont, err := ilibgo.LoadFontFromData("helvR10", adobe.Font_helvR10())
	if err != nil {
		return err
	}

	// First pass: load fonts and measure row heights.
	type row struct {
		name string
		font *ilibgo.Font
		px   int
	}
	rows := make([]row, 0, len(catalog))
	total := topMargin
	for _, sp := range catalog {
		f, err := ilibgo.LoadFontFromData(sp.name, sp.lines())
		if err != nil {
			return fmt.Errorf("load %s: %w", sp.name, err)
		}
		px, _ := ilibgo.GetFontSize(f)
		h := px
		if labelFontPx > h {
			h = labelFontPx
		}
		rows = append(rows, row{name: sp.name, font: f, px: px})
		total += h + rowPadding
	}
	height := total + topMargin

	white, _ := ilibgo.AllocNamedColor("white")
	img := ilibgo.CreateImageWithBackground(width, height, white)
	gc := ilibgo.CreateGraphicsContext()
	black, _ := ilibgo.AllocNamedColor("black")
	gray, _ := ilibgo.AllocNamedColor("gray")

	y := topMargin
	for _, r := range rows {
		rowH := r.px
		if labelFontPx > rowH {
			rowH = labelFontPx
		}
		baseline := y + rowH

		// Font-name label in a fixed small font.
		ilibgo.SetFont(&gc, labelFont)
		ilibgo.SetForeground(&gc, gray)
		img.DrawString(gc, leftMargin, baseline, r.name)

		// Sample text in the font itself.
		ilibgo.SetFont(&gc, r.font)
		ilibgo.SetForeground(&gc, black)
		img.DrawString(gc, textColumn, baseline, text)

		y += rowH + rowPadding
	}

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("create %s: %w", out, err)
	}
	defer f.Close()
	if err := ilibgo.WriteImageFile(f, img, ilibgo.FormatPNG); err != nil {
		return fmt.Errorf("write %s: %w", out, err)
	}
	fmt.Printf("wrote %s (%dx%d, %d fonts)\n", out, width, height, len(rows))
	return nil
}
