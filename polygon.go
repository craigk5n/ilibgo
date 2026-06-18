package ilibgo

// History:
// 08-Jul-2022	Craig Knudsen craig@k5n.us
//            	Converted from C to Go
// 222-Nov-1999	Craig Knudsen cknudsen@cknudsen.com
//            	Created

type lineType struct {
	x1, y1, x2, y2 int
	slope          float64
}

// Calculate slope.  Keep in mind that our y coordinate is reverse
// compared to the standard math coordinate system.
func setLineSlope(line *lineType) {
	if line.x2 == line.x1 {
		line.slope = 0
	} else {
		line.slope = float64(line.y2-line.y1) /
			float64(line.x2-line.x1)
	}
}

// Does a line include the specified y value?
func lineIncludesYValue(line lineType, yval int) bool {
	if line.y1 <= yval && line.y2 >= yval {
		return true
	}
	if line.y2 <= yval && line.y1 >= yval {
		return true
	}
	return false
}

func getIntersectionXValue(line lineType, yval int) int {
	if line.x1 == line.x2 {
		return line.x1
	} else if line.y1 == line.y2 {
		return line.x1
	} else {
		/* calc b now
		 ** y = mx + b
		 ** b = y - mx
		 */
		b := float64(line.y2) - (line.slope * float64(line.x2))

		/*
		 ** now determine x value
		 ** x = (y - b) / m)
		 */
		ret := int((float64(yval) - b) / line.slope)
		return ret
	}
}

func (image *Image) DrawPolygon(gc GraphicsContext, points []Point) error {
	for loop := 1; loop < len(points); loop++ {
		image.DrawLine(gc, points[loop-1].X, points[loop-1].Y,
			points[loop].X, points[loop].Y)
	}

	return nil
}

// Fill a polygonn
func (image *Image) FillPolygon(gc GraphicsContext, points []Point) error {
	gc.lineWidth = 1

	// create an array of lines
	var lines []lineType = make([]lineType, 0)
	for loop := 1; loop < len(points); loop++ {
		var line lineType = lineType{x1: points[loop-1].X, y1: points[loop-1].Y, x2: points[loop].X, y2: points[loop].Y}
		setLineSlope(&line)
		lines = append(lines, line)
	}
	// last line connects first and last points
	var line lineType = lineType{x1: points[0].X, y1: points[0].Y, x2: points[len(points)-1].X, y2: points[len(points)-1].Y}
	setLineSlope(&line)
	lines = append(lines, line)

	/* debugging code
	   for ( loop = 0; loop < nlines; loop++ ) {
	     printf ( "Line %d: (%d,%d) to (%d,%d) with slope = %.2f\n",
	       loop, lines[loop].x1, lines[loop].y1, lines[loop].x2, lines[loop].y2,
	       (float)lines[loop].slope );
	   }
	*/

	// calculate the min and max y values
	minY := points[0].Y
	maxY := minY
	for loop := 1; loop < len(points); loop++ {
		minY = min(points[loop].Y, minY)
		maxY = max(points[loop].Y, maxY)
	}

	// now loop through from lowest y to top y
	for yloop := minY; yloop <= maxY; yloop++ {
		// now determine which lines of this polygon intersect this y value
		found := 0
		var left, right int
		for loop := 0; loop < len(lines); loop++ {
			if lineIncludesYValue(lines[loop], yloop) {
				// don't know if this is the left-most or right-most if this is
				// first point found...
				if found == 0 {
					// get intersection of this line and y
					if lines[loop].y1 == lines[loop].y2 {
						left = min(lines[loop].x1, lines[loop].x2)
						right = max(lines[loop].x1, lines[loop].x2)
					} else {
						left = getIntersectionXValue(lines[loop], yloop)
						right = left
					}
					found++
				} else {
					if lines[loop].y1 == lines[loop].y2 {
						left = min(left, lines[loop].x1)
						left = min(left, lines[loop].x2)
						right = max(right, lines[loop].x1)
						right = max(right, lines[loop].x2)
					} else {
						xval := getIntersectionXValue(lines[loop], yloop)
						left = min(left, xval)
						right = max(right, xval)
					}
					found++
					/* printf ( "left = %d, right = %d\n", left, right ); */
				}
			}
		}

		if found >= 1 {
			// Spans both the multi-intersection case and a single
			// intersection (point or horizontal edge at this row).
			image.DrawLine(gc, left, yloop, right, yloop)
		}
	}

	return nil
}
