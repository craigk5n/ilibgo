// iconvert converts an image from one format to another. The output format is
// chosen from the output file's extension.
//
// Usage:
//
//	iconvert [options] input output
//
//	  -info   print the input image's dimensions and format, then exit
//
// Examples:
//
//	iconvert photo.png photo.jpg
//	iconvert -info logo.gif
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/craigk5n/ilibgo"
)

func main() {
	info := flag.Bool("info", false, "print input dimensions/format and exit")
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if *info {
		if len(args) != 1 {
			usage()
			os.Exit(2)
		}
		if err := printInfo(args[0]); err != nil {
			fail(err)
		}
		return
	}

	if len(args) != 2 {
		usage()
		os.Exit(2)
	}
	if err := convert(args[0], args[1]); err != nil {
		fail(err)
	}
}

func convert(inPath, outPath string) error {
	format, err := ilibgo.FileType(outPath)
	if err != nil {
		return fmt.Errorf("output format: %w", err)
	}

	img, err := readImage(inPath)
	if err != nil {
		return err
	}

	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", outPath, err)
	}
	defer out.Close()
	if err := ilibgo.WriteImageFile(out, img, format); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}
	fmt.Printf("wrote %s (%dx%d)\n", outPath, ilibgo.ImageWidth(img), ilibgo.ImageHeight(img))
	return nil
}

func printInfo(inPath string) error {
	img, err := readImage(inPath)
	if err != nil {
		return err
	}
	format, _ := ilibgo.FileType(inPath)
	fmt.Printf("%s: %dx%d, format %v\n", inPath, ilibgo.ImageWidth(img), ilibgo.ImageHeight(img), format)
	return nil
}

func readImage(path string) (*ilibgo.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	img, err := ilibgo.ReadImageFile(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return img, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: iconvert [options] input output")
	fmt.Fprintln(os.Stderr, "       iconvert -info input")
	flag.PrintDefaults()
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "iconvert: %v\n", err)
	os.Exit(1)
}
