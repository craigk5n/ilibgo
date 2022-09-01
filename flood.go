package ilibgo

// History:
// 08-Jul-2022 Craig Knudsen craig@k5n.us
//             Converted from C to Go
// 25-Oct-2004 Craig Knudsen   cknudsen@cknudsen.com
//             Created

func colorsMatch(color1 Color, color2 Color) bool {
	// Just ignore alpha
	r1, g1, b1, _ := color1.color.RGBA()
	r2, g2, b2, _ := color2.color.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2
}

func FloodFill(image *Image, gc GraphicsContext, x int, y int) error {
	/* Get the color we are replacing with the flood fill */
	origColor := GetPoint(image, x, y)

	/* find left side, filling along the way */
	fillL := x
	inLine := true
	for inLine {
		SetPoint(image, gc, fillL, y)
		fillL--
		color := GetPoint(image, fillL, y)

		if fillL < 0 {
			inLine = false
		} else {
			inLine = colorsMatch(color, origColor)
		}
	}
	fillL++
	fillR := x
	/* find right side, filling along the way */
	inLine = true
	for inLine {
		SetPoint(image, gc, fillR, y)
		fillR++
		color := GetPoint(image, fillR, y)
		if fillR >= image.width {
			inLine = false
		} else {
			inLine = colorsMatch(color, origColor)
		}
	}
	fillR--

	/* search top and bottom */
	for i := fillL; i <= fillR; i++ {
		color := GetPoint(image, i, y-1)
		if y > 0 && colorsMatch(color, origColor) {
			FloodFill(image, gc, i, y-1)
		}
		color = GetPoint(image, i, y+1)
		if y < image.height && colorsMatch(color, origColor) {
			FloodFill(image, gc, i, y+1)
		}
	}

	return nil
}
