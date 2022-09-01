package ilibgo

// Description:
//	Copy an area of an image to another image.
//
// History:
//  09-Aug-2022	Craig Knudsen craig@k5n.us
//      	Converted from C to go
//	15-Aug-2001	Craig Knudsen	cknudsen@cknudsen.com
//			Fixed bug in ICopyImageScaled
//			(thanks Gal Steinitz for this fix)
//	23-Jul-1999	Craig Knudsen   cknudsen@cknudsen.com
//			Added ICopyImageScaled
//	11-Nov-1998	Craig Knudsen	cknudsen@cknudsen.com
//			Allow transparent values to not be copied.
//	20-May-1996	Craig Knudsen	cknudsen@cknudsen.com
//			Created

// Copy a rectangle portion of a source image to a destination image
func CopyImage(source *Image, dest *Image, gc GraphicsContext, src_x int, src_y int, width int, height int,
	dest_x int, dest_y int) error {
	SetForeground(&gc, newIColor(0, 0, 0, 255))

	y := dest_y
	for row := src_y; row < src_y+height; row++ {
		x := dest_x
		for col := src_x; col < src_x+width; col++ {
			c := GetPoint(source, col, row)
			SetForeground(&gc, c)
			SetPoint(dest, gc, x, y)
			x++
		}
		y++
	}
	return nil
}

// This allows the user to scale up or down the source image onto
// the destination image.
func CopyImageScaled(source *Image, dest *Image,
	src_x int, src_y int, src_width int, src_height int,
	dest_x int, dest_y int, dest_width int, dest_height int) error {

	// When scaling down, we might want to add an algorithm for averaging
	// a series of source pixels into the destination pixel.  For now,
	// we just grab one color.
	var gc GraphicsContext
	scalex := float64(dest_width) / float64(src_width)
	scaley := float64(dest_height) / float64(src_height)
	for y := dest_y; y < dest_y+dest_height; y++ {
		for x := dest_x; x < dest_x+dest_width; x++ {
			// get location from source image for this location
			// x2,y2 is location in source image.
			tempx := float64(src_x) + float64(x-dest_x)/scalex
			x2 := int(tempx)
			tempy := float64(src_y) + float64(y-dest_y)/scaley
			y2 := int(tempy)
			gc.foreground = GetPoint(source, x2, y2)
			SetPoint(dest, gc, x, y)
		}
	}
	return nil
}
