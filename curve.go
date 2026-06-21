package ilibgo

import (
	"fmt"
	"math"
)

// Smooth curves: cubic Bezier paths and Catmull-Rom splines. Both are
// flattened to short line segments and drawn with DrawLine, so they inherit the
// graphics context's line style and anti-aliasing. Ported from IDrawCurve.c.

func ptDist(a, b Point) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// curveSamples returns the number of flattening samples for a segment of the
// given approximate length: ~1 sample every 3 pixels, bounded so degenerate
// input stays cheap.
func curveSamples(length float64) int {
	n := int(length / 3.0)
	if n < 4 {
		n = 4
	}
	if n > 4096 {
		n = 4096
	}
	return n
}

// DrawBezier draws a chain of cubic Bezier segments. points[0] is the start;
// every following group of three points is (control1, control2, end) of one
// cubic, so len(points) must be 4, 7, 10, ... (1 + 3k).
func (image *Image) DrawBezier(gc GraphicsContext, points []Point) error {
	n := len(points)
	if n < 4 || (n-1)%3 != 0 {
		return fmt.Errorf("ilibgo: DrawBezier: need 1+3k points (4,7,10,...), got %d", n)
	}

	for seg := 0; seg+3 < n; seg += 3 {
		p0, p1 := points[seg], points[seg+1]
		p2, p3 := points[seg+2], points[seg+3]
		samples := curveSamples(ptDist(p0, p1) + ptDist(p1, p2) + ptDist(p2, p3))
		prevx, prevy := p0.X, p0.Y
		for i := 1; i <= samples; i++ {
			t := float64(i) / float64(samples)
			mt := 1.0 - t
			b0 := mt * mt * mt
			b1 := 3 * mt * mt * t
			b2 := 3 * mt * t * t
			b3 := t * t * t
			x := int(b0*float64(p0.X) + b1*float64(p1.X) + b2*float64(p2.X) + b3*float64(p3.X) + 0.5)
			y := int(b0*float64(p0.Y) + b1*float64(p1.Y) + b2*float64(p2.Y) + b3*float64(p3.Y) + 0.5)
			image.DrawLine(gc, prevx, prevy, x, y)
			prevx, prevy = x, y
		}
	}
	return nil
}

// DrawSpline draws a Catmull-Rom spline passing through all the given points
// (the curve interpolates each point). Endpoints are clamped (the first/last
// point is duplicated). Requires at least 2 points; two points draw a straight
// line.
func (image *Image) DrawSpline(gc GraphicsContext, points []Point) error {
	n := len(points)
	if n < 2 {
		return fmt.Errorf("ilibgo: DrawSpline: need at least 2 points, got %d", n)
	}

	for seg := 0; seg < n-1; seg++ {
		var p0 Point
		if seg > 0 {
			p0 = points[seg-1]
		} else {
			p0 = points[0]
		}
		p1 := points[seg]
		p2 := points[seg+1]
		var p3 Point
		if seg+2 < n {
			p3 = points[seg+2]
		} else {
			p3 = points[n-1]
		}
		samples := curveSamples(ptDist(p1, p2) * 1.5)
		prevx, prevy := p1.X, p1.Y
		for i := 1; i <= samples; i++ {
			t := float64(i) / float64(samples)
			t2 := t * t
			t3 := t2 * t
			x := 0.5 * ((2.0 * float64(p1.X)) +
				(-float64(p0.X)+float64(p2.X))*t +
				(2.0*float64(p0.X)-5.0*float64(p1.X)+4.0*float64(p2.X)-float64(p3.X))*t2 +
				(-float64(p0.X)+3.0*float64(p1.X)-3.0*float64(p2.X)+float64(p3.X))*t3)
			y := 0.5 * ((2.0 * float64(p1.Y)) +
				(-float64(p0.Y)+float64(p2.Y))*t +
				(2.0*float64(p0.Y)-5.0*float64(p1.Y)+4.0*float64(p2.Y)-float64(p3.Y))*t2 +
				(-float64(p0.Y)+3.0*float64(p1.Y)-3.0*float64(p2.Y)+float64(p3.Y))*t3)
			ix := int(x + 0.5)
			iy := int(y + 0.5)
			image.DrawLine(gc, prevx, prevy, ix, iy)
			prevx, prevy = ix, iy
		}
	}
	return nil
}
