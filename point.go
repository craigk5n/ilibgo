package ilibgo

// Pt is a convenience constructor for a Point at (x, y).
func Pt(x int, y int) Point {
	return Point{X: x, Y: y}
}

// DrawPoint sets the point at (x, y) using the graphics context's foreground
// color. (Alias for SetPoint.)
func (image *Image) DrawPoint(gc GraphicsContext, x int, y int) error {
	return image.SetPoint(gc, x, y)
}

// SetPoint sets the point at (x, y) using the graphics context's foreground
// color.
func (image *Image) SetPoint(gc GraphicsContext, x int, y int) error {
	image.data.Set(x, y, gc.foreground.color)
	return nil
}

// GetPoint returns the color at (x, y).
func (image *Image) GetPoint(x int, y int) Color {
	var ret Color
	ret.color = image.data.RGBAAt(x, y)
	return ret
}
