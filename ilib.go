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
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
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

// Encode writes img to w in the given format. It is the streaming core of the
// package's output path: w may be any io.Writer (a file, bytes.Buffer, HTTP
// response, gzip stream, etc.). The caller owns w and is responsible for
// closing it.
func Encode(w io.Writer, img *Image, format ImageFormat) error {
	switch format {
	case FormatPNG:
		return png.Encode(w, img.data)
	case FormatJPEG:
		// TODO: Allow user to specify quality
		options := jpeg.Options{Quality: jpeg.DefaultQuality}
		return jpeg.Encode(w, img.data, &options)
	case FormatGIF:
		options := gif.Options{NumColors: 256}
		return gif.Encode(w, img.data, &options)
	case FormatBMP:
		return bmp.Encode(w, img.data)
	case FormatTIFF:
		// x/image/tiff's encoder only supports Uncompressed and Deflate;
		// LZW is decode-only and returns "unsupported compression".
		options := tiff.Options{Compression: tiff.Deflate}
		return tiff.Encode(w, img.data, &options)
	case FormatPPM:
		return ppm.Encode(w, img.data)
	case FormatPGM:
		// TODO: Can we force PGM here?
		return errors.New("writing PGM not yet supported")
	case FormatPBM:
		// TODO: Can we force PBM here?
		return errors.New("writing PBM not yet supported")
	case FormatXPM:
		// TODO
		return errors.New("writing XPM not yet supported")
	}
	return fmt.Errorf("ilibgo: unknown image format %d", format)
}

// Decode reads an image from r (any io.Reader) and returns it as an *Image.
// The format is detected automatically from the stream contents via the
// registered image decoders. The caller owns r.
func Decode(r io.Reader) (*Image, error) {
	// image.Decode handles all registered image types
	var ret Image
	img, _, err := image.Decode(r)
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

// WriteImageFile writes an image to an already-open file (or any io.Writer).
// The writer is left open for the caller to close. It is a thin wrapper over
// Encode.
func WriteImageFile(f io.Writer, img *Image, format ImageFormat) error {
	return Encode(f, img, format)
}

// ReadImageFile creates an image from an already-open file (or any io.Reader).
// The reader is left open for the caller to close. It is a thin wrapper over
// Decode.
func ReadImageFile(f io.Reader) (*Image, error) {
	return Decode(f)
}

// SaveImageFile encodes img to the file at path, creating or truncating it.
// It is a convenience wrapper that opens the file, encodes, and closes it,
// joining any encode and close errors so a failed flush is not lost.
func SaveImageFile(path string, img *Image, format ImageFormat) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	encErr := Encode(f, img, format)
	closeErr := f.Close()
	return errors.Join(encErr, closeErr)
}

// LoadImageFile opens the image file at path, decodes it, and closes the file.
func LoadImageFile(path string) (*Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(f)
}
