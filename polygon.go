package ilibgo

import "fmt"

// History:
// 08-Jul-2022	Craig Knudsen craig@k5n.us
//            	Converted from C to Go
// 222-Nov-1999	Craig Knudsen cknudsen@cknudsen.com
//            	Created

type lineType struct {
	x1, y1, x2, y2 int
	slope          float64
}

func max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func min(a int, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
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

func DrawPolygon(image *Image, gc GraphicsContext, points []Point) error {
	for loop := 1; loop < len(points); loop++ {
		DrawLine(image, gc, points[loop-1].x, points[loop-1].y,
			points[loop].x, points[loop].y)
	}

	return nil
}

// Fill a polygonn
func FillPolygon(image *Image, gc GraphicsContext, points []Point) error {
	gc.lineWidth = 1

	// create an array of lines
	var lines []lineType = make([]lineType, 0)
	for loop := 1; loop < len(points); loop++ {
		var line lineType = lineType{x1: points[loop-1].x, y1: points[loop-1].y, x2: points[loop].x, y2: points[loop].y}
		setLineSlope(&line)
		lines = append(lines, line)
	}
	// last line connects first and last points
	var line lineType = lineType{x1: points[0].x, y1: points[0].y, x2: points[len(points)-1].x, y2: points[len(points)-1].y}
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
	minY := points[0].y
	maxY := minY
	for loop := 1; loop < len(points); loop++ {
		minY = min(points[loop].y, minY)
		maxY = max(points[loop].y, maxY)
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

		if found >= 2 {
			DrawLine(image, gc, left, yloop, right, yloop)
		} else if found == 1 && left == right {
			DrawLine(image, gc, left, yloop, right, yloop)
		} else if found == 1 && left != right {
			/* Eek.  This really shouldn't happen */
			return fmt.Errorf("ilib bug no. 34534345")
		}
	}

	return nil
}
