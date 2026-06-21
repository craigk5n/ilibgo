package ilibgo

import "sort"

// Color reduction (quantization) via median cut, ported from IQuantize.c. Used
// to fit an image into a limited palette, most notably for 8-bit GIF output.
//
// The histogram works at 5 bits/channel (32768 cells). That precision only
// matters once we are *already* reducing a >maxColors image (a lossy
// operation); images with few enough colors are passed through exactly. Alpha
// is preserved unchanged; only the RGB channels are quantized.

const (
	qBits   = 5
	qShift  = 8 - qBits
	qLevels = 1 << qBits       // 32
	qCells  = 1 << (qBits * 3) // 32768
)

func qPack(r, g, b uint8) int {
	return (int(r>>qShift) << (qBits * 2)) | (int(g>>qShift) << qBits) | int(b>>qShift)
}

type qBucket struct {
	r5, g5, b5 uint8   // 5-bit channel values
	count      int     // pixels in this cell
	sr, sg, sb float64 // sums of original 8-bit values (for averaging)
}

type qBox struct {
	lo, hi int // [lo, hi) slice of the bucket array
}

// qBoxExtent returns the longest axis (0=r,1=g,2=b) of the buckets in box and
// its extent.
func qBoxExtent(pop []qBucket, box qBox) (axis, extent int) {
	rmin, rmax := qLevels, -1
	gmin, gmax := qLevels, -1
	bmin, bmax := qLevels, -1
	for i := box.lo; i < box.hi; i++ {
		r, g, b := int(pop[i].r5), int(pop[i].g5), int(pop[i].b5)
		if r < rmin {
			rmin = r
		}
		if r > rmax {
			rmax = r
		}
		if g < gmin {
			gmin = g
		}
		if g > gmax {
			gmax = g
		}
		if b < bmin {
			bmin = b
		}
		if b > bmax {
			bmax = b
		}
	}
	rext, gext, bext := rmax-rmin, gmax-gmin, bmax-bmin
	if rext >= gext && rext >= bext {
		return 0, rext
	}
	if gext >= bext {
		return 1, gext
	}
	return 2, bext
}

func qAxisValue(b qBucket, axis int) uint8 {
	switch axis {
	case 0:
		return b.r5
	case 1:
		return b.g5
	default:
		return b.b5
	}
}

// qMedianCut splits the populated buckets into <= maxColors boxes; it writes
// each box's average color into palette and returns the palette together with a
// lookup table so lut[qPack(r,g,b)] is the palette index for any cell in a box.
func qMedianCut(pop []qBucket, maxColors int) (palette [][3]uint8, lut []int) {
	lut = make([]int, qCells)
	boxes := make([]qBox, 1, maxColors)
	boxes[0] = qBox{lo: 0, hi: len(pop)}

	for len(boxes) < maxColors {
		best, bestExt, bestAxis := -1, 0, 0
		for i := range boxes {
			if boxes[i].hi-boxes[i].lo < 2 {
				continue
			}
			axis, ext := qBoxExtent(pop, boxes[i])
			if ext > bestExt {
				bestExt = ext
				best = i
				bestAxis = axis
			}
		}
		if best < 0 {
			break // nothing left to split
		}

		axis := bestAxis
		slice := pop[boxes[best].lo:boxes[best].hi]
		sort.SliceStable(slice, func(a, b int) bool {
			return qAxisValue(slice[a], axis) < qAxisValue(slice[b], axis)
		})

		// Split at the median by pixel count.
		total := 0
		for i := boxes[best].lo; i < boxes[best].hi; i++ {
			total += pop[i].count
		}
		acc := 0
		split := boxes[best].lo + 1
		for i := boxes[best].lo; i < boxes[best].hi-1; i++ {
			acc += pop[i].count
			if acc*2 >= total {
				split = i + 1
				break
			}
		}

		newBox := qBox{lo: split, hi: boxes[best].hi}
		boxes[best].hi = split
		boxes = append(boxes, newBox)
	}

	palette = make([][3]uint8, len(boxes))
	for i := range boxes {
		var sr, sg, sb float64
		var cnt int
		for j := boxes[i].lo; j < boxes[i].hi; j++ {
			sr += pop[j].sr
			sg += pop[j].sg
			sb += pop[j].sb
			cnt += pop[j].count
		}
		if cnt == 0 {
			cnt = 1
		}
		palette[i] = [3]uint8{
			uint8(sr/float64(cnt) + 0.5),
			uint8(sg/float64(cnt) + 0.5),
			uint8(sb/float64(cnt) + 0.5),
		}
		for j := boxes[i].lo; j < boxes[i].hi; j++ {
			lut[qPack(pop[j].r5<<qShift, pop[j].g5<<qShift, pop[j].b5<<qShift)] = i
		}
	}
	return palette, lut
}

// qWithinLimit reports whether the image's RGB pixels contain at most maxColors
// distinct colors, bailing out as soon as one too many is seen.
func (img *Image) qWithinLimit(maxColors int) bool {
	seen := make(map[uint32]struct{}, maxColors+1)
	pix := img.data.Pix
	for i := 0; i+3 < len(pix); i += 4 {
		key := uint32(pix[i])<<16 | uint32(pix[i+1])<<8 | uint32(pix[i+2])
		if _, ok := seen[key]; ok {
			continue
		}
		if len(seen) >= maxColors {
			return false
		}
		seen[key] = struct{}{}
	}
	return true
}

// ReduceColors reduces the image to at most maxColors distinct RGB colors using
// median-cut quantization. When the image already has few enough colors it is
// left unchanged. Alpha is preserved. Mirrors C IReduceColors.
func (img *Image) ReduceColors(maxColors int) error {
	if maxColors < 1 {
		maxColors = 1
	}
	if maxColors > 256 {
		maxColors = 256
	}
	if img.width <= 0 || img.height <= 0 {
		return nil
	}
	if img.qWithinLimit(maxColors) {
		return nil
	}

	pix := img.data.Pix

	// Build the histogram of populated 5-bit cells.
	count := make([]int, qCells)
	var sumr, sumg, sumb [qCells]float64
	for i := 0; i+3 < len(pix); i += 4 {
		r, g, b := pix[i], pix[i+1], pix[i+2]
		cell := qPack(r, g, b)
		count[cell]++
		sumr[cell] += float64(r)
		sumg[cell] += float64(g)
		sumb[cell] += float64(b)
	}

	pop := make([]qBucket, 0, qCells)
	for i := 0; i < qCells; i++ {
		if count[i] == 0 {
			continue
		}
		pop = append(pop, qBucket{
			r5:    uint8((i >> (qBits * 2)) & (qLevels - 1)),
			g5:    uint8((i >> qBits) & (qLevels - 1)),
			b5:    uint8(i & (qLevels - 1)),
			count: count[i],
			sr:    sumr[i],
			sg:    sumg[i],
			sb:    sumb[i],
		})
	}
	if len(pop) == 0 {
		return nil
	}

	palette, lut := qMedianCut(pop, maxColors)
	if len(palette) == 0 {
		return nil
	}

	for i := 0; i+3 < len(pix); i += 4 {
		idx := lut[qPack(pix[i], pix[i+1], pix[i+2])]
		pix[i] = palette[idx][0]
		pix[i+1] = palette[idx][1]
		pix[i+2] = palette[idx][2]
	}
	return nil
}
