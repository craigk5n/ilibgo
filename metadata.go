package ilibgo

import "image"

// Image-level metadata and whole-image operations ported from the C library's
// IImage.c: DuplicateImage, SetComment/GetComment, SetTransparent/GetTransparent.

// DuplicateImage returns a deep copy of the image: a new image with identical
// dimensions and pixel data, plus the same comment and transparent color.
// Mirrors C IDuplicateImage.
func (img *Image) DuplicateImage() *Image {
	dup := &Image{
		width:    img.width,
		height:   img.height,
		data:     image.NewRGBA(image.Rect(0, 0, img.width, img.height)),
		comments: img.comments,
	}
	copy(dup.data.Pix, img.data.Pix)
	if img.transparent != nil {
		t := *img.transparent
		dup.transparent = &t
	}
	return dup
}

// SetComment sets the free-form text comment carried with the image. Mirrors C
// ISetComment.
func (img *Image) SetComment(comment string) {
	img.comments = comment
}

// GetComment returns the image's comment (empty if none was set). Mirrors C
// IGetComment.
func (img *Image) GetComment() string {
	return img.comments
}

// SetTransparent marks the given color as the image's transparent color, used
// by formats that support a single transparent color. Mirrors C
// ISetTransparent.
func (img *Image) SetTransparent(color Color) {
	c := color
	img.transparent = &c
}

// GetTransparent returns the image's transparent color and true if one has been
// set; otherwise it returns the zero Color and false. Mirrors C
// IGetTransparent (which returns INoTransparentColor when none is set).
func (img *Image) GetTransparent() (Color, bool) {
	if img.transparent == nil {
		return Color{}, false
	}
	return *img.transparent, true
}
