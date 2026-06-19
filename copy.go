package ilibgo

import (
	"image"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

// Description:
//	Copy an area of an image to another image.
//
// History:
//  09-Aug-2022	Craig Knudsen craig@k5n.us
//      	Converted from C to go
//	15-Aug-2001	Craig Knudsen	cknudsen@cknudsen.com
//			Fixed bug in ICopyImageScaled
//			(thanks Gal Steinitz for this fix)
//	23-Jul-1999	Craig Knudsen   cknudsen@cknudsen.com
//			Added ICopyImageScaled
//	11-Nov-1998	Craig Knudsen	cknudsen@cknudsen.com
//			Allow transparent values to not be copied.
//	20-May-1996	Craig Knudsen	cknudsen@cknudsen.com
//			Created

// CopyImage copies a rectangular region of source onto this (destination)
// image. The gc parameter is retained for API compatibility but is unused; the
// copy is a straight pixel transfer clipped to both images' bounds.
func (dest *Image) CopyImage(source *Image, gc GraphicsContext, srcX int, srcY int, width int, height int,
	destX int, destY int) error {
	dstRect := image.Rect(destX, destY, destX+width, destY+height)
	srcPoint := image.Point{X: srcX, Y: srcY}
	draw.Draw(dest.data, dstRect, source.data, srcPoint, draw.Src)
	return nil
}

// ScaleQuality selects the resampling filter used by CopyImageScaledQuality.
// Higher-quality filters are smoother but slower; nearest-neighbor is fastest
// and blockiest.
type ScaleQuality int

const (
	// ScaleNearestNeighbor picks the closest source pixel. Fastest and
	// blockiest; matches the behavior of CopyImageScaled.
	ScaleNearestNeighbor ScaleQuality = iota
	// ScaleApproxBiLinear is a fast bilinear approximation — a good default
	// for thumbnails.
	ScaleApproxBiLinear
	// ScaleBiLinear is full bilinear interpolation (smoother than the
	// approximate variant).
	ScaleBiLinear
	// ScaleCatmullRom is a high-quality bicubic filter; best for downscaling
	// photographs where sharpness matters.
	ScaleCatmullRom
)

// CopyImageScaledQuality scales a region of source onto this (destination)
// image using the chosen resampling filter, delegating to
// golang.org/x/image/draw. Unlike CopyImageScaled (nearest-neighbor only), it
// can produce smooth, anti-aliased results. The destination region is
// overwritten (no alpha blending).
func (dest *Image) CopyImageScaledQuality(source *Image,
	srcX int, srcY int, srcWidth int, srcHeight int,
	destX int, destY int, destWidth int, destHeight int,
	quality ScaleQuality) error {

	var interp xdraw.Interpolator
	switch quality {
	case ScaleApproxBiLinear:
		interp = xdraw.ApproxBiLinear
	case ScaleBiLinear:
		interp = xdraw.BiLinear
	case ScaleCatmullRom:
		interp = xdraw.CatmullRom
	default:
		interp = xdraw.NearestNeighbor
	}
	dstRect := image.Rect(destX, destY, destX+destWidth, destY+destHeight)
	srcRect := image.Rect(srcX, srcY, srcX+srcWidth, srcY+srcHeight)
	interp.Scale(dest.data, dstRect, source.data, srcRect, xdraw.Src, nil)
	return nil
}

// CopyImageScaled scales a region of source onto this (destination) image
// using nearest-neighbor sampling. For smoother results, use
// CopyImageScaledQuality.
func (dest *Image) CopyImageScaled(source *Image,
	srcX int, srcY int, srcWidth int, srcHeight int,
	destX int, destY int, destWidth int, destHeight int) error {

	// When scaling down, we might want to add an algorithm for averaging
	// a series of source pixels into the destination pixel.  For now,
	// we just grab one color.
	var gc GraphicsContext
	scaleX := float64(destWidth) / float64(srcWidth)
	scaleY := float64(destHeight) / float64(srcHeight)
	for y := destY; y < destY+destHeight; y++ {
		for x := destX; x < destX+destWidth; x++ {
			// get location from source image for this location
			// x2,y2 is location in source image.
			tempX := float64(srcX) + float64(x-destX)/scaleX
			x2 := int(tempX)
			tempY := float64(srcY) + float64(y-destY)/scaleY
			y2 := int(tempY)
			gc.foreground = source.GetPoint(x2, y2)
			dest.SetPoint(gc, x, y)
		}
	}
	return nil
}
