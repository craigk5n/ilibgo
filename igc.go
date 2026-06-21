package ilibgo

import "errors"

// Create a graphics context
func CreateGraphicsContext() GraphicsContext {
	var gc GraphicsContext

	gc.foreground = NewColor(0, 0, 0, 255)       // black
	gc.background = NewColor(255, 255, 255, 255) // white
	gc.lineWidth = 1
	gc.lineStyle = LineSolid
	gc.textStyle = TextNormal
	gc.blendMode = BlendReplace

	return gc
}

func SetLineWidth(gc *GraphicsContext, lineWidth int) {
	// TODO: Support lineWidth greater than 3
	if lineWidth > 3 {
		lineWidth = 3
	}
	gc.lineWidth = lineWidth
}

func SetBackground(gc *GraphicsContext, backgroundColor Color) {
	gc.background = backgroundColor
}

func SetForeground(gc *GraphicsContext, foregroundColor Color) {
	gc.foreground = foregroundColor
}

func SetLineStyle(gc *GraphicsContext, lineStyle LineStyle) {
	gc.lineStyle = lineStyle
}

func SetTextStyle(gc *GraphicsContext, textStyle TextStyle) {
	gc.textStyle = textStyle
}

func SetFont(gc *GraphicsContext, font *Font) {
	gc.font = font
	gc.antiAliasedFont = false
}

// SetBlendMode sets the pixel compositing mode of a graphics context.
// BlendReplace (the default) overwrites pixels; BlendOver composites the
// foreground over the destination using the foreground's alpha.
func SetBlendMode(gc *GraphicsContext, mode BlendMode) {
	gc.blendMode = mode
}

// SetAntiAlias enables or disables anti-aliased rendering of drawing
// primitives (thin solid lines, arcs, circles, ellipses, and fills).
func SetAntiAlias(gc *GraphicsContext, on bool) {
	gc.antiAlias = on
}

// SetAntiAliasedFont sets the graphics context's font and requests
// anti-aliased rendering for it. For bitmap (BDF) fonts this is an
// experimental smoothing of the glyph bitmap; TrueType fonts always render
// anti-aliased. Equivalent to SetFont followed by enabling the flag.
func SetAntiAliasedFont(gc *GraphicsContext, font *Font) {
	gc.font = font
	gc.antiAliasedFont = true
}

func GetFontSize(font *Font) (height int, err error) {
	if font == nil {
		return 0, errors.New("font is nil")
	}
	if font.isTrueType() {
		return font.height, nil
	}
	if font.font == nil {
		return 0, errors.New("font is nil")
	}
	return font.font.pixelSize, nil
}
