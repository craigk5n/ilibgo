package ilibgo

type ImageFormat uint8

const (
	DefaultFormatString = "ppm"
	IlibVersion         = "2.0"
	IlibVersionDate     = "11 Aug 2022"

	FormatGIF       ImageFormat = 0
	FormatPPM       ImageFormat = 1
	FormatPGM       ImageFormat = 2
	FormatPBM       ImageFormat = 3
	FormatXPM       ImageFormat = 4
	FormatXBM       ImageFormat = 5
	FormatPNG       ImageFormat = 6
	FormatJPEG      ImageFormat = 7
	FormatTIFF      ImageFormat = 8
	FormatBMP       ImageFormat = 9
	NumberOfFormats uint        = 10
)

// Line styles for use in LineStyle
type LineStyle int

const (
	LineSolid      LineStyle = 1 // default
	LineOnOffDash  LineStyle = 2 // Not yet implemented
	LineDoubleDash LineStyle = 3 // Not yet implemented
)

// Fill styles for use with FillStyle
type FillStyle int

const (
	FillSolid          FillStyle = 1 // default
	FillTiled          FillStyle = 2 // Not yet implemented
	FillStippled       FillStyle = 3 // Not yet implemented
	FillOpaqueStippled FillStyle = 4 // Not yet implemented
)

type TextStyle int

const (
	TextNormal    TextStyle = 1 // default
	TextEtchedIn  TextStyle = 2 // text appears etched into background
	TextEtchedOut TextStyle = 3 // text appears etched out of background
	TextShadowed  TextStyle = 4 // text has drop shadow that fades into background
)

type TextDirection int

const (
	TextLeftToRight TextDirection = 1 // default
	TextBottomToTop TextDirection = 2
	TextTopToBottom TextDirection = 3
)
