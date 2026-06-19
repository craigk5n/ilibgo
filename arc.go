package ilibgo

import (
	"fmt"
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
func (image *Image) DrawArc(gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	if r1 < 0 || r2 < 0 {
		return fmt.Errorf("ilibgo: DrawArc: negative radius (r1=%d, r2=%d)", r1, r2)
	}
	var myx, myy, lastx, lasty, N, loop int

	/* because our y is upside down, make all angles their negative */
	a1 = 360 - a1
	a2 = 360 - a2

	N = int(math.Abs(a2-a1)) + 8
	a := a1 * 2.0 * math.Pi / 360.0
	da := (a2 - a1) * (2.0 * math.Pi / 360.0) / (float64(N) - 1.0)
	// Step the angle incrementally with a rotation matrix instead of calling
	// Cos/Sin every iteration (see IDEAS §5.5).
	cosCur, sinCur := math.Cos(a), math.Sin(a)
	cosDa, sinDa := math.Cos(da), math.Sin(da)
	for loop = 0; loop < N; loop++ {
		myx = x + int(float64(r1)*cosCur)
		myy = y + int(float64(r2)*sinCur)
		if loop > 0 {
			image.DrawLine(gc, lastx, lasty, myx, myy)
		}
		lastx = myx
		lasty = myy
		cosCur, sinCur = cosCur*cosDa-sinCur*sinDa, sinCur*cosDa+cosCur*sinDa
	}

	return nil
}

// Draw an arc and connect it to the center point.
func (image *Image) DrawEnclosedArc(gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	if r1 < 0 || r2 < 0 {
		return fmt.Errorf("ilibgo: DrawEnclosedArc: negative radius (r1=%d, r2=%d)", r1, r2)
	}
	// because our y is upside down, make all angles their negative
	a1 = 360 - a1
	a2 = 360 - a2

	N := int(math.Abs(a2-a1)) + 8
	a := a1 * 2.0 * math.Pi / 360.0
	da := float64(a2-a1) * (2.0 * math.Pi / 360.0) / float64(N-1)
	lastx := 0
	lasty := 0
	cosCur, sinCur := math.Cos(a), math.Sin(a)
	cosDa, sinDa := math.Cos(da), math.Sin(da)
	for loop := 0; loop < N; loop++ {
		myx := x + int(float64(r1)*cosCur)
		myy := y + int(float64(r2)*sinCur)
		if loop > 0 {
			image.DrawLine(gc, lastx, lasty, myx, myy)
		}
		if loop == N-1 || loop == 0 {
			image.DrawLine(gc, x, y, myx, myy)
		}
		lastx = myx
		lasty = myy
		cosCur, sinCur = cosCur*cosDa-sinCur*sinDa, sinCur*cosDa+cosCur*sinDa
	}

	return nil
}

// Fill an arc
func (image *Image) FillArc(gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	if r1 < 0 || r2 < 0 {
		return fmt.Errorf("ilibgo: FillArc: negative radius (r1=%d, r2=%d)", r1, r2)
	}
	var points []Point = make([]Point, 0)

	// because our y is upside down, make all angles their negative
	a1 = 360 - a1
	a2 = 360 - a2

	N := int(math.Abs(a2-a1)) + 9
	a := a1 * 2.0 * math.Pi / 360.0
	da := (a2 - a1) * (2.0 * math.Pi / 360.0) / float64(N-1)
	cosCur, sinCur := math.Cos(a), math.Sin(a)
	cosDa, sinDa := math.Cos(da), math.Sin(da)
	for loop := 0; loop < N; loop++ {
		var p Point
		p.X = x + int(float64(r1)*cosCur)
		p.Y = y + int(float64(r2)*sinCur)
		points = append(points, p)
		cosCur, sinCur = cosCur*cosDa-sinCur*sinDa, sinCur*cosDa+cosCur*sinDa
	}

	// if we're not drawing a circle, add in the center point
	if a2-a1 < 359.9 {
		var p Point
		p.X = x
		p.Y = y
		points = append(points, p)
	}
	return image.FillPolygon(gc, points)
}
