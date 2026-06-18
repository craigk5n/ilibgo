// Package qr is a small, dependency-free QR Code encoder.
//
// It encodes data in byte mode (which can represent any input), selects the
// smallest QR version that fits for the chosen error-correction level, applies
// Reed-Solomon error correction over GF(256), and chooses the mask pattern
// with the lowest penalty. The result is a square grid of dark/light modules.
//
// The algorithm follows the QR Code standard (ISO/IEC 18004) and the structure
// of Nayuki's well-known reference encoder.
package qr

import "fmt"

// Ecc is a QR error-correction level (higher tolerates more damage but holds
// less data).
type Ecc int

const (
	Low      Ecc = iota // recovers ~7%
	Medium              // recovers ~15%
	Quartile            // recovers ~25%
	High                // recovers ~30%
)

func (e Ecc) formatBits() int {
	return [...]int{1, 0, 3, 2}[e]
}

// Code is an encoded QR symbol: a size x size grid of modules.
type Code struct {
	version    int
	size       int
	modules    [][]bool // true = dark
	isFunction [][]bool
}

// Version returns the QR version (1..40); the symbol is 17+4*version modules
// on a side.
func (c *Code) Version() int { return c.version }

// Size returns the side length in modules (excluding the quiet zone).
func (c *Code) Size() int { return c.size }

// Module reports whether the module at (x, y) is dark.
func (c *Code) Module(x, y int) bool { return c.modules[y][x] }

// EncodeText encodes s (as UTF-8 bytes) at the given error-correction level.
func EncodeText(s string, ecl Ecc) (*Code, error) {
	return Encode([]byte(s), ecl)
}

// Encode encodes data in byte mode at the given error-correction level.
func Encode(data []byte, ecl Ecc) (*Code, error) {
	if ecl < Low || ecl > High {
		return nil, fmt.Errorf("qr: invalid error-correction level %d", int(ecl))
	}

	// Find the smallest version that can hold the data.
	version := 0
	for v := 1; v <= 40; v++ {
		capacityBits := getNumDataCodewords(v, ecl) * 8
		if 4+charCountBits(v)+len(data)*8 <= capacityBits {
			version = v
			break
		}
	}
	if version == 0 {
		return nil, fmt.Errorf("qr: data of %d bytes is too long for any version at this error-correction level", len(data))
	}

	// Build the bit stream: mode indicator, char count, data, terminator, padding.
	var bb bitBuffer
	bb.appendBits(0x4, 4) // byte mode
	bb.appendBits(uint(len(data)), charCountBits(version))
	for _, b := range data {
		bb.appendBits(uint(b), 8)
	}
	capacityBits := getNumDataCodewords(version, ecl) * 8
	bb.appendBits(0, min(4, capacityBits-len(bb)))
	bb.appendBits(0, (8-len(bb)%8)%8)
	for pad := 0xEC; len(bb) < capacityBits; pad ^= 0xEC ^ 0x11 {
		bb.appendBits(uint(pad), 8)
	}

	dataCodewords := bb.bytes()
	allCodewords := addEccAndInterleave(dataCodewords, version, ecl)

	c := &Code{version: version, size: version*4 + 17}
	c.modules = make([][]bool, c.size)
	c.isFunction = make([][]bool, c.size)
	for i := range c.modules {
		c.modules[i] = make([]bool, c.size)
		c.isFunction[i] = make([]bool, c.size)
	}

	c.drawFunctionPatterns(ecl)
	c.drawCodewords(allCodewords)

	// Pick the mask with the lowest penalty.
	bestMask, minPenalty := 0, -1
	for m := 0; m < 8; m++ {
		c.applyMask(m)
		c.drawFormatBits(ecl, m)
		p := c.penaltyScore()
		if minPenalty == -1 || p < minPenalty {
			minPenalty = p
			bestMask = m
		}
		c.applyMask(m) // XOR again to undo
	}
	c.applyMask(bestMask)
	c.drawFormatBits(ecl, bestMask)

	return c, nil
}

func charCountBits(version int) int {
	// Byte mode character-count field width.
	if version <= 9 {
		return 8
	}
	return 16
}

// --- Bit buffer -----------------------------------------------------------

type bitBuffer []bool

func (bb *bitBuffer) appendBits(val uint, n int) {
	for i := n - 1; i >= 0; i-- {
		*bb = append(*bb, (val>>uint(i))&1 != 0)
	}
}

func (bb bitBuffer) bytes() []byte {
	out := make([]byte, len(bb)/8)
	for i, bit := range bb {
		if bit {
			out[i>>3] |= 1 << uint(7-(i&7))
		}
	}
	return out
}

// --- Capacity / structure tables ------------------------------------------

// getNumRawDataModules returns the number of data-and-EC modules available in
// a QR symbol of the given version (i.e. excluding all function patterns).
func getNumRawDataModules(ver int) int {
	size := ver*4 + 17
	result := size * size
	result -= 8 * 8 * 3       // three finder patterns with separators
	result -= 15*2 + 1        // two format-info strips and the dark module
	result -= (size - 16) * 2 // two timing patterns
	if ver >= 2 {
		numAlign := ver/7 + 2
		result -= (numAlign - 1) * (numAlign - 1) * 25 // alignment patterns
		result -= (numAlign - 2) * 2 * 20              // overlap with timing patterns
		if ver >= 7 {
			result -= 6 * 3 * 2 // two version-info blocks
		}
	}
	return result
}

func getNumDataCodewords(ver int, ecl Ecc) int {
	return getNumRawDataModules(ver)/8 -
		eccCodewordsPerBlock[ecl][ver]*numEccBlocks[ecl][ver]
}

// Indexed [ecl][version]; index 0 (version) is unused padding.
var eccCodewordsPerBlock = [4][41]int{
	{-1, 7, 10, 15, 20, 26, 18, 20, 24, 30, 18, 20, 24, 26, 30, 22, 24, 28, 30, 28, 28, 28, 28, 30, 30, 26, 28, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30},
	{-1, 10, 16, 26, 18, 24, 16, 18, 22, 22, 26, 30, 22, 22, 24, 24, 28, 28, 26, 26, 26, 26, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28},
	{-1, 13, 22, 18, 26, 18, 24, 18, 22, 20, 24, 28, 26, 24, 20, 30, 24, 28, 28, 26, 30, 28, 30, 30, 30, 30, 28, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30},
	{-1, 17, 28, 22, 16, 22, 28, 26, 26, 24, 28, 24, 28, 22, 24, 24, 30, 28, 28, 26, 28, 30, 24, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30},
}

var numEccBlocks = [4][41]int{
	{-1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 4, 4, 4, 4, 4, 6, 6, 6, 6, 7, 8, 8, 9, 9, 10, 12, 12, 12, 13, 14, 15, 16, 17, 18, 19, 19, 20, 21, 22, 24, 25},
	{-1, 1, 1, 1, 2, 2, 4, 4, 4, 5, 5, 5, 8, 9, 9, 10, 10, 11, 13, 14, 16, 17, 17, 18, 20, 21, 23, 25, 26, 28, 29, 31, 33, 35, 37, 38, 40, 43, 45, 47, 49},
	{-1, 1, 1, 2, 2, 4, 4, 6, 6, 8, 8, 8, 10, 12, 16, 12, 17, 16, 18, 21, 20, 23, 23, 25, 27, 29, 34, 34, 35, 38, 40, 43, 45, 48, 51, 53, 56, 59, 62, 65, 68},
	{-1, 1, 1, 2, 4, 4, 4, 5, 6, 8, 8, 11, 11, 16, 16, 18, 16, 19, 21, 25, 25, 25, 34, 30, 32, 35, 37, 40, 42, 45, 48, 51, 54, 57, 60, 63, 66, 70, 74, 77, 81},
}

// --- Reed-Solomon over GF(256) --------------------------------------------

func reedSolomonMultiply(x, y int) int {
	z := 0
	for i := 7; i >= 0; i-- {
		z = (z << 1) ^ ((z >> 7) * 0x11D)
		z ^= ((y >> uint(i)) & 1) * x
	}
	return z & 0xFF
}

func reedSolomonDivisor(degree int) []int {
	result := make([]int, degree)
	result[degree-1] = 1
	root := 1
	for i := 0; i < degree; i++ {
		for j := 0; j < degree; j++ {
			result[j] = reedSolomonMultiply(result[j], root)
			if j+1 < degree {
				result[j] ^= result[j+1]
			}
		}
		root = reedSolomonMultiply(root, 0x02)
	}
	return result
}

func reedSolomonRemainder(data, divisor []int) []int {
	result := make([]int, len(divisor))
	for _, b := range data {
		factor := b ^ result[0]
		copy(result, result[1:])
		result[len(result)-1] = 0
		for i := range result {
			result[i] ^= reedSolomonMultiply(divisor[i], factor)
		}
	}
	return result
}

func addEccAndInterleave(data []byte, ver int, ecl Ecc) []byte {
	numBlocks := numEccBlocks[ecl][ver]
	blockEccLen := eccCodewordsPerBlock[ecl][ver]
	rawCodewords := getNumRawDataModules(ver) / 8
	numShortBlocks := numBlocks - rawCodewords%numBlocks
	shortBlockLen := rawCodewords / numBlocks

	divisor := reedSolomonDivisor(blockEccLen)
	blocks := make([][]int, numBlocks)
	k := 0
	for i := 0; i < numBlocks; i++ {
		datLen := shortBlockLen - blockEccLen
		if i >= numShortBlocks {
			datLen++
		}
		dat := make([]int, datLen)
		for j := 0; j < datLen; j++ {
			dat[j] = int(data[k+j])
		}
		k += datLen
		// block is sized shortBlockLen+1; short blocks leave a padding slot.
		block := make([]int, shortBlockLen+1)
		copy(block, dat)
		ecc := reedSolomonRemainder(dat, divisor)
		copy(block[len(block)-blockEccLen:], ecc)
		blocks[i] = block
	}

	result := make([]byte, 0, rawCodewords)
	for i := 0; i < len(blocks[0]); i++ {
		for j := 0; j < numBlocks; j++ {
			// Skip the padding slot in short blocks.
			if i != shortBlockLen-blockEccLen || j >= numShortBlocks {
				result = append(result, byte(blocks[j][i]))
			}
		}
	}
	return result
}

// --- Matrix drawing -------------------------------------------------------

func (c *Code) setFunctionModule(x, y int, isDark bool) {
	c.modules[y][x] = isDark
	c.isFunction[y][x] = true
}

func (c *Code) drawFunctionPatterns(ecl Ecc) {
	for i := 0; i < c.size; i++ {
		c.setFunctionModule(6, i, i%2 == 0)
		c.setFunctionModule(i, 6, i%2 == 0)
	}
	c.drawFinderPattern(3, 3)
	c.drawFinderPattern(c.size-4, 3)
	c.drawFinderPattern(3, c.size-4)

	pos := c.alignmentPatternPositions()
	n := len(pos)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if (i == 0 && j == 0) || (i == 0 && j == n-1) || (i == n-1 && j == 0) {
				continue // overlaps finder patterns
			}
			c.drawAlignmentPattern(pos[i], pos[j])
		}
	}

	c.drawFormatBits(ecl, 0) // reserve the area; overwritten after masking
	c.drawVersion()
}

func (c *Code) drawFinderPattern(x, y int) {
	for dy := -4; dy <= 4; dy++ {
		for dx := -4; dx <= 4; dx++ {
			dist := absMax(dx, dy)
			xx, yy := x+dx, y+dy
			if 0 <= xx && xx < c.size && 0 <= yy && yy < c.size {
				c.setFunctionModule(xx, yy, dist != 2 && dist != 4)
			}
		}
	}
}

func (c *Code) drawAlignmentPattern(x, y int) {
	for dy := -2; dy <= 2; dy++ {
		for dx := -2; dx <= 2; dx++ {
			c.setFunctionModule(x+dx, y+dy, absMax(dx, dy) != 1)
		}
	}
}

func (c *Code) drawFormatBits(ecl Ecc, mask int) {
	data := ecl.formatBits()<<3 | mask
	rem := data
	for i := 0; i < 10; i++ {
		rem = (rem << 1) ^ ((rem >> 9) * 0x537)
	}
	bits := (data<<10 | rem) ^ 0x5412 // 15 bits

	for i := 0; i <= 5; i++ {
		c.setFunctionModule(8, i, getBit(bits, i))
	}
	c.setFunctionModule(8, 7, getBit(bits, 6))
	c.setFunctionModule(8, 8, getBit(bits, 7))
	c.setFunctionModule(7, 8, getBit(bits, 8))
	for i := 9; i < 15; i++ {
		c.setFunctionModule(14-i, 8, getBit(bits, i))
	}

	for i := 0; i < 8; i++ {
		c.setFunctionModule(c.size-1-i, 8, getBit(bits, i))
	}
	for i := 8; i < 15; i++ {
		c.setFunctionModule(8, c.size-15+i, getBit(bits, i))
	}
	c.setFunctionModule(8, c.size-8, true) // always-dark module
}

func (c *Code) drawVersion() {
	if c.version < 7 {
		return
	}
	rem := c.version
	for i := 0; i < 12; i++ {
		rem = (rem << 1) ^ ((rem >> 11) * 0x1F25)
	}
	bits := c.version<<12 | rem // 18 bits
	for i := 0; i < 18; i++ {
		bit := getBit(bits, i)
		a, b := c.size-11+i%3, i/3
		c.setFunctionModule(a, b, bit)
		c.setFunctionModule(b, a, bit)
	}
}

func (c *Code) drawCodewords(data []byte) {
	i := 0 // bit index into data
	for right := c.size - 1; right >= 1; right -= 2 {
		if right == 6 {
			right = 5 // skip the vertical timing column
		}
		for vert := 0; vert < c.size; vert++ {
			for j := 0; j < 2; j++ {
				x := right - j
				upward := (right+1)&2 == 0
				y := vert
				if upward {
					y = c.size - 1 - vert
				}
				if !c.isFunction[y][x] && i < len(data)*8 {
					c.modules[y][x] = getBit(int(data[i>>3]), 7-(i&7))
					i++
				}
			}
		}
	}
}

func (c *Code) applyMask(mask int) {
	for y := 0; y < c.size; y++ {
		for x := 0; x < c.size; x++ {
			if c.isFunction[y][x] {
				continue
			}
			var invert bool
			switch mask {
			case 0:
				invert = (x+y)%2 == 0
			case 1:
				invert = y%2 == 0
			case 2:
				invert = x%3 == 0
			case 3:
				invert = (x+y)%3 == 0
			case 4:
				invert = (x/3+y/2)%2 == 0
			case 5:
				invert = x*y%2+x*y%3 == 0
			case 6:
				invert = (x*y%2+x*y%3)%2 == 0
			case 7:
				invert = ((x+y)%2+x*y%3)%2 == 0
			}
			if invert {
				c.modules[y][x] = !c.modules[y][x]
			}
		}
	}
}

// --- Penalty scoring (mask selection) -------------------------------------

const (
	penaltyN1 = 3
	penaltyN2 = 3
	penaltyN3 = 40
	penaltyN4 = 10
)

func (c *Code) penaltyScore() int {
	result := 0
	size := c.size

	// Rule 1 + finder-like rule 3, by rows.
	for y := 0; y < size; y++ {
		runColor := false
		runX := 0
		var hist [7]int
		for x := 0; x < size; x++ {
			if c.modules[y][x] == runColor {
				runX++
				if runX == 5 {
					result += penaltyN1
				} else if runX > 5 {
					result++
				}
			} else {
				c.finderPenaltyAddHistory(runX, &hist)
				if !runColor {
					result += c.finderPenaltyCount(hist) * penaltyN3
				}
				runColor = c.modules[y][x]
				runX = 1
			}
		}
		result += c.finderPenaltyTerminate(runColor, runX, &hist) * penaltyN3
	}
	// Rule 1 + 3, by columns.
	for x := 0; x < size; x++ {
		runColor := false
		runY := 0
		var hist [7]int
		for y := 0; y < size; y++ {
			if c.modules[y][x] == runColor {
				runY++
				if runY == 5 {
					result += penaltyN1
				} else if runY > 5 {
					result++
				}
			} else {
				c.finderPenaltyAddHistory(runY, &hist)
				if !runColor {
					result += c.finderPenaltyCount(hist) * penaltyN3
				}
				runColor = c.modules[y][x]
				runY = 1
			}
		}
		result += c.finderPenaltyTerminate(runColor, runY, &hist) * penaltyN3
	}

	// Rule 2: 2x2 blocks of the same color.
	for y := 0; y < size-1; y++ {
		for x := 0; x < size-1; x++ {
			color := c.modules[y][x]
			if color == c.modules[y][x+1] && color == c.modules[y+1][x] && color == c.modules[y+1][x+1] {
				result += penaltyN2
			}
		}
	}

	// Rule 4: overall dark/light balance.
	dark := 0
	for _, row := range c.modules {
		for _, color := range row {
			if color {
				dark++
			}
		}
	}
	total := size * size
	k := (abs(dark*20-total*10)+total-1)/total - 1
	result += k * penaltyN4

	return result
}

func (c *Code) finderPenaltyCount(hist [7]int) int {
	n := hist[1]
	core := n > 0 && hist[2] == n && hist[3] == n*3 && hist[4] == n && hist[5] == n
	res := 0
	if core && hist[0] >= n*4 && hist[6] >= n {
		res++
	}
	if core && hist[6] >= n*4 && hist[0] >= n {
		res++
	}
	return res
}

func (c *Code) finderPenaltyTerminate(currentRunColor bool, currentRunLength int, hist *[7]int) int {
	if currentRunColor { // dark run: add to history, then count the trailing light border
		c.finderPenaltyAddHistory(currentRunLength, hist)
		currentRunLength = 0
	}
	currentRunLength += c.size // light border
	c.finderPenaltyAddHistory(currentRunLength, hist)
	return c.finderPenaltyCount(*hist)
}

func (c *Code) finderPenaltyAddHistory(currentRunLength int, hist *[7]int) {
	if hist[0] == 0 {
		currentRunLength += c.size // add the leading light border
	}
	for i := 6; i >= 1; i-- {
		hist[i] = hist[i-1]
	}
	hist[0] = currentRunLength
}

// alignmentPatternPositions returns the center coordinates of the alignment
// patterns for this version.
func (c *Code) alignmentPatternPositions() []int {
	ver := c.version
	if ver == 1 {
		return nil
	}
	numAlign := ver/7 + 2
	step := 26
	if ver != 32 {
		step = (ver*4 + numAlign*2 + 1) / (numAlign*2 - 2) * 2
	}
	result := make([]int, numAlign)
	result[0] = 6
	pos := c.size - 7
	for i := numAlign - 1; i >= 1; i-- {
		result[i] = pos
		pos -= step
	}
	return result
}

// --- small helpers --------------------------------------------------------

func getBit(x, i int) bool { return (x>>uint(i))&1 != 0 }

func absMax(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
