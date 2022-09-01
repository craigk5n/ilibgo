// Create an go file from a BDF font file so that the BDF font can be bundled
// with the Go file.  This allows an application to distribute the font
// within the binary rather than having to also include the BDF font file
// separately and then load the BDF font file.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	outfile := ""
	infile := ""
	packageName := "ilibgo/fonts"

	flag.StringVar(&outfile, "outfile", "", "output go filename")
	flag.StringVar(&infile, "infile", "", "input BDF font filename")
	flag.StringVar(&packageName, "package", packageName, "package name to use within Go files")
	flag.Parse()
	if len(infile) == 0 {
		fmt.Printf("Error: an input BDF font file must be specified with -infile\n")
		os.Exit(1)
	}
	fp, err := os.Open(infile)
	if err != nil {
		fmt.Printf("error opening file %s: %v\n", infile, err)
		os.Exit(1)
	}
	defer fp.Close()

	// Create the output filename based on the input BDF input filename
	base := filepath.Base(infile)
	base = strings.Replace(base, ".bdf", "", 1)
	fmt.Printf("Base: %s\n", base)

	if len(outfile) == 0 {
		outfile = fmt.Sprintf("%s.go", base)
	}
	fmt.Printf("Output file: %s\n", outfile)

	outfp, err := os.Create(outfile) // creates a file at current directory
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outfp.Close()

	fmt.Fprintf(outfp, "package %s\n\n", packageName)
	fmt.Fprintf(outfp, "// Return font data for the %s BDF font\n", base)
	fmt.Fprintf(outfp, "func Font_%s() []string {\n", base)
	fmt.Fprintf(outfp, "  return font_%s\n}\n\n", base)
	fmt.Fprintf(outfp, "var font_%s []string = []string{\n", base)

	fileScanner := bufio.NewScanner(fp)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		line = strings.ReplaceAll(line, "\"", "\\\"")
		fmt.Fprintf(outfp, "  \"%s\",\n", line)
	}
	fmt.Fprintf(outfp, "}\n")
}
