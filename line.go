package ilibgo

import (
	"math"
)

// History:
// 08-Jul-2022 Craig Knudsen craig@k5n.us
//             Converted from C to Go
// 20-May-1996 Craig Knudsen   cknudsen@cknudsen.com
//             Created

const onOffPixels float64 = 3.0

func DrawLine(image *Image, gc GraphicsContext, x1 int, y1 int, x2 int, y2 int) error {
	var myx, myy int
	var slope, myslope, curx, cury float64
	done := false
	var temp int
	drawCount := 0.0
	onOffSize := 0.0

	/* x2 should always be greater than x1 */
	if x2 < x1 {
		temp = x2
		x2 = x1
		x1 = temp
		temp = y2
		y2 = y1
		y1 = temp
	}

	/* remember, our coordinate system is reversed for y */
	if x1 == x2 {
		if y2 < y1 {
			// swap y1 & y2
			temp = x2
			x2 = x1
			x1 = temp
			temp = y2
			y2 = y1
			y1 = temp
		}
	} else {
		slope = -(float64(y2) - float64(y1)) / (float64(x1) - float64(x2))
	}
	curx = float64(x1)
	cury = float64(y1)

	// handle dashes
	if gc.lineStyle == LineOnOffDash {
		if x1 == x2 || math.Abs(slope) < 0.1 || math.Abs(slope) > 10.0 {
			onOffSize = onOffPixels
		} else {
			myslope = math.Abs(slope)
			if myslope > 1.0 {
				myslope = 1.0 / myslope
			}
			// myslope now between 0 and 1.0
			onOffSize = onOffPixels + (0.41 * myslope)
		}
	}

	for !done {
		if x1 == x2 {
			// Slope is infinite
			if cury >= float64(y2) {
				done = true
			} else {
				cury += 1.0
			}
		} else if slope > 1.0 {
			if cury >= float64(y2) {
				done = true
			} else {
				cury += 1.0
				curx += 1.0 / slope
			}
		} else if slope < -1.0 {
			if cury <= float64(y2) {
				done = true
			} else {
				cury -= 1.0
				curx -= 1.0 / slope
			}
		} else if slope >= 0.0 {
			if curx >= float64(x2) {
				done = true
			} else {
				curx += 1.0
				cury += slope
			}
		} else if slope < 0.0 {
			if curx >= float64(x2) {
				done = true
			} else {
				curx += 1.0
				cury += slope
			}
		}

		if gc.lineStyle == LineOnOffDash {
			drawCount += 1.0
			if (int(math.Floor(drawCount/onOffSize)) % 2) == 1 {
				continue
			}
		}

		if !done {
			myx = int(curx)
			myy = int(cury)
			switch gc.lineWidth {
			default:
			case 0:
				DrawPoint(image, gc, myx, myy)
			case 1:
				DrawPoint(image, gc, myx, myy)
			case 2:
				DrawPoint(image, gc, myx, myy)
				DrawPoint(image, gc, myx-1, myy)
				DrawPoint(image, gc, myx-1, myy-1)
				DrawPoint(image, gc, myx, myy-1)
			case 3:
				DrawPoint(image, gc, myx, myy)
				DrawPoint(image, gc, myx-1, myy)
				DrawPoint(image, gc, myx+1, myy-1)
				DrawPoint(image, gc, myx, myy-1)
				DrawPoint(image, gc, myx, myy+1)
			}
		}
	}

	return nil
}
