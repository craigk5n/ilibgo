// bdfinfo prints the metadata of an X11 BDF bitmap font: foundry, family,
// face name, slant, weight, spacing, pixel size, ascent/descent, and the
// number of defined glyphs. It is a non-rendering consumer of the font parser,
// handy for inspecting a .bdf before using it.
//
// Usage:
//
//	bdfinfo font.bdf [font2.bdf ...]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/craigk5n/ilibgo"
)

func main() {
	flag.Usage = usage
	flag.Parse()
	paths := flag.Args()
	if len(paths) == 0 {
		usage()
		os.Exit(2)
	}

	status := 0
	for _, path := range paths {
		if err := printFont(path); err != nil {
			fmt.Fprintf(os.Stderr, "bdfinfo: %v\n", err)
			status = 1
		}
	}
	os.Exit(status)
}

func printFont(path string) error {
	name := filepath.Base(path)
	f, err := ilibgo.LoadFontFromFile(path, name)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	spacing := "fixed"
	if f.Proportional() {
		spacing = "proportional"
	}

	fmt.Printf("%s\n", path)
	fmt.Printf("  foundry:     %s\n", orNA(f.Foundry()))
	fmt.Printf("  family:      %s\n", orNA(f.Family()))
	fmt.Printf("  face name:   %s\n", orNA(f.FaceName()))
	fmt.Printf("  slant:       %s\n", orNA(f.Slant()))
	fmt.Printf("  weight:      %s\n", orNA(f.Weight()))
	fmt.Printf("  spacing:     %s\n", spacing)
	fmt.Printf("  pixel size:  %d\n", f.PixelSize())
	fmt.Printf("  ascent:      %d\n", f.Ascent())
	fmt.Printf("  descent:     %d\n", f.Descent())
	fmt.Printf("  glyphs:      %d\n", f.GlyphCount())
	return nil
}

func orNA(s string) string {
	if s == "" {
		return "(n/a)"
	}
	return s
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: bdfinfo font.bdf [font2.bdf ...]")
	flag.PrintDefaults()
}
