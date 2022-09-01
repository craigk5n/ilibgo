package ilibgo

// Set the point at the specified location using the foreground
// color of the IGC parameter. (This is an alias to ISetPoint.)
func DrawPoint(image *Image, gc GraphicsContext, x int, y int) error {
	return SetPoint(image, gc, x, y)
}

// Set the point at the specified location using the foreground
// color of the IGC parameter.
func SetPoint(image *Image, gc GraphicsContext, x int, y int) error {
	image.data.Set(x, y, gc.foreground.color)
	return nil
}

// Get the color at a specific point.
func GetPoint(image *Image, x int, y int) Color {
	var ret Color
	ret.color = image.data.RGBAAt(x, y)
	return ret
}
