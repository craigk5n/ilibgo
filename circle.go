package ilibgo

// DrawCircle draws the outline of a circle.
func (image *Image) DrawCircle(gc GraphicsContext, x int, y int, r int) error {
	return image.DrawArc(gc, x, y, r, r, 0.0, 360.0)
}

// FillCircle draws a filled circle.
func (image *Image) FillCircle(gc GraphicsContext, x int, y int, r int) error {
	return image.FillArc(gc, x, y, r, r, 0.0, 360.0)
}

// DrawEllipse draws the outline of an ellipse with radii r1 and r2. A single
// full sweep is used so anti-aliased drawing uses the dedicated ellipse
// rasterizer (the result is the same full ellipse either way).
func (image *Image) DrawEllipse(gc GraphicsContext, x int, y int, r1 int, r2 int) error {
	return image.DrawArc(gc, x, y, r1, r2, 0.0, 360.0)
}

// FillEllipse draws a filled ellipse with radii r1 and r2.
func (image *Image) FillEllipse(gc GraphicsContext, x int, y int, r1 int, r2 int) error {
	return image.FillArc(gc, x, y, r1, r2, 0.0, 360.0)
}
