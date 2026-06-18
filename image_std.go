package ilibgo

import (
	"image"
	"image/color"
	"image/draw"
)

// Image satisfies the standard library's image.Image and draw.Image
// interfaces by delegating to its backing *image.RGBA. This lets callers pass
// an *Image directly to image/draw, the image encoders (image/png, etc.),
// golang.org/x/image scalers, and anything else accepting those interfaces.
var (
	_ image.Image = (*Image)(nil)
	_ draw.Image  = (*Image)(nil)
)

// ColorModel returns the image's color model (color.RGBAModel).
func (img *Image) ColorModel() color.Model { return img.data.ColorModel() }

// Bounds returns the domain over which the image's pixels are defined.
func (img *Image) Bounds() image.Rectangle { return img.data.Bounds() }

// At returns the color of the pixel at (x, y).
func (img *Image) At(x int, y int) color.Color { return img.data.At(x, y) }

// Set sets the pixel at (x, y) to c, satisfying draw.Image.
func (img *Image) Set(x int, y int, c color.Color) { img.data.Set(x, y, c) }
