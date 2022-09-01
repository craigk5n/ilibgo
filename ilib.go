// Image Library for Go
//
// Ilib is a library (and some tools and examples) written in Go
// that can read, create, manipulate and save images.  It is capable
// of using [X11 BDF fonts] for drawing text.  That means you get
// lots of fonts to use.  You can even create your
// own if you know how to create an X11 BDF font.  It should be able
// to read any image file that the base Go image package supports.
// Copyright (C) 2001-2022 Craig Knudsen, craig@k5n.us
// http://github.com/craigk5n/ilibgo
//
// [X11 BDF fonts]: https://gitlab.freedesktop.org/xorg/font
package ilibgo

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path/filepath"

	"github.com/lmittmann/ppm"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	"os"
	"strings"
)

type Point struct {
	x int
	y int
}

type Color struct {
	color color.RGBA
}

type ImageOption uint8

const (
	OptionNone       ImageOption = 0
	OptionGrayscale  ImageOption = 1
	OptionGreyscale  ImageOption = OptionGrayscale
	OptionAscii      ImageOption = 2 // ascii output for pbm/pgm/ppm images
	OptionInterlaced ImageOption = 4 // interlaced output (GIF)
)

// Default color values
const BlackPixel int = 0
const WhitePixel int = 1

type Image struct {
	width    int
	height   int
	comments string // TODO: save off as metadata in image
	data     *image.RGBA
}

type Font struct {
	name string
	font *BdfFont
	// Add additional support font types (truetype, etc.) here
}

type GraphicsContext struct {
	foreground Color
	background Color
	font       *Font
	lineWidth  int
	lineStyle  LineStyle
	textStyle  TextStyle
}

// Check to see if the specified image format (as represented as the file suffix like "png")
// is supported.  Some image formats
func IsSupportedFormat(stringFormat string) bool {
	_, e := FormatStringToType(stringFormat)
	return e == nil
}

// Convert a string representation of an image format (typically the file extension like "png"),
// and convert it to the corresponding IImageFormat type.  An error is returned if the
// type is not recognized or not supported.
func FormatStringToType(formatString string) (ImageFormat, error) {
	formatString = strings.ReplaceAll(formatString, ".", "")
	switch strings.ToLower(formatString) {
	case "gif":
		return FormatGIF, nil
	case "ppm":
		return FormatPPM, nil
	case "pgm":
		return FormatPGM, nil
	case "pbm":
		return FormatPBM, nil
	case "xpm":
		return FormatXPM, nil
	case "xbm":
		return FormatXBM, nil
	case "png":
		return FormatPNG, nil
	case "jpeg":
		return FormatJPEG, nil
	case "jpg":
		return FormatJPEG, nil
	case "bmp":
		return FormatBMP, nil
	case "tiff":
		return FormatTIFF, nil
	}
	// Default to PPM and return error.
	return FormatPPM, errors.New("Unrecognized file type '" + formatString + "'")
}

func FileType(filename string) (ImageFormat, error) {
	ext := filepath.Ext(filename)
	return FormatStringToType(ext)
}

func ImageWidth(img *Image) int {
	return img.width
}

func ImageHeight(img *Image) int {
	return img.height
}

func CreateImage(width int, height int) *Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var ret Image = Image{data: img, width: width, height: height}
	return &ret
}

func CreateImageWithBackground(width int, height int, background Color) *Image {
	ret := CreateImage(width, height)
	// TODO: Use IFillRect
	gc := CreateGraphicsContext()
	gc.foreground = background
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			SetPoint(ret, gc, x, y)
		}
	}
	return ret
}

// Writes an image to a file.
// The file is left open for the caller to close.
func WriteImageFile(f *os.File, img *Image, format ImageFormat) error {
	var err error
	defer f.Close()
	switch format {
	case FormatPNG:
		png.Encode(f, img.data)
	case FormatJPEG:
		// TODO: Allow user to specify quality
		var options jpeg.Options = jpeg.Options{Quality: jpeg.DefaultQuality}
		jpeg.Encode(f, img.data, &options)
	case FormatGIF:
		var options gif.Options = gif.Options{NumColors: 256}
		gif.Encode(f, img.data, &options)
	case FormatBMP:
		bmp.Encode(f, img.data)
	case FormatTIFF:
		var tiffOptions tiff.Options = tiff.Options{Compression: tiff.LZW}
		tiff.Encode(f, img.data, &tiffOptions)
	case FormatPPM:
		ppm.Encode(f, img.data)
	case FormatPGM:
		// TODO: Can we force PGM here?
		err = errors.New("writing PGM not yet supported")
	case FormatPBM:
		// TODO: Can we force PBM here?
		err = errors.New("writing PBM not yet supported")
	case FormatXPM:
		// TODO
		err = errors.New("writing XPM not yet supported")
	}
	return err
}

// Creates an image from an image file.
// The file is left open for the caller to close.
func ReadImageFile(f *os.File) (*Image, error) {
	// image.Decode handles all registered image types
	var ret Image
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	b := img.Bounds()
	ret.width = b.Size().X
	ret.height = b.Size().Y
	ret.data = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(ret.data, ret.data.Bounds(), img, b.Min, draw.Src)

	return &ret, nil
}
