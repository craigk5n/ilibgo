// iresize scales an image to a new size using a selectable resampling filter.
// It demonstrates CopyImageScaledQuality. The output format is chosen from the
// output file's extension.
//
// Usage:
//
//	iresize [options] input output
//
//	  -w, -h     target dimensions in pixels (0 = derive from the other,
//	             preserving aspect ratio; at least one must be > 0)
//	  -filter    nearest | approx | bilinear | catmullrom (default catmullrom)
//
// Examples:
//
//	iresize -w 320 photo.jpg thumb.png
//	iresize -w 800 -h 600 -filter bilinear in.png out.png
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/craigk5n/ilibgo"
)

func main() {
	w := flag.Int("w", 0, "target width (0 = derive from height)")
	h := flag.Int("h", 0, "target height (0 = derive from width)")
	filterName := flag.String("filter", "catmullrom", "nearest|approx|bilinear|catmullrom")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 || (*w <= 0 && *h <= 0) {
		usage()
		os.Exit(2)
	}

	filter, err := parseFilter(*filterName)
	if err != nil {
		fail(err)
	}
	if err := resize(args[0], args[1], *w, *h, filter); err != nil {
		fail(err)
	}
}

func resize(inPath, outPath string, w, h int, filter ilibgo.ScaleQuality) error {
	src, err := ilibgo.LoadImageFile(inPath)
	if err != nil {
		return err
	}
	sw, sh := ilibgo.ImageWidth(src), ilibgo.ImageHeight(src)
	if sw < 1 || sh < 1 {
		return fmt.Errorf("source has zero dimension")
	}

	dw, dh := targetSize(sw, sh, w, h)

	format, err := ilibgo.FileType(outPath)
	if err != nil {
		return fmt.Errorf("output format: %w", err)
	}

	dst := ilibgo.CreateImage(dw, dh)
	if err := dst.CopyImageScaledQuality(src, 0, 0, sw, sh, 0, 0, dw, dh, filter); err != nil {
		return err
	}
	if err := ilibgo.SaveImageFile(outPath, dst, format); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}
	fmt.Printf("wrote %s (%dx%d -> %dx%d)\n", outPath, sw, sh, dw, dh)
	return nil
}

// targetSize resolves the requested width/height, filling in a zero dimension
// from the other to preserve the source aspect ratio.
func targetSize(sw, sh, w, h int) (int, int) {
	switch {
	case w > 0 && h > 0:
		return w, h
	case w > 0:
		return w, max(1, int(float64(w)*float64(sh)/float64(sw)+0.5))
	default:
		return max(1, int(float64(h)*float64(sw)/float64(sh)+0.5)), h
	}
}

func parseFilter(name string) (ilibgo.ScaleQuality, error) {
	switch name {
	case "nearest":
		return ilibgo.ScaleNearestNeighbor, nil
	case "approx":
		return ilibgo.ScaleApproxBiLinear, nil
	case "bilinear":
		return ilibgo.ScaleBiLinear, nil
	case "catmullrom":
		return ilibgo.ScaleCatmullRom, nil
	}
	return 0, fmt.Errorf("unknown filter %q (want nearest|approx|bilinear|catmullrom)", name)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: iresize [-w N] [-h N] [-filter F] input output")
	flag.PrintDefaults()
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "iresize: %v\n", err)
	os.Exit(1)
}
