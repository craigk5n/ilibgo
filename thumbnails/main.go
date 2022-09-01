// Build an image that contains an index of a bunch of images.
//
// Usage:
//     thumbnails [-w N] [-h N] [-f] [-s] [-h] outimage inimage1 ...
//
//	-w N	sets the icon width, Example: -w120
//	-H N	sets the icon height, Example: -h120
//	-f	go ahead and overwrite the output if it already exists
//	-s	silent
//	-html f	specify an output file to write an HTML client-side image map
//
//	File format types will be determined by their filename
//	extension (".jpg", ".gif", ".png", etc.)
//
// Note:
//	You can use '-' for the outimage to write to standard output:
//	  ./thumbnails - file1.jpg file2.jpg file3.jpg > index.png
//	Note: you can only use PNG for output this way.
package main

// History:
//      17-Aug-2022	Craig Knudsen craig@k5n.us
//			Converted from C to Go
//		09-Dec-1999	Craig Knudsen	cknudsen@radix.net
//			Display help if no arguments are given
//		28-Sep-1999	Craig Knudsen	cknudsen@radix.net
//			Fixed bug where cols was not init to 0
//		26-Jul-1999	Craig Knudsen	cknudsen@radix.net
//			Added 3D look
//			Added -html option to write HTML image maps
//		23-Jul-1999	Craig Knudsen	cknudsen@radix.net
//			Created

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/craigk5n/ilibgo"
	font "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)

const topShadow uint8 = 232
const bottomShandow uint8 = 112
const background uint8 = 192

func printUsage() {
	fmt.Printf("Usage:\n  index [options] outfile infile1 infile2 ...\n")
	fmt.Printf("\nOptions:\n")
	fmt.Printf("\t-h       show this help information\n")
	fmt.Printf("\t-f       overwrite outfile if it already exists\n")
	fmt.Printf("\t-w N     use pixel width of N for icons\n")
	fmt.Printf("\t-H N     use pixel height of N for icons\n")
	fmt.Printf("\t-s       run silently (no output to stdout)\n")
	fmt.Printf("\t-html f  specify an output file to write an HTML client-side image map\n")
	os.Exit(1)
}

func fatalError(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func fileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func main() {
	force := false
	silent := false
	iconW := 80
	iconH := 60
	outfile := ""
	infiles := make([]string, 0)
	format := ilibgo.FormatPNG
	border := 5
	textspace := 12
	var mapFp *os.File
	var outFp *os.File
	var err error

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		//fmt.Printf("Arg %d: %s\n", i, arg)
		if arg == "-h" || arg == "-help" {
			printUsage()
		} else if arg == "-version" {
			fmt.Printf("Ilib version %s\n", ilibgo.IlibVersion)
			os.Exit(0)
		} else if arg == "-html" {
			if i+1 >= len(os.Args) {
				fatalError("Parameter -html requres a filename")
			}
			i++
			if fileExists(os.Args[i]) && !force {
				fmt.Printf("Output file %s already exists\n", os.Args[i])
				os.Exit(1)
			}
			mapFp, err = os.Create(os.Args[i])
			if err != nil {
				fmt.Printf("Error writing to file %s: %v\n", os.Args[i], err)
				os.Exit(1)
			}
			defer mapFp.Close()
		} else if arg == "-w" {
			if i+1 >= len(os.Args) {
				fatalError("Parameter -w requres an integer parameter")
			}
			i++
			iconW, err = strconv.Atoi(os.Args[i])
			if err != nil {
				fatalError("Invalid parameter for -w")
			}
		} else if arg == "-H" {
			if i+1 >= len(os.Args) {
				fatalError("Parameter -H requres an integer parameter")
			}
			i++
			iconH, err = strconv.Atoi(os.Args[i])
			if err != nil {
				fatalError("Invalid parameter for -h")
			}
		} else if arg == "-f" {
			force = true
		} else if arg == "-s" {
			silent = true
		} else if len(outfile) == 0 {
			outfile = os.Args[i]
			if fileExists(outfile) && !force {
				fatalError("Output file already exists.  Use -f to force overwrite.")
			}
			outFp, err = os.Create(outfile)
			if err != nil {
				fatalError(fmt.Sprintf("Error writing output file %s: %v\n", outfile, err))
			}
		} else {
			// Must be an input image file
			if !fileExists(os.Args[i]) {
				fatalError(fmt.Sprintf("Input file does not exist: %s", os.Args[i]))
			}
			infiles = append(infiles, os.Args[i])
		}
	}

	if len(outfile) == 0 || len(infiles) == 0 {
		printUsage()
	}

	/* load font */
	font, _ := ilibgo.LoadFontFromData("helvR08", font.Font_helvR08())

	gc := ilibgo.CreateGraphicsContext()

	/* determine size of output image */
	cols := 0
	for cols = 0; cols*cols < len(infiles); cols++ {
	}
	rows := len(infiles) / cols
	if len(infiles) > (rows * cols) {
		rows++
	}
	w := border + cols*(iconW+border)
	h := border + rows*(iconH+border+textspace)

	bg, _ := ilibgo.AllocColor(background, background, background)
	outImage := ilibgo.CreateImageWithBackground(w, h, bg)
	black, _ := ilibgo.AllocColor(0, 0, 0)
	ts, _ := ilibgo.AllocColor(topShadow, topShadow, topShadow)
	bs, _ := ilibgo.AllocColor(bottomShandow, bottomShandow, bottomShandow)
	ilibgo.SetForeground(&gc, bg)
	ilibgo.FillRectangle(outImage, gc, 0, 0, w, h)

	ilibgo.SetForeground(&gc, black)
	ilibgo.SetFont(&gc, font)

	if mapFp != nil {
		fmt.Fprintf(mapFp,
			"<html><head><title>Image Index</title></head>\n")
		fmt.Fprintf(mapFp, "<body bgcolor=\"#%02x%02x%02x\"><h2>Image Index</h2>\n",
			background, background, background)
		fmt.Fprintf(mapFp, "<map name=\"image_index\">\n")
	}
	row := 0
	col := 0
	for i := 0; i < len(infiles); i++ {
		filename := filepath.Base(infiles[i])
		_, err := ilibgo.FileType(infiles[i])
		if err != nil {
			fmt.Printf("Uknown file type file %s: %v\n", infiles[i], err)
			continue
		}
		fp, err := os.Open(infiles[i])
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", infiles[i], err)
			continue
		}
		img, err := ilibgo.ReadImageFile(fp)
		fp.Close()
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", infiles[i], err)
			continue
		}
		x := border + col*(iconW+border)
		y := border + row*(iconH+border+textspace)
		/* determine scaled size.  maintain proportions */
		scale := 1
		thisW := ilibgo.ImageWidth(img) / scale
		thisH := ilibgo.ImageHeight(img) / scale
		for thisW > iconW || thisH > iconH {
			scale++
			thisW = ilibgo.ImageWidth(img) / scale
			thisH = ilibgo.ImageHeight(img) / scale
		}
		ilibgo.SetForeground(&gc, ts)
		ilibgo.DrawLine(outImage, gc, x-1, y-1, x+iconW+1, y-1)
		ilibgo.DrawLine(outImage, gc, x-1, y-1, x-1, y+iconH+1)
		ilibgo.SetForeground(&gc, bs)
		ilibgo.DrawLine(outImage, gc, x-1, y+iconH+1, x+iconW+1, y+iconH+1)
		ilibgo.DrawLine(outImage, gc, x+iconW+1, y-1, x+iconW+1, y+iconH+1)
		ilibgo.SetForeground(&gc, black)
		ilibgo.CopyImageScaled(img, outImage,
			0, 0, ilibgo.ImageWidth(img), ilibgo.ImageHeight(img),
			x+(iconW-thisW)/2, y+(iconH-thisH)/2, thisW, thisH)
		/* don't write more text than will fit under the image. */
		l := len(filename)
		name := filename
		strw, _ := ilibgo.TextWidth(gc, font, filename)
		for ; strw > iconW && len(name) > 0; l-- {
			name = filename[0 : l-1]
			strw, _ = ilibgo.TextWidth(gc, font, name)
		}
		ilibgo.DrawString(outImage, gc, x, y+iconH+textspace-2,
			name)
		if mapFp != nil {
			fmt.Fprintf(mapFp,
				"<area shape=\"rect\" coords=\"%d,%d,%d,%d\" href=\"%s\">\n",
				x, y, x+iconW, y+iconH, infiles[i])
		}
		col++
		if col >= cols {
			col = 0
			row++
		}
	}
	if mapFp != nil {
		fmt.Fprintf(mapFp,
			"<img src=\"%s\" width=\"%d\" height=\"%d\" usemap=\"#image_index\" border=\"0\">\n",
			outfile, ilibgo.ImageWidth(outImage), ilibgo.ImageHeight(outImage))
		fmt.Fprintf(mapFp, "</map>\n")
		fmt.Fprintf(mapFp,
			"<p><font size=\"-1\">Generated with Ilib v%s.</font>\n",
			ilibgo.IlibVersion)
		fmt.Fprintf(mapFp, "</body></html>\n")
		mapFp.Close()
	}

	if !silent {
		fmt.Printf("Writing %dx%d image to %s.\n", ilibgo.ImageWidth(outImage),
			ilibgo.ImageHeight(outImage), outfile)
	}
	ilibgo.WriteImageFile(outFp, outImage, format)
	outFp.Close()
}
