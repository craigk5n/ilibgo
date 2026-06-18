package ilibgo

// This file collects backward-compatibility aliases. The drawing operations
// were converted from free functions taking an *Image to methods on *Image
// (e.g. FillRectangle(img, gc, ...) -> img.FillRectangle(gc, ...)), and the
// arc/circle functions also dropped a leading "I" left over from the original
// C API. The functions below preserve the old call sites by forwarding to the
// canonical methods; they may be removed in a future major release.

// --- Point ---

// Deprecated: use (*Image).DrawPoint.
func DrawPoint(img *Image, gc GraphicsContext, x int, y int) error {
	return img.DrawPoint(gc, x, y)
}

// Deprecated: use (*Image).SetPoint.
func SetPoint(img *Image, gc GraphicsContext, x int, y int) error {
	return img.SetPoint(gc, x, y)
}

// Deprecated: use (*Image).GetPoint.
func GetPoint(img *Image, x int, y int) Color {
	return img.GetPoint(x, y)
}

// --- Line / rectangle ---

// Deprecated: use (*Image).DrawLine.
func DrawLine(img *Image, gc GraphicsContext, x1 int, y1 int, x2 int, y2 int) error {
	return img.DrawLine(gc, x1, y1, x2, y2)
}

// Deprecated: use (*Image).DrawRectangle.
func DrawRectangle(img *Image, gc GraphicsContext, x int, y int, width int, height int) error {
	return img.DrawRectangle(gc, x, y, width, height)
}

// Deprecated: use (*Image).FillRectangle.
func FillRectangle(img *Image, gc GraphicsContext, x int, y int, width int, height int) error {
	return img.FillRectangle(gc, x, y, width, height)
}

// --- Arcs / circles ---

// Deprecated: use (*Image).DrawArc.
func DrawArc(img *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	return img.DrawArc(gc, x, y, r1, r2, a1, a2)
}

// Deprecated: use (*Image).DrawEnclosedArc.
func DrawEnclosedArc(img *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	return img.DrawEnclosedArc(gc, x, y, r1, r2, a1, a2)
}

// Deprecated: use (*Image).FillArc.
func FillArc(img *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	return img.FillArc(gc, x, y, r1, r2, a1, a2)
}

// Deprecated: use (*Image).DrawCircle.
func DrawCircle(img *Image, gc GraphicsContext, x int, y int, r int) error {
	return img.DrawCircle(gc, x, y, r)
}

// Deprecated: use (*Image).FillCircle.
func FillCircle(img *Image, gc GraphicsContext, x int, y int, r int) error {
	return img.FillCircle(gc, x, y, r)
}

// Deprecated: use (*Image).DrawEllipse.
func DrawEllipse(img *Image, gc GraphicsContext, x int, y int, r1 int, r2 int) error {
	return img.DrawEllipse(gc, x, y, r1, r2)
}

// --- Polygons / flood fill ---

// Deprecated: use (*Image).DrawPolygon.
func DrawPolygon(img *Image, gc GraphicsContext, points []Point) error {
	return img.DrawPolygon(gc, points)
}

// Deprecated: use (*Image).FillPolygon.
func FillPolygon(img *Image, gc GraphicsContext, points []Point) error {
	return img.FillPolygon(gc, points)
}

// Deprecated: use (*Image).FloodFill.
func FloodFill(img *Image, gc GraphicsContext, x int, y int) error {
	return img.FloodFill(gc, x, y)
}

// --- Copy ---

// Deprecated: use (*Image).CopyImage on the destination image.
func CopyImage(source *Image, dest *Image, gc GraphicsContext, src_x int, src_y int, width int, height int,
	dest_x int, dest_y int) error {
	return dest.CopyImage(source, gc, src_x, src_y, width, height, dest_x, dest_y)
}

// Deprecated: use (*Image).CopyImageScaled on the destination image.
func CopyImageScaled(source *Image, dest *Image,
	src_x int, src_y int, src_width int, src_height int,
	dest_x int, dest_y int, dest_width int, dest_height int) error {
	return dest.CopyImageScaled(source, src_x, src_y, src_width, src_height, dest_x, dest_y, dest_width, dest_height)
}

// --- Text ---

// Deprecated: use (*Image).DrawString.
func DrawString(img *Image, gc GraphicsContext, x int, y int, text string) {
	img.DrawString(gc, x, y, text)
}

// Deprecated: use (*Image).DrawStringRotated.
func DrawStringRotated(img *Image, gc GraphicsContext, x int, y int, text string, direction TextDirection) {
	img.DrawStringRotated(gc, x, y, text, direction)
}

// Deprecated: use (*Image).DrawStringRotatedAngle.
func DrawStringRotatedAngle(img *Image, gc GraphicsContext, x int, y int, text string, angle float64) {
	img.DrawStringRotatedAngle(gc, x, y, text, angle)
}

// --- Legacy "I"-prefixed arc/circle names (forward to the renamed methods) ---

// Deprecated: use (*Image).DrawArc.
func IDrawArc(img *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	return img.DrawArc(gc, x, y, r1, r2, a1, a2)
}

// Deprecated: use (*Image).DrawEnclosedArc.
func IDrawEnclosedArc(img *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	return img.DrawEnclosedArc(gc, x, y, r1, r2, a1, a2)
}

// Deprecated: use (*Image).FillArc.
func IFillArc(img *Image, gc GraphicsContext, x int, y int, r1 int, r2 int, a1 float64, a2 float64) error {
	return img.FillArc(gc, x, y, r1, r2, a1, a2)
}

// Deprecated: use (*Image).DrawCircle.
func IDrawCircle(img *Image, gc GraphicsContext, x int, y int, r int) error {
	return img.DrawCircle(gc, x, y, r)
}

// Deprecated: use (*Image).FillCircle.
func IFillCircle(img *Image, gc GraphicsContext, x int, y int, r int) error {
	return img.FillCircle(gc, x, y, r)
}

// Deprecated: use (*Image).DrawEllipse.
func IDrawEllipse(img *Image, gc GraphicsContext, x int, y int, r1 int, r2 int) error {
	return img.DrawEllipse(gc, x, y, r1, r2)
}
