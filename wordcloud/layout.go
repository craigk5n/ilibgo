package main

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"math/rand"
	"sort"

	"github.com/craigk5n/ilibgo"
)

// ascentRatio approximates the fraction of a line's height that sits above the
// baseline. The library's Font.Ascent() reports the full line height for
// TrueType faces rather than the true baseline ascent, so we estimate it to seat
// glyphs within the scratch image before cropping. Typical Latin faces fall in
// the 0.75-0.85 range.
const ascentRatio = 0.8

// boxPadding is the gap (in pixels) kept clear around each word so neighbours
// never touch.
const boxPadding = 3

// occCell is the side length (in pixels) of one occupancy-grid cell. The grid is
// a downscaled view of the canvas: collision tests and stamping happen at cell
// granularity, which keeps the integral-image scan fast while staying fine
// enough that words pack tightly. Smaller = tighter packing but slower.
const occCell = 3

// inkThreshold is the alpha above which a rendered pixel counts as glyph "ink".
// Anti-aliased edges below this are treated as empty so faint halos do not block
// neighbours.
const inkThreshold = 32

// preferHorizontalDefault is the probability a word is laid out horizontally when
// the config does not specify one. The remainder are rotated 90°, which lets
// vertical words interlock into the gaps between horizontal ones.
const preferHorizontalDefault = 0.9

// placement is a word resolved to a canvas position and a pre-rendered sprite.
// The sprite carries the glyph pixels (colored, premultiplied alpha, transparent
// elsewhere) cropped to its inked bounds; (x, y) is its top-left on the canvas.
type placement struct {
	text   string
	x, y   int
	sprite *image.RGBA
}

// fontLoader returns a *Font at a given point size, caching by rounded size so
// repeated sizes are only rasterized once.
type fontLoader func(size float64) (*ilibgo.Font, error)

// layout places words into a width x height canvas using the approach taken by
// mature word-cloud tools (e.g. Python's wordcloud): each word is rasterized to
// a glyph mask, then dropped at a random free position rather than spiralled out
// from the center. Collision is tested against an integral-image occupancy map
// at glyph-pixel granularity, so small words nest into the whitespace of larger
// ones, and the random placement spreads big words across the whole canvas
// instead of stacking them in the middle. Words are placed largest-first (a big
// word that lands late has little room); words that find no free spot are
// returned in skipped.
func layout(words []sizedWord, cfg Config, load fontLoader, rng *rand.Rand) (placed []placement, skipped []string, err error) {
	// Largest first: big words need the open canvas, small ones fill the gaps.
	sorted := make([]sizedWord, len(words))
	copy(sorted, words)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].size > sorted[j].size })

	preferH := cfg.PreferHorizontal
	if preferH <= 0 {
		preferH = preferHorizontalDefault
	}
	padCells := (boxPadding + occCell - 1) / occCell
	marginCells := (cfg.Margin + occCell - 1) / occCell
	occ := newOccMap(cfg.Width, cfg.Height, occCell)

	for _, w := range sorted {
		font, ferr := load(w.size)
		if ferr != nil {
			return nil, nil, ferr
		}
		col, cerr := resolveColor(w.color, cfg.Palette, rng)
		if cerr != nil {
			return nil, nil, cerr
		}

		// Most words horizontal; a minority rotated so they interlock.
		rotate := rng.Float64() >= preferH
		sprite, serr := renderSprite(w.text, col, font, rotate)
		if serr != nil {
			return nil, nil, serr
		}
		sb := sprite.Bounds()
		sw, sh := sb.Dx(), sb.Dy()
		if sw > cfg.Width || sh > cfg.Height {
			skipped = append(skipped, w.text)
			continue
		}

		gw, gh, cells := spriteCells(sprite, occCell)
		gx, gy, ok := occ.findSpot(gw, gh, padCells, marginCells, rng)
		if !ok {
			skipped = append(skipped, w.text)
			continue
		}
		occ.stamp(gx, gy, cells)

		// Cell -> pixel, clamped so grid rounding never pushes the sprite off-canvas.
		px := clamp(gx*occCell, 0, cfg.Width-sw)
		py := clamp(gy*occCell, 0, cfg.Height-sh)
		placed = append(placed, placement{text: w.text, x: px, y: py, sprite: sprite})
	}
	return placed, skipped, nil
}

// renderSprite rasterizes one word into a tight RGBA sprite: glyphs drawn in col
// on a transparent background, cropped to their inked bounds, optionally rotated
// 90° (reading bottom-to-top). Pixels are premultiplied so the sprite composites
// correctly with image/draw's Over operator (ilibgo stores straight alpha).
func renderSprite(text string, col ilibgo.Color, font *ilibgo.Font, rotate bool) (*image.RGBA, error) {
	gc := ilibgo.CreateGraphicsContext()
	ilibgo.SetFont(&gc, font)
	ilibgo.SetForeground(&gc, col)
	ilibgo.SetAntiAlias(&gc, true)

	tw, th, err := ilibgo.TextDimensions(gc, font, text)
	if err != nil {
		return nil, fmt.Errorf("measure %q: %w", text, err)
	}
	const margin = 4
	scratchW, scratchH := tw+2*margin, th+2*margin
	scratch := ilibgo.CreateImage(scratchW, scratchH)
	baseline := margin + int(math.Round(ascentRatio*float64(th)))
	scratch.DrawString(gc, margin, baseline, text)

	// Copy straight-alpha pixels out, premultiply, and track the inked bounds.
	full := image.NewRGBA(image.Rect(0, 0, scratchW, scratchH))
	minX, minY, maxX, maxY := scratchW, scratchH, -1, -1
	for y := 0; y < scratchH; y++ {
		for x := 0; x < scratchW; x++ {
			r, g, b, a, _ := scratch.GetPixelAlpha(x, y)
			if a > inkThreshold {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
			i := full.PixOffset(x, y)
			full.Pix[i+0] = uint8(r * a / 255)
			full.Pix[i+1] = uint8(g * a / 255)
			full.Pix[i+2] = uint8(b * a / 255)
			full.Pix[i+3] = uint8(a)
		}
	}
	if maxX < minX || maxY < minY {
		// No ink (e.g. a lone space): a 1x1 transparent sprite occupies nothing.
		return image.NewRGBA(image.Rect(0, 0, 1, 1)), nil
	}

	cw, ch := maxX-minX+1, maxY-minY+1
	cropped := image.NewRGBA(image.Rect(0, 0, cw, ch))
	draw.Draw(cropped, cropped.Bounds(), full, image.Point{X: minX, Y: minY}, draw.Src)
	if rotate {
		cropped = rotate90(cropped)
	}
	return cropped, nil
}

// rotate90 rotates an RGBA image 90° counter-clockwise (so rotated words read
// bottom-to-top, the usual word-cloud convention).
func rotate90(src *image.RGBA) *image.RGBA {
	b := src.Bounds()
	sw, sh := b.Dx(), b.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, sh, sw))
	for y := 0; y < sh; y++ {
		for x := 0; x < sw; x++ {
			dst.Set(y, sw-1-x, src.RGBAAt(b.Min.X+x, b.Min.Y+y))
		}
	}
	return dst
}

// spriteCells returns the sprite's footprint in occupancy-grid cells: its width
// and height in cells, and the offsets of the cells that actually contain ink.
// The bounding box (gw x gh) is used to test for a free spot; only the inked
// cells are marked occupied, leaving a word's internal whitespace free for
// smaller words to nest into.
func spriteCells(sprite *image.RGBA, cell int) (gw, gh int, cells [][2]int) {
	b := sprite.Bounds()
	sw, sh := b.Dx(), b.Dy()
	gw = (sw + cell - 1) / cell
	gh = (sh + cell - 1) / cell
	seen := make([]bool, gw*gh)
	for y := 0; y < sh; y++ {
		for x := 0; x < sw; x++ {
			if sprite.RGBAAt(b.Min.X+x, b.Min.Y+y).A > inkThreshold {
				cx, cy := x/cell, y/cell
				if idx := cy*gw + cx; !seen[idx] {
					seen[idx] = true
					cells = append(cells, [2]int{cx, cy})
				}
			}
		}
	}
	return gw, gh, cells
}

// occMap is a downscaled occupancy grid with a summed-area table (integral image)
// for O(1) "is this rectangle free?" queries. A cell is occupied once any placed
// glyph pixel falls in it. The integral is rebuilt after each placement, which is
// cheap at grid resolution.
type occMap struct {
	cell     int
	w, h     int   // grid dimensions, in cells
	occ      []uint8
	integral []int // (w+1) * (h+1)
}

func newOccMap(pxW, pxH, cell int) *occMap {
	w := (pxW + cell - 1) / cell
	h := (pxH + cell - 1) / cell
	return &occMap{
		cell:     cell,
		w:        w,
		h:        h,
		occ:      make([]uint8, w*h),
		integral: make([]int, (w+1)*(h+1)),
	}
}

// regionSum returns the number of occupied cells in the grid rectangle
// [gx0,gx1) x [gy0,gy1), clamped to the grid.
func (m *occMap) regionSum(gx0, gy0, gx1, gy1 int) int {
	gx0 = clamp(gx0, 0, m.w)
	gy0 = clamp(gy0, 0, m.h)
	gx1 = clamp(gx1, 0, m.w)
	gy1 = clamp(gy1, 0, m.h)
	stride := m.w + 1
	return m.integral[gy1*stride+gx1] - m.integral[gy0*stride+gx1] -
		m.integral[gy1*stride+gx0] + m.integral[gy0*stride+gx0]
}

// findSpot scans every grid position where a gw x gh box (inflated by pad cells
// of breathing room) lands entirely on free cells, and returns one chosen
// uniformly at random via reservoir sampling — the random pick is what spreads
// words, including big ones, across the whole canvas. margin keeps the box that
// many cells away from the canvas edge, so words are never placed flush against
// the border.
func (m *occMap) findSpot(gw, gh, pad, margin int, rng *rand.Rand) (gx, gy int, ok bool) {
	minX, minY := margin, margin
	maxX, maxY := m.w-gw-margin, m.h-gh-margin
	if maxX < minX || maxY < minY {
		return 0, 0, false
	}
	count := 0
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if m.regionSum(x-pad, y-pad, x+gw+pad, y+gh+pad) == 0 {
				count++
				if rng.Intn(count) == 0 {
					gx, gy, ok = x, y, true
				}
			}
		}
	}
	return gx, gy, ok
}

// stamp marks the sprite's inked cells (placed at top-left gx,gy) as occupied and
// rebuilds the integral image.
func (m *occMap) stamp(gx, gy int, cells [][2]int) {
	for _, c := range cells {
		x, y := gx+c[0], gy+c[1]
		if x >= 0 && x < m.w && y >= 0 && y < m.h {
			m.occ[y*m.w+x] = 1
		}
	}
	stride := m.w + 1
	for y := 0; y < m.h; y++ {
		for x := 0; x < m.w; x++ {
			m.integral[(y+1)*stride+(x+1)] = int(m.occ[y*m.w+x]) +
				m.integral[y*stride+(x+1)] + m.integral[(y+1)*stride+x] -
				m.integral[y*stride+x]
		}
	}
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// resolveColor picks the word's color: an explicit per-word override (name or
// hex) when set, otherwise a random pick from the palette.
func resolveColor(override string, palette []string, rng *rand.Rand) (ilibgo.Color, error) {
	if override != "" {
		return parseColor(override)
	}
	return parseColor(palette[rng.Intn(len(palette))])
}

// render draws the background and composites every placed sprite onto a fresh
// image with source-over alpha blending.
func render(cfg Config, placed []placement) (*ilibgo.Image, error) {
	bg, err := parseColor(cfg.Background)
	if err != nil {
		return nil, fmt.Errorf("background color: %w", err)
	}
	img := ilibgo.CreateImageWithBackground(cfg.Width, cfg.Height, bg)
	for _, p := range placed {
		b := p.sprite.Bounds()
		r := image.Rect(p.x, p.y, p.x+b.Dx(), p.y+b.Dy())
		draw.Draw(img, r, p.sprite, b.Min, draw.Over)
	}
	return img, nil
}
