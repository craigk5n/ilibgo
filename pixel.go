package ilibgo

import "fmt"

// Explicit per-pixel color access, ported from the C library's IPixel.c. Unlike
// SetPoint, these functions take color components directly (not from a graphics
// context) and always overwrite the pixel (no blending). Color components are
// in the range [0, 255]; coordinates must be within the image bounds.

func (img *Image) inBounds(x, y int) bool {
	return x >= 0 && x < img.width && y >= 0 && y < img.height
}

func validComponents(vals ...int) error {
	for _, v := range vals {
		if v < 0 || v > 255 {
			return fmt.Errorf("ilibgo: color component %d out of range [0,255]", v)
		}
	}
	return nil
}

// SetPixelAlpha sets the pixel at (x, y) to the given red/green/blue/alpha
// values, overwriting any existing pixel. Mirrors C ISetPixelAlpha.
func (img *Image) SetPixelAlpha(x, y, red, green, blue, alpha int) error {
	if err := validComponents(red, green, blue, alpha); err != nil {
		return err
	}
	if !img.inBounds(x, y) {
		return fmt.Errorf("ilibgo: SetPixelAlpha: (%d,%d) out of bounds", x, y)
	}
	i := img.data.PixOffset(x, y)
	pix := img.data.Pix
	pix[i] = uint8(red)
	pix[i+1] = uint8(green)
	pix[i+2] = uint8(blue)
	pix[i+3] = uint8(alpha)
	return nil
}

// SetPixel sets the pixel at (x, y) to the given red/green/blue values (fully
// opaque). Mirrors C ISetPixel.
func (img *Image) SetPixel(x, y, red, green, blue int) error {
	return img.SetPixelAlpha(x, y, red, green, blue, 255)
}

// GetPixelAlpha returns the red/green/blue/alpha values of the pixel at
// (x, y). Mirrors C IGetPixelAlpha.
func (img *Image) GetPixelAlpha(x, y int) (red, green, blue, alpha int, err error) {
	if !img.inBounds(x, y) {
		return 0, 0, 0, 0, fmt.Errorf("ilibgo: GetPixelAlpha: (%d,%d) out of bounds", x, y)
	}
	i := img.data.PixOffset(x, y)
	pix := img.data.Pix
	return int(pix[i]), int(pix[i+1]), int(pix[i+2]), int(pix[i+3]), nil
}

// GetPixel returns the red/green/blue values of the pixel at (x, y),
// discarding alpha. Mirrors C IGetPixel.
func (img *Image) GetPixel(x, y int) (red, green, blue int, err error) {
	r, g, b, _, e := img.GetPixelAlpha(x, y)
	return r, g, b, e
}
