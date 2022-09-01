package ilibgo

import (
	"errors"
	"math"
)

// History:
//      15-Aug-2022     Craig Knudsen craig@k5n.us
//                      Converted from C to go
//      21-Jan-2000     Geovan Rodriguez <geovan@cigb.edu.cu>
//                      Added IDrawStringRotatedAngle()
//      23-Aug-1999     Craig Knudsen   cknudsen@cknudsen.com
//                      Added support for text styles:
//                      ITEXT_NORMAL, ITEXT_ETCHED_IN,
//                      ITEXT_ETCHED_OUT, ITEXT_SHADOWED
//      21-Aug-1999     Craig Knudsen   cknudsen@cknudsen.com
//                      Added IDrawStringRotated()
//                      Removed anti-aliasing stuff since it didn't really
//                      work.  Eventually, true type fonts will be used
//                      for this (FreeType lib).
//      18-May-1998     Craig Knudsen   cknudsen@cknudsen.com
//                      Added support for anti-aliasing fonts by using a
//                      font 2X bigger than we need.
//      20-May-1996     Craig Knudsen   cknudsen@cknudsen.com
//                      Created

// Get the width (in pixels) of the specified text using
// the font currently set in the graphics context.
func TextWidth(gc GraphicsContext, font *Font, text string) (width int, err error) {
	width, _, err = TextDimensions(gc, font, text)
	return width, err
}

// Get the height (in pixels) of the specified text using
// the font currently set in the graphics context.
func TextHeight(gc GraphicsContext, font *Font, text string) (height int, err error) {
	_, height, err = TextDimensions(gc, font, text)
	return height, err
}

func TextDimensions(gc GraphicsContext, font *Font, text string) (width int, height int, err error) {
	charNum := 0
	retWidth := 0
	retHeight := 0
	charx := 0
	chary := 0
	//var bitData []rune

	if font == nil {
		return 0, 0, errors.New("no font set")
	}
	fontHeight, _ := GetFontSize(font)
	retHeight = fontHeight
	for _, char := range text {
		if char == '\012' {
			// new line
			charx = 0
			chary += fontHeight
			retHeight += fontHeight
			charNum = 0
			continue
		} else if char == '\t' {
			// tab
			character := FontBDFGetRune(font.font, char)
			charx += (8 - (charNum % 8)) * character.actualWidth
			if charx > retWidth {
				retWidth = charx
			}
			continue
		} else if char == '\033' {
			// Escape
			// TODO: handle escape sequences
		} else {
			character := FontBDFGetRune(font.font, char)
			charx += character.actualWidth
			charNum++
			if charx > retWidth {
				retWidth = charx
			}
		}
	}
	return retWidth, retHeight, nil
}

func DrawString(image *Image, gc GraphicsContext, x int, y int, text string) {
	DrawStringRotated(image, gc, x, y, text, TextLeftToRight)
}

// calculate values for topshadow and bottomshadow
func makeTopAndBottomShadow(color Color) (topColor Color, bottomColor Color) {
	var topR, topG, topB uint8
	if color.color.R > 205 {
		topR = 255
	} else {
		topR = color.color.R + 50
	}
	if color.color.G > 205 {
		topG = 255
	} else {
		topG = color.color.G + 50
	}
	if color.color.B > 205 {
		topB = 255
	} else {
		topB = color.color.B + 50
	}
	topColor, _ = AllocColor(topR, topG, topB)
	var bottomR, bottomG, bottomB uint8
	if color.color.R < 50 {
		bottomR = 0
	} else {
		bottomR = color.color.R - 50
	}
	if color.color.G < 50 {
		bottomG = 0
	} else {
		bottomG = color.color.G - 50
	}
	if color.color.B < 50 {
		bottomB = 0
	} else {
		bottomB = color.color.B - 50
	}
	bottomColor, _ = AllocColor(bottomR, bottomG, bottomB)
	return topColor, bottomColor
}

// calculate values for shadows
func makeShadows(incolor Color, numShadows int) []Color {
	var rinc, ginc, binc, rstart, gstart, bstart, loop uint8

	rstart = incolor.color.R / 2
	gstart = incolor.color.G / 2
	bstart = incolor.color.B / 2

	rinc = (incolor.color.R - rstart) / uint8(numShadows)
	ginc = (incolor.color.G - gstart) / uint8(numShadows)
	binc = (incolor.color.B - bstart) / uint8(numShadows)

	ret := make([]Color, numShadows)
	for loop = 0; loop < uint8(numShadows); loop++ {
		ret[loop], _ = AllocColor(rstart+rinc*loop, gstart+ginc*loop, bstart+binc*loop)
	}
	return ret
}

func DrawStringRotated(image *Image, gc GraphicsContext, x int, y int, text string, direction TextDirection) {
	var top, bottom Color
	var fontHeight int

	origFg := gc.foreground
	switch gc.textStyle {
	case TextNormal:
		drawStringRotated90(image, gc, x, y, text, direction)
	case TextEtchedIn, TextEtchedOut:
		if gc.textStyle == TextEtchedOut {
			top, bottom = makeTopAndBottomShadow(gc.background)
		} else {
			top, bottom = makeTopAndBottomShadow(gc.background)
		}
		SetForeground(&gc, top)
		drawStringRotated90(image, gc, x-1, y-1, text, direction)
		SetForeground(&gc, bottom)
		drawStringRotated90(image, gc, x+1, y+1, text, direction)
		gc.foreground = gc.background
		drawStringRotated90(image, gc, x, y, text, direction)
	case TextShadowed:
		fontHeight, _ = GetFontSize(gc.font)
		nshadows := fontHeight / 5
		shadows := makeShadows(gc.background, nshadows)
		gc.foreground = origFg
		for loop := nshadows; loop > 0; loop-- {
			SetForeground(&gc, shadows[loop-1])
			drawStringRotated90(image, gc, x+loop, y+loop, text, direction)
		}
		gc.foreground = origFg
		drawStringRotated90(image, gc, x, y, text, direction)
	}
}

func drawStringRotated90(image *Image, gc GraphicsContext, x int, y int, text string, direction TextDirection) {
	charNum := 0
	charx := x
	chary := y
	fontHeight, _ := GetFontSize(gc.font)
	var bdfChar *BdfChar

	chars := []rune(text)
	for loop := 0; loop < len(chars); loop++ {
		char := string(chars[loop])
		if char == "\n" {
			switch direction {
			case TextLeftToRight:
				charx = x
				chary += fontHeight
			case TextTopToBottom:
				chary = y
				charx -= fontHeight
			case TextBottomToTop:
				chary = y
				charx += fontHeight
			}
			charNum = 0
			continue
		} else if char == "\t" {
			bdfChar = FontBDFGetChar(gc.font.font, char)
			switch direction {
			case TextLeftToRight:
				charx += (8 - (charNum % 8)) * bdfChar.actualWidth
			case TextTopToBottom:
				chary -= (8 - (charNum % 8)) * bdfChar.actualWidth
			case TextBottomToTop:
				chary += (8 - (charNum % 8)) * bdfChar.actualWidth
			}
			continue
			//} else if char == '\033' {
			//	// Handle escape
			//	loop2 := 0
			//	for loop2 := 1; text[loop+loop2:loop+loop2] != ";" && (loop+loop2) < len(string) && loop2 < 32; {
			//		char = char + text[loop+loop2:loop+loop2]
			//	}
			//	loop2++
			//	loop += loop2
		}
		bdfChar = FontBDFGetChar(gc.font.font, char)
		if bdfChar != nil {
			for loop3 := 0; loop3 < bdfChar.height; loop3++ {
				switch direction {
				case TextLeftToRight:
					myy := chary - (bdfChar.height + bdfChar.yoffset) + loop3
					for loop4 := 0; loop4 < bdfChar.width; loop4++ {
						if bdfChar.data[loop3*bdfChar.width+loop4] {
							myx := charx + bdfChar.xoffset + loop4
							SetPoint(image, gc, myx, myy)
						}
					}
				case TextTopToBottom:
					myx := charx + (bdfChar.height + bdfChar.yoffset) - loop3
					for loop4 := 0; loop4 < bdfChar.width; loop4++ {
						if bdfChar.data[loop3*bdfChar.width+loop4] {
							myy := chary + bdfChar.xoffset + loop4
							SetPoint(image, gc, myx, myy)
						}
					}
				case TextBottomToTop:
					myx := charx - (bdfChar.height + bdfChar.yoffset) + loop3
					for loop4 := 0; loop4 < bdfChar.width; loop4++ {
						if bdfChar.data[loop3*bdfChar.width+loop4] {
							myy := chary - bdfChar.xoffset - loop4
							SetPoint(image, gc, myx, myy)
						}
					}
				}
			}
			switch direction {
			case TextLeftToRight:
				charx += bdfChar.actualWidth
			case TextTopToBottom:
				chary += bdfChar.actualWidth
			case TextBottomToTop:
				chary -= bdfChar.actualWidth
			}
			charNum++
		}
	}
}

func DrawStringRotatedAngle(image *Image, gc GraphicsContext, x int, y int, text string, angle float64) {
	var bdfChar *BdfChar
	var myx, myy int
	charNum := 0
	var x1, y1, x2, y2 float64
	charx := x
	chary := y
	fontHeight, _ := GetFontSize(gc.font)

	for loop := 0; loop < len(text); loop++ {
		char := text[loop:loop]
		if char == "\n" {
			chary += fontHeight
			charx = x
			charNum = 0
		} else if char == "\t" {
			bdfChar = FontBDFGetChar(gc.font.font, char)
			chary += (8 - (charNum % 8)) * bdfChar.actualWidth
			continue
			//} else if *ptr == '\033' {
			//	// Handle escape
			//	// Handle escape
			//	loop2 := 0
			//	for loop2 := 1; text[loop+loop2:loop+loop2] != ";" && (loop+loop2) < len(string) && loop2 < 32; {
			//		char = char + text[loop+loop2:loop+loop2]
			//	}
			//	loop2++
			//	loop += loop2
		}
		bdfChar = FontBDFGetChar(gc.font.font, char)
		if bdfChar != nil {
			for loop3 := 0; loop3 < bdfChar.height; loop3++ {
				for loop4 := 0; loop4 < bdfChar.width; loop4++ {
					if bdfChar.data[loop3*bdfChar.width+loop4] {
						x1 = float64(loop4 + bdfChar.xoffset)
						y1 = float64(loop3 + gc.font.font.pixelSize - bdfChar.height)
						x2 = x1*math.Cos(angle*math.Pi/180.0) +
							y1*math.Sin(angle*math.Pi/180.0)
						y2 = -1*x1*math.Sin(angle*math.Pi/180.0) +
							y1*math.Cos(angle*math.Pi/180.0)
						myy = chary - (gc.font.font.pixelSize + bdfChar.yoffset) + int(y2)
						myx = charx + int(x2)
						SetPoint(image, gc, myx, myy)
					}
				}
			}
			x1 = float64(bdfChar.actualWidth + bdfChar.xoffset)
			y1 = 0
			x2 = x1*math.Cos(angle*math.Pi/180.0) +
				y1*math.Sin(angle*math.Pi/180.0)
			y2 = -1*x1*math.Sin(angle*math.Pi/180.0) +
				y1*math.Cos(angle*math.Pi/180.0)
			charx += int(x2)
			chary += int(y2)
			charNum++
		}
	}
}
