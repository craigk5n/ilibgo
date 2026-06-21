package ilibgo

import "math"

// ArcPoints holds the key coordinates of an arc computed by ArcProperties:
// the points on the ellipse at the start angle (A1) and end angle (A2), and the
// point at the angle halfway between them (Middle). These are handy for placing
// labels or connecting lines (e.g. pie-chart slice callouts).
type ArcPoints struct {
	A1X, A1Y         int // point at the first angle (a1)
	A2X, A2Y         int // point at the second angle (a2)
	MiddleX, MiddleY int // point at the midpoint angle (a1+a2)/2
}

// ArcProperties computes the start, end, and midpoint coordinates of the arc
// centered at (x, y) with radii r1 (x) and r2 (y) spanning angles a1..a2
// (degrees, 0..360). Mirrors the C library's IArcProperties; the angles use the
// same convention as DrawArc (the image y axis points down, so angles are
// internally negated).
func ArcProperties(x, y, r1, r2 int, a1, a2 float64) ArcPoints {
	const deg2rad = 2.0 * math.Pi / 360.0

	// because our y is upside down, make all angles their negative
	a1 = 360 - a1
	a2 = 360 - a2

	var p ArcPoints
	p.A1X = x + int(float64(r1)*math.Cos(deg2rad*a1))
	p.A1Y = y + int(float64(r2)*math.Sin(deg2rad*a1))
	p.A2X = x + int(float64(r1)*math.Cos(deg2rad*a2))
	p.A2Y = y + int(float64(r2)*math.Sin(deg2rad*a2))

	a := (a1 + a2) / 2
	p.MiddleX = x + int(float64(r1)*math.Cos(deg2rad*a))
	p.MiddleY = y + int(float64(r2)*math.Sin(deg2rad*a))
	return p
}
