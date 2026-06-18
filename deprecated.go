package ilibgo

// This file collects backward-compatibility aliases for functions that were
// renamed to drop the leading "I" (a holdover from the original C API). The
// aliases simply forward to the canonical names and may be removed in a future
// major release.

// Deprecated: use DrawArc instead.
func IDrawArc(image *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	return DrawArc(image, gc, x, y, r1, r2, a1, a2)
}

// Deprecated: use DrawEnclosedArc instead.
func IDrawEnclosedArc(image *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	return DrawEnclosedArc(image, gc, x, y, r1, r2, a1, a2)
}

// Deprecated: use FillArc instead.
func IFillArc(image *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	return FillArc(image, gc, x, y, r1, r2, a1, a2)
}

// Deprecated: use DrawCircle instead.
func IDrawCircle(image *Image, gc GraphicsContext, x int, y int, r int) error {
	return DrawCircle(image, gc, x, y, r)
}

// Deprecated: use FillCircle instead.
func IFillCircle(image *Image, gc GraphicsContext, x int, y int, r int) error {
	return FillCircle(image, gc, x, y, r)
}

// Deprecated: use DrawEllipse instead.
func IDrawEllipse(image *Image, gc GraphicsContext, x int, y int, r1 int, r2 int) error {
	return DrawEllipse(image, gc, x, y, r1, r2)
}
