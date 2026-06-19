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
	"golang.org/x/image/font"
	"golang.org/x/image/tiff"

	"os"
	"strings"
)

// Point is an (X, Y) coordinate. Its fields are exported so callers can
// construct points directly (e.g. for DrawPolygon/FillPolygon), either with a
// struct literal (Point{X: 1, Y: 2}) or the Pt helper.
type Point struct {
	X int
	Y int
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

// Font is a loaded font. It is backed either by a bitmap BDF font (font != nil)
// or by a scalable TrueType/OpenType face at a fixed size (face != nil).
type Font struct {
	name   string
	font   *BdfFont  // BDF bitmap font (nil for TrueType)
	face   font.Face // TrueType/OpenType face (nil for BDF)
	height int       // cached pixel line height for face fonts
}

// isTrueType reports whether the font is backed by a scalable face.
func (f *Font) isTrueType() bool { return f != nil && f.face != nil }

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
	draw.Draw(ret.data, ret.data.Bounds(), &image.Uniform{C: background.color}, image.Point{}, draw.Src)
	return ret
}

// Writes an image to a file.
// The file is left open for the caller to close.
func WriteImageFile(f *os.File, img *Image, format ImageFormat) error {
	var err error
	switch format {
	case FormatPNG:
		err = png.Encode(f, img.data)
	case FormatJPEG:
		// TODO: Allow user to specify quality
		var options jpeg.Options = jpeg.Options{Quality: jpeg.DefaultQuality}
		err = jpeg.Encode(f, img.data, &options)
	case FormatGIF:
		var options gif.Options = gif.Options{NumColors: 256}
		err = gif.Encode(f, img.data, &options)
	case FormatBMP:
		err = bmp.Encode(f, img.data)
	case FormatTIFF:
		// x/image/tiff's encoder only supports Uncompressed and Deflate;
		// LZW is decode-only and returns "unsupported compression".
		var tiffOptions tiff.Options = tiff.Options{Compression: tiff.Deflate}
		err = tiff.Encode(f, img.data, &tiffOptions)
	case FormatPPM:
		err = ppm.Encode(f, img.data)
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
