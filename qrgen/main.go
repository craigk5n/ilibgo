// qrgen renders a QR code to a PNG image.
//
// The text to encode is taken from the command-line arguments, or from stdin
// when no arguments are given.
//
// Usage:
//
//	qrgen [options] text...
//
//	  -out string    output PNG file (default "qr.png")
//	  -ecc string    error-correction level: L, M, Q, or H (default "M")
//	  -scale int     pixels per module (default 8)
//	  -border int    quiet-zone width in modules (default 4)
//
// Examples:
//
//	qrgen -out site.png https://github.com/craigk5n/ilibgo
//	echo "hello" | qrgen -ecc H -scale 12
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/craigk5n/ilibgo"
	"github.com/craigk5n/ilibgo/qr"
)

func main() {
	out := flag.String("out", "qr.png", "output PNG file")
	eccName := flag.String("ecc", "M", "error-correction level: L, M, Q, or H")
	scale := flag.Int("scale", 8, "pixels per module")
	border := flag.Int("border", 4, "quiet-zone width in modules")
	flag.Parse()

	ecl, err := parseEcc(*eccName)
	if err != nil {
		fail(err)
	}
	if *scale < 1 || *border < 0 {
		fail(fmt.Errorf("scale must be >= 1 and border >= 0"))
	}

	text, err := inputText(flag.Args())
	if err != nil {
		fail(err)
	}
	if text == "" {
		fail(fmt.Errorf("no text to encode (pass arguments or pipe via stdin)"))
	}

	code, err := qr.EncodeText(text, ecl)
	if err != nil {
		fail(err)
	}

	img := render(code, *scale, *border)
	f, err := os.Create(*out)
	if err != nil {
		fail(fmt.Errorf("create %s: %w", *out, err))
	}
	defer f.Close()
	if err := ilibgo.WriteImageFile(f, img, ilibgo.FormatPNG); err != nil {
		fail(fmt.Errorf("write %s: %w", *out, err))
	}
	fmt.Printf("wrote %s (version %d, %dx%d modules, %d px)\n",
		*out, code.Version(), code.Size(), code.Size(), (code.Size()+2*(*border))*(*scale))
}

func render(code *qr.Code, scale, border int) *ilibgo.Image {
	dim := (code.Size() + 2*border) * scale
	white, _ := ilibgo.AllocNamedColor("white")
	img := ilibgo.CreateImageWithBackground(dim, dim, white)

	gc := ilibgo.CreateGraphicsContext()
	black, _ := ilibgo.AllocNamedColor("black")
	ilibgo.SetForeground(&gc, black)

	for y := 0; y < code.Size(); y++ {
		for x := 0; x < code.Size(); x++ {
			if code.Module(x, y) {
				px := (border + x) * scale
				py := (border + y) * scale
				img.FillRectangle(gc, px, py, scale, scale)
			}
		}
	}
	return img
}

func parseEcc(name string) (qr.Ecc, error) {
	switch strings.ToUpper(strings.TrimSpace(name)) {
	case "L":
		return qr.Low, nil
	case "M":
		return qr.Medium, nil
	case "Q":
		return qr.Quartile, nil
	case "H":
		return qr.High, nil
	}
	return qr.Medium, fmt.Errorf("invalid ecc level %q (want L, M, Q, or H)", name)
}

func inputText(args []string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(data), "\r\n"), nil
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "qrgen: %v\n", err)
	os.Exit(1)
}
