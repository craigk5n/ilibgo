package ilibgo

import (
	"math"
)

/*
 * History:
 *      08-Jul-2022     craig@k5n.us
 *                      Converted from C to Go
 *      28-Nov-99       Craig Knudsen   cknudsen@cknudsen.com
 *                      Added IDrawEnclosedArc()
 *      19-Nov-99       Craig Knudsen   cknudsen@cknudsen.com
 *                      Created
 *
 ****************************************************************************/

// Draw an arc.  Both arc1 and arc2 are in degrees from 0 to 360.
func IDrawArc(image *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	var myx, myy, lastx, lasty, N, loop int

	/* because our y is upside down, make all angles their negative */
	a1 = 360 - a1
	a2 = 360 - a2

	N = int(math.Abs(a2-a1)) + 8
	a := a1 * 2.0 * math.Pi / 360.0
	da := (a2 - a1) * (2.0 * math.Pi / 360.0) / (float64(N) - 1.0)
	for loop = 0; loop < N; loop++ {
		myx = x + int(float64(r1)*math.Cos(a+float64(loop)*da))
		myy = y + int(float64(r2)*math.Sin(a+float64(loop)*da))
		if loop > 0 {
			DrawLine(image, gc, lastx, lasty, myx, myy)
		}
		lastx = myx
		lasty = myy
	}

	return nil
}

// Draw an arc and connect it to the center point.
func IDrawEnclosedArc(image *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	// because our y is upside down, make all angles their negative
	a1 = 360 - a1
	a2 = 360 - a2

	N := int(math.Abs(a2-a1)) + 8
	a := a1 * 2.0 * math.Pi / 360.0
	da := float64(a2-a1) * (2.0 * math.Pi / 360.0) / float64(N-1)
	lastx := 0
	lasty := 0
	for loop := 0; loop < N; loop++ {
		myx := x + int(r1*int(math.Cos(a+float64(loop)*da)))
		myy := y + int(r2*int(math.Sin(a+float64(loop)*da)))
		if loop > 0 {
			DrawLine(image, gc, lastx, lasty, myx, myy)
		}
		if loop == N-1 || loop == 0 {
			DrawLine(image, gc, x, y, myx, myy)
		}
		lastx = myx
		lasty = myy
	}

	return nil
}

// Fill an arc
func IFillArc(image *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	var points []Point = make([]Point, 0)

	// because our y is upside down, make all angles their negative
	a1 = 360 - a1
	a2 = 360 - a2

	N := int(math.Abs(a2-a1)) + 9
	a := a1 * 2.0 * math.Pi / 360.0
	da := (a2 - a1) * (2.0 * math.Pi / 360.0) / float64(N-1)
	for loop := 0; loop < N; loop++ {
		var p Point
		p.x = x + int(float64(r1)*math.Cos(a+float64(loop)*da))
		p.y = y + int(float64(r2)*math.Sin(a+float64(loop)*da))
		points = append(points, p)
	}

	// if we're not drawing a circle, add in the center point
	if a2-a1 < 359.9 {
		var p Point
		p.x = x
		p.y = y
		points = append(points, p)
		FillPolygon(image, gc, points)
	} else {
		FillPolygon(image, gc, points)
	}

	return nil
}
