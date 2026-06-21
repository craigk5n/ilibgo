package ilibgo

import "math"

// Anti-aliased rasterizers used when a graphics context has anti-aliasing
// enabled (SetAntiAlias). Ported from the C library's IDrawLin.c, IDrawArc.c,
// IFillArc.c and IFillPol.c. All of them composite the foreground with
// fractional coverage via blendPoint.

const aaSupersample = 4 // NxN sub-samples per pixel for supersampled fills

func fpart(x float64) float64  { return x - math.Floor(x) }
func rfpart(x float64) float64 { return 1.0 - fpart(x) }

// aaPlot composites the foreground at (px,py) with fractional coverage c (0..1).
func (img *Image) aaPlot(gc GraphicsContext, px, py int, c float64) {
	if c <= 0.0 {
		return
	}
	cover := int(c*255.0 + 0.5)
	if cover <= 0 {
		return
	}
	if cover > 255 {
		cover = 255
	}
	img.blendPoint(gc, px, py, uint32(cover))
}

// drawLineAA draws Xiaolin Wu's anti-aliased line between integer endpoints.
// Used for width-1 solid lines when GC anti-aliasing is on.
func (img *Image) drawLineAA(gc GraphicsContext, ix0, iy0, ix1, iy1 int) {
	x0, y0 := float64(ix0), float64(iy0)
	x1, y1 := float64(ix1), float64(iy1)
	steep := math.Abs(y1-y0) > math.Abs(x1-x0)

	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	dx := x1 - x0
	dy := y1 - y0
	gradient := 1.0
	if dx != 0.0 {
		gradient = dy / dx
	}

	// first endpoint
	xend := math.Floor(x0 + 0.5)
	yend := y0 + gradient*(xend-x0)
	xgap := rfpart(x0 + 0.5)
	xpxl1 := int(xend)
	if steep {
		img.aaPlot(gc, int(math.Floor(yend)), xpxl1, rfpart(yend)*xgap)
		img.aaPlot(gc, int(math.Floor(yend))+1, xpxl1, fpart(yend)*xgap)
	} else {
		img.aaPlot(gc, xpxl1, int(math.Floor(yend)), rfpart(yend)*xgap)
		img.aaPlot(gc, xpxl1, int(math.Floor(yend))+1, fpart(yend)*xgap)
	}
	intery := yend + gradient

	// second endpoint
	xend = math.Floor(x1 + 0.5)
	yend = y1 + gradient*(xend-x1)
	xgap = fpart(x1 + 0.5)
	xpxl2 := int(xend)
	if steep {
		img.aaPlot(gc, int(math.Floor(yend)), xpxl2, rfpart(yend)*xgap)
		img.aaPlot(gc, int(math.Floor(yend))+1, xpxl2, fpart(yend)*xgap)
	} else {
		img.aaPlot(gc, xpxl2, int(math.Floor(yend)), rfpart(yend)*xgap)
		img.aaPlot(gc, xpxl2, int(math.Floor(yend))+1, fpart(yend)*xgap)
	}

	// main span
	for x := xpxl1 + 1; x < xpxl2; x++ {
		if steep {
			img.aaPlot(gc, int(math.Floor(intery)), x, rfpart(intery))
			img.aaPlot(gc, int(math.Floor(intery))+1, x, fpart(intery))
		} else {
			img.aaPlot(gc, x, int(math.Floor(intery)), rfpart(intery))
			img.aaPlot(gc, x, int(math.Floor(intery))+1, fpart(intery))
		}
		intery += gradient
	}
}

// circPlot8 blends the foreground at the 8 symmetric octant points of (dx,dy)
// about the circle center, each with fractional coverage cov.
func (img *Image) circPlot8(gc GraphicsContext, cx, cy, dx, dy int, cov float64) {
	c := int(cov*255.0 + 0.5)
	if c <= 0 {
		return
	}
	if c > 255 {
		c = 255
	}
	cu := uint32(c)
	img.blendPoint(gc, cx+dx, cy+dy, cu)
	img.blendPoint(gc, cx-dx, cy+dy, cu)
	img.blendPoint(gc, cx+dx, cy-dy, cu)
	img.blendPoint(gc, cx-dx, cy-dy, cu)
	// The (dy,dx) reflection coincides with the above on the 45-degree
	// diagonal; skip it there so those pixels are not blended twice.
	if dx != dy {
		img.blendPoint(gc, cx+dy, cy+dx, cu)
		img.blendPoint(gc, cx-dy, cy+dx, cu)
		img.blendPoint(gc, cx+dy, cy-dx, cu)
		img.blendPoint(gc, cx-dy, cy-dx, cu)
	}
}

// aaCircle draws Xiaolin Wu's anti-aliased circle outline (radius r about
// (cx,cy)).
func (img *Image) aaCircle(gc GraphicsContext, cx, cy, r int) {
	if r < 1 {
		img.blendPoint(gc, cx, cy, 255)
		return
	}
	rr := float64(r) * float64(r)
	for xx := 0; ; xx++ {
		yy := math.Sqrt(rr - float64(xx)*float64(xx))
		yi := int(math.Floor(yy))
		if xx > yi { // past 45 degrees; octant symmetry covers the rest
			break
		}
		frac := yy - float64(yi)
		img.circPlot8(gc, cx, cy, xx, yi, 1.0-frac)
		img.circPlot8(gc, cx, cy, xx, yi+1, frac)
	}
}

// aaEllipseBounds clips the (cx±rx, cy±ry) bounding box (expanded by pad) to the
// image and returns it.
func (img *Image) aaEllipseBounds(cx, cy, rx, ry, pad int) (x0, y0, x1, y1 int) {
	x0 = cx - rx - pad
	y0 = cy - ry - pad
	x1 = cx + rx + pad
	y1 = cy + ry + pad
	if x0 < 0 {
		x0 = 0
	}
	if y0 < 0 {
		y0 = 0
	}
	if x1 >= img.width {
		x1 = img.width - 1
	}
	if y1 >= img.height {
		y1 = img.height - 1
	}
	return
}

// ellipseOutlineCoverage returns the approximate 1px-wide outline coverage at
// (x,y) for the ellipse f = (dx/rx)^2 + (dy/ry)^2 - 1, using f divided by the
// gradient magnitude as the distance to the curve. Returns (coverage, ok).
func ellipseOutlineCoverage(x, y, cx, cy int, rx2, ry2 float64) (float64, bool) {
	dx := float64(x - cx)
	dy := float64(y - cy)
	f := dx*dx/rx2 + dy*dy/ry2 - 1.0
	gx := 2.0 * dx / rx2
	gy := 2.0 * dy / ry2
	grad := math.Sqrt(gx*gx + gy*gy)
	if grad < 1e-9 {
		return 0, false
	}
	cov := 1.0 - math.Abs(f)/grad
	if cov <= 0.0 {
		return 0, false
	}
	if cov > 1.0 {
		cov = 1.0
	}
	return cov, true
}

// aaEllipseOutline draws an anti-aliased ellipse outline (radii rx,ry about
// (cx,cy)) via the implicit-distance estimate.
func (img *Image) aaEllipseOutline(gc GraphicsContext, cx, cy, rx, ry int) {
	if rx < 1 || ry < 1 {
		return
	}
	rx2 := float64(rx) * float64(rx)
	ry2 := float64(ry) * float64(ry)
	x0, y0, x1, y1 := img.aaEllipseBounds(cx, cy, rx, ry, 1)
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if cov, ok := ellipseOutlineCoverage(x, y, cx, cy, rx2, ry2); ok {
				img.blendPoint(gc, x, y, uint32(cov*255.0+0.5))
			}
		}
	}
}

// arcAngleInRange reports whether the point (dx,dy) is within the arc's angular
// range. lo<=hi are transformed-degree bounds; the +/-360 tests handle the wrap
// at 0/360.
func arcAngleInRange(dx, dy, rx, ry, lo, hi float64) bool {
	deg := math.Atan2(dy/ry, dx/rx) * (180.0 / math.Pi)
	if deg < 0.0 {
		deg += 360.0
	}
	return (deg >= lo && deg <= hi) ||
		(deg+360.0 >= lo && deg+360.0 <= hi) ||
		(deg-360.0 >= lo && deg-360.0 <= hi)
}

// aaArcOutline draws an anti-aliased arc outline: the implicit-distance ellipse
// outline restricted to the angular range lo..hi (transformed degrees).
func (img *Image) aaArcOutline(gc GraphicsContext, cx, cy, rx, ry int, lo, hi float64) {
	if rx < 1 || ry < 1 {
		return
	}
	rx2 := float64(rx) * float64(rx)
	ry2 := float64(ry) * float64(ry)
	x0, y0, x1, y1 := img.aaEllipseBounds(cx, cy, rx, ry, 1)
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			cov, ok := ellipseOutlineCoverage(x, y, cx, cy, rx2, ry2)
			if !ok {
				continue
			}
			if !arcAngleInRange(float64(x-cx), float64(y-cy), float64(rx), float64(ry), lo, hi) {
				continue
			}
			img.blendPoint(gc, x, y, uint32(cov*255.0+0.5))
		}
	}
}

// fillEllipseAA fills an anti-aliased ellipse via NxN supersampled coverage
// against the true ellipse equation.
func (img *Image) fillEllipseAA(gc GraphicsContext, cx, cy, rx, ry int) {
	if rx < 1 || ry < 1 {
		return
	}
	rx2 := float64(rx) * float64(rx)
	ry2 := float64(ry) * float64(ry)
	x0, y0, x1, y1 := img.aaEllipseBounds(cx, cy, rx, ry, 0)
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			cnt := 0
			for sy := 0; sy < aaSupersample; sy++ {
				for sx := 0; sx < aaSupersample; sx++ {
					dx := (float64(x) + (float64(sx)+0.5)/aaSupersample) - float64(cx)
					dy := (float64(y) + (float64(sy)+0.5)/aaSupersample) - float64(cy)
					if dx*dx/rx2+dy*dy/ry2 <= 1.0 {
						cnt++
					}
				}
			}
			if cnt > 0 {
				img.blendPoint(gc, x, y, uint32(cnt*255/(aaSupersample*aaSupersample)))
			}
		}
	}
}

// fillArcAA fills an anti-aliased pie wedge: supersampled coverage of the
// ellipse sector between transformed-degree bounds lo..hi. A tiny disk at the
// center keeps the apex connected.
func (img *Image) fillArcAA(gc GraphicsContext, cx, cy, rx, ry int, lo, hi float64) {
	if rx < 1 || ry < 1 {
		return
	}
	rx2 := float64(rx) * float64(rx)
	ry2 := float64(ry) * float64(ry)
	x0, y0, x1, y1 := img.aaEllipseBounds(cx, cy, rx, ry, 0)
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			cnt := 0
			for sy := 0; sy < aaSupersample; sy++ {
				for sx := 0; sx < aaSupersample; sx++ {
					dx := (float64(x) + (float64(sx)+0.5)/aaSupersample) - float64(cx)
					dy := (float64(y) + (float64(sy)+0.5)/aaSupersample) - float64(cy)
					if dx*dx/rx2+dy*dy/ry2 <= 1.0 &&
						(dx*dx+dy*dy < 1.0 ||
							arcAngleInRange(dx, dy, float64(rx), float64(ry), lo, hi)) {
						cnt++
					}
				}
			}
			if cnt > 0 {
				img.blendPoint(gc, x, y, uint32(cnt*255/(aaSupersample*aaSupersample)))
			}
		}
	}
}

// pointInPolygon is an even-odd point-in-polygon test (works for convex or
// concave polygons).
func pointInPolygon(p []Point, px, py float64) bool {
	inside := false
	n := len(p)
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		if (float64(p[i].Y) > py) != (float64(p[j].Y) > py) &&
			px < float64(p[j].X-p[i].X)*(py-float64(p[i].Y))/float64(p[j].Y-p[i].Y)+float64(p[i].X) {
			inside = !inside
		}
	}
	return inside
}

// fillPolygonAA fills a polygon with anti-aliased edges via NxN supersampled
// coverage. Handles convex or concave polygons.
func (img *Image) fillPolygonAA(gc GraphicsContext, pts []Point) {
	minx, maxx := pts[0].X, pts[0].X
	miny, maxy := pts[0].Y, pts[0].Y
	for i := 1; i < len(pts); i++ {
		minx = min(minx, pts[i].X)
		maxx = max(maxx, pts[i].X)
		miny = min(miny, pts[i].Y)
		maxy = max(maxy, pts[i].Y)
	}
	if minx < 0 {
		minx = 0
	}
	if miny < 0 {
		miny = 0
	}
	if maxx >= img.width {
		maxx = img.width - 1
	}
	if maxy >= img.height {
		maxy = img.height - 1
	}

	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			cnt := 0
			for sy := 0; sy < aaSupersample; sy++ {
				for sx := 0; sx < aaSupersample; sx++ {
					if pointInPolygon(pts,
						float64(x)+(float64(sx)+0.5)/aaSupersample,
						float64(y)+(float64(sy)+0.5)/aaSupersample) {
						cnt++
					}
				}
			}
			if cnt > 0 {
				img.blendPoint(gc, x, y, uint32(cnt*255/(aaSupersample*aaSupersample)))
			}
		}
	}
}
