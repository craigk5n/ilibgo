package ilibgo

import "errors"

// Create a graphics context
func CreateGraphicsContext() GraphicsContext {
	var gc GraphicsContext

	gc.foreground = newIColor(0, 0, 0, 255)       // black
	gc.background = newIColor(255, 255, 255, 255) // white
	gc.lineWidth = 1
	gc.lineStyle = LineSolid
	gc.textStyle = TextNormal

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
}

func GetFontSize(font *Font) (height int, err error) {
	if font == nil || font.font == nil {
		return 0, errors.New("font is nil")
	}
	// Only support BDF fonts for now.
	return font.font.pixelSize, nil
}
