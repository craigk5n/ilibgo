package ilibgo

// Draw a circle
func IDrawCircle(image *Image, gc GraphicsContext, x int, y int, r int) error {
	return IDrawArc(image, gc, x, y, r, r, 0.0, 360.0)
}

// Fill a circle
func IFillCircle(image *Image, gc GraphicsContext, x int, y int, r int) error {
	return IFillArc(image, gc, x, y, r, r, 0.0, 360.0)
}

func IDrawEllipse(image *Image, gc GraphicsContext, x int, y int, r1 int, r2 int) error {
	ret := IDrawArc(image, gc, x, y, r1, r2, 0.0, 90.0)
	if ret != nil {
		return ret
	}
	return IDrawArc(image, gc, x, y, r1, r2, 90.0, 360.0)
}
