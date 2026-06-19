package ilibgo

// Font metadata accessors. A Font's fields are unexported, so these provide
// read-only access to the descriptive information parsed from a BDF font (or,
// for a scalable font, the basics that apply). They are safe to call on a nil
// receiver and return zero values for fields that don't apply.

// Name returns the name the font was loaded with.
func (f *Font) Name() string {
	if f == nil {
		return ""
	}
	return f.name
}

// IsTrueType reports whether the font is backed by a scalable TrueType/OpenType
// face rather than a bitmap BDF font.
func (f *Font) IsTrueType() bool {
	return f.isTrueType()
}

// Foundry returns the BDF FOUNDRY (e.g. "Adobe"). Empty for TrueType fonts.
func (f *Font) Foundry() string {
	if f == nil || f.font == nil {
		return ""
	}
	return f.font.foundry
}

// Family returns the BDF FAMILY_NAME (e.g. "Helvetica"). Empty for TrueType.
func (f *Font) Family() string {
	if f == nil || f.font == nil {
		return ""
	}
	return f.font.family
}

// FaceName returns the BDF FACE_NAME. Empty for TrueType fonts.
func (f *Font) FaceName() string {
	if f == nil || f.font == nil {
		return ""
	}
	return f.font.faceName
}

// Slant returns the BDF SLANT code ("R" roman, "I" italic, "O" oblique).
// Empty for TrueType fonts.
func (f *Font) Slant() string {
	if f == nil || f.font == nil {
		return ""
	}
	return f.font.slant
}

// Weight returns the BDF WEIGHT_NAME (e.g. "Medium", "Bold"). Empty for
// TrueType fonts.
func (f *Font) Weight() string {
	if f == nil || f.font == nil {
		return ""
	}
	return f.font.weight
}

// Proportional reports whether the BDF font uses proportional (vs fixed)
// spacing. False for TrueType fonts.
func (f *Font) Proportional() bool {
	if f == nil || f.font == nil {
		return false
	}
	return f.font.proportional
}

// PixelSize returns the BDF PIXEL_SIZE. Zero for TrueType fonts.
func (f *Font) PixelSize() int {
	if f == nil || f.font == nil {
		return 0
	}
	return f.font.pixelSize
}

// Ascent returns the font ascent in pixels. For TrueType fonts this is derived
// from the rasterized face height.
func (f *Font) Ascent() int {
	if f == nil {
		return 0
	}
	if f.font != nil {
		return f.font.fontAscent
	}
	return f.height
}

// Descent returns the BDF font descent in pixels. Zero for TrueType fonts.
func (f *Font) Descent() int {
	if f == nil || f.font == nil {
		return 0
	}
	return f.font.fontDescent
}

// GlyphCount returns the number of defined glyphs (ASCII/Latin-1 plus any
// extra named glyphs). Zero for TrueType fonts.
func (f *Font) GlyphCount() int {
	if f == nil || f.font == nil {
		return 0
	}
	n := 0
	for i := range f.font.chars {
		if len(f.font.chars[i].data) > 0 {
			n++
		}
	}
	return n + len(f.font.otherChars)
}
