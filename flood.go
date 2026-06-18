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

// FloodFill replaces the contiguous region of pixels matching the color at
// (x, y) with the graphics context's foreground color. It uses an explicit
// scanline stack rather than recursion so that large regions cannot overflow
// the goroutine stack.
func (image *Image) FloodFill(gc GraphicsContext, x int, y int) error {
	if x < 0 || x >= image.width || y < 0 || y >= image.height {
		return nil
	}

	/* The color we are replacing with the flood fill */
	origColor := image.GetPoint(x, y)

	// Nothing to do if the seed already has the fill color; filling anyway
	// would never stop matching origColor and could loop forever.
	if colorsMatch(origColor, gc.foreground) {
		return nil
	}

	stack := []Point{{X: x, Y: y}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Skip pixels already recolored (or never matching) since being queued.
		if !colorsMatch(image.GetPoint(p.X, p.Y), origColor) {
			continue
		}

		/* find left side of this run */
		fillL := p.X
		for fillL-1 >= 0 && colorsMatch(image.GetPoint(fillL-1, p.Y), origColor) {
			fillL--
		}
		/* find right side of this run */
		fillR := p.X
		for fillR+1 < image.width && colorsMatch(image.GetPoint(fillR+1, p.Y), origColor) {
			fillR++
		}

		/* fill the run and queue matching pixels above and below */
		for i := fillL; i <= fillR; i++ {
			image.SetPoint(gc, i, p.Y)
			if p.Y-1 >= 0 && colorsMatch(image.GetPoint(i, p.Y-1), origColor) {
				stack = append(stack, Point{X: i, Y: p.Y - 1})
			}
			if p.Y+1 < image.height && colorsMatch(image.GetPoint(i, p.Y+1), origColor) {
				stack = append(stack, Point{X: i, Y: p.Y + 1})
			}
		}
	}

	return nil
}
