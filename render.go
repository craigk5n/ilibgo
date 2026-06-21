package ilibgo

// Low-level pixel compositing shared by the drawing primitives. Ported from the
// C library's _ISetPoint / _IDrawPoint / _IBlendPoint (IPixel.c, IlibP.h).
//
// The package stores straight (non-premultiplied) 8-bit RGBA in the backing
// image.RGBA's Pix slice, matching how NewColor and the rest of the library
// already treat color values. The blending math below is therefore expressed in
// straight alpha, identical to the C implementation.

// div255 rounds n/255 for non-negative integers (C IDIV255).
func div255(n uint32) uint32 { return (n + 127) / 255 }

// drawPoint writes the graphics context's foreground at (x, y). With BlendOver
// the foreground is composited (source-over) using its alpha; otherwise it
// overwrites the pixel. Out-of-bounds coordinates are ignored. This is the Go
// equivalent of the C _IDrawPoint / _ISetPoint macro.
func (img *Image) drawPoint(gc GraphicsContext, x int, y int) {
	if gc.blendMode == BlendOver {
		img.blendPoint(gc, x, y, 255)
		return
	}
	if x < 0 || x >= img.width || y < 0 || y >= img.height {
		return
	}
	i := img.data.PixOffset(x, y)
	pix := img.data.Pix
	pix[i] = gc.foreground.color.R
	pix[i+1] = gc.foreground.color.G
	pix[i+2] = gc.foreground.color.B
	pix[i+3] = gc.foreground.color.A
}

// blendPoint composites the graphics context's foreground onto the pixel at
// (x, y) with fractional edge coverage cover (0..255, where 255 is full
// coverage). The effective source alpha is the foreground alpha scaled by the
// coverage. Used by the anti-aliased rasterizers and by BlendOver. This is the
// Go equivalent of the C _IBlendPoint.
func (img *Image) blendPoint(gc GraphicsContext, x int, y int, cover uint32) {
	if x < 0 || x >= img.width || y < 0 || y >= img.height {
		return
	}
	fg := gc.foreground.color

	// Effective source alpha = color alpha scaled by edge coverage.
	sa := div255(uint32(fg.A) * cover)
	if sa == 0 {
		return
	}
	inv := 255 - sa
	sr := uint32(fg.R)
	sg := uint32(fg.G)
	sb := uint32(fg.B)

	i := img.data.PixOffset(x, y)
	pix := img.data.Pix

	// Straight-alpha source-over onto the RGBA destination.
	da := uint32(pix[i+3])
	dcontrib := div255(da * inv) // dst alpha contribution after (1-sa)
	// oa = sa + round(da*(255-sa)/255) <= sa + (255-sa) = 255, so the uint8
	// cast below never overflows.
	oa := sa + dcontrib
	if oa == 0 {
		pix[i] = 0
		pix[i+1] = 0
		pix[i+2] = 0
		pix[i+3] = 0
		return
	}
	pix[i] = uint8((sr*sa + uint32(pix[i])*dcontrib + oa/2) / oa)
	pix[i+1] = uint8((sg*sa + uint32(pix[i+1])*dcontrib + oa/2) / oa)
	pix[i+2] = uint8((sb*sa + uint32(pix[i+2])*dcontrib + oa/2) / oa)
	pix[i+3] = uint8(oa)
}
