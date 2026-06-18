package ilibgo

import (
	"image"
	"image/draw"
)

// DrawRectangle draws the outline of a rectangle.
func (image *Image) DrawRectangle(gc GraphicsContext, x int, y int, width int, height int) error {
	image.DrawLine(gc, x, y, x+width, y)
	image.DrawLine(gc, x, y, x, y+height)
	image.DrawLine(gc, x, y+height, x+width, y+height)
	image.DrawLine(gc, x+width, y, x+width, y+height)

	return nil
}

// FillRectangle draws a filled rectangle. The rectangle is clipped to the
// image bounds, so off-image or negative coordinates are handled safely.
func (img *Image) FillRectangle(gc GraphicsContext, x int, y int, width int, height int) error {
	r := image.Rect(x, y, x+width, y+height)
	draw.Draw(img.data, r, &image.Uniform{C: gc.foreground.color}, image.Point{}, draw.Src)
	return nil
}
