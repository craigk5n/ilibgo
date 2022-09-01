package ilibgo

// Draw a rectangle
func DrawRectangle(image *Image, gc GraphicsContext, x int, y int, width int, height int) error {
	DrawLine(image, gc, x, y, x+width, y)
	DrawLine(image, gc, x, y, x, y+height)
	DrawLine(image, gc, x, y+height, x+width, y+height)
	DrawLine(image, gc, x+width, y, x+width, y+height)

	return nil
}

// Draw a filled rectangle
func FillRectangle(image *Image, gc GraphicsContext, x int, y int, width int, height int) error {
	for row := y; row < y+height && row < image.height; row++ {
		for col := x; col < x+width && col < image.width; col++ {
			if row >= 0 && col >= 0 {
				SetPoint(image, gc, col, row)
			}
		}
	}

	return nil
}
