package ilibgo

// Draw a circle
func DrawCircle(image *Image, gc GraphicsContext, x int, y int, r int) error {
	return DrawArc(image, gc, x, y, r, r, 0.0, 360.0)
}

// Fill a circle
func FillCircle(image *Image, gc GraphicsContext, x int, y int, r int) error {
	return FillArc(image, gc, x, y, r, r, 0.0, 360.0)
}

func DrawEllipse(image *Image, gc GraphicsContext, x int, y int, r1 int, r2 int) error {
	ret := DrawArc(image, gc, x, y, r1, r2, 0.0, 90.0)
	if ret != nil {
		return ret
	}
	return DrawArc(image, gc, x, y, r1, r2, 90.0, 360.0)
}
