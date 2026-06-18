package ilibgo

import (
	"image"
	"image/draw"
)

// Draw a rectangle
func DrawRectangle(image *Image, gc GraphicsContext, x int, y int, width int, height int) error {
	DrawLine(image, gc, x, y, x+width, y)
	DrawLine(image, gc, x, y, x, y+height)
	DrawLine(image, gc, x, y+height, x+width, y+height)
	DrawLine(image, gc, x+width, y, x+width, y+height)

	return nil
}

// Draw a filled rectangle. The rectangle is clipped to the image bounds, so
// off-image or negative coordinates are handled safely.
func FillRectangle(img *Image, gc GraphicsContext, x int, y int, width int, height int) error {
	r := image.Rect(x, y, x+width, y+height)
	draw.Draw(img.data, r, &image.Uniform{C: gc.foreground.color}, image.Point{}, draw.Src)
	return nil
}
