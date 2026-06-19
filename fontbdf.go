package ilibgo

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// maxGlyphDim bounds a single glyph's bounding-box width/height. Real BDF
// glyphs are at most a few hundred pixels; this cap rejects malformed files
// that would otherwise request an enormous bitmap allocation.
const maxGlyphDim = 4096

// TODO: support these translations
var charTranslations map[string]string = map[string]string{
	"space":        " ",
	"exclam":       "!",
	"quotedbl":     "\"",
	"numbersign":   "#",
	"dollar":       "$",
	"percent":      "%",
	"ampersand":    "&",
	"quoteright":   "\"",
	"parenleft":    "(",
	"parenright":   ")",
	"asterisk":     "*",
	"plus":         "+",
	"comma":        ",",
	"minus":        "-",
	"period":       ".",
	"slash":        "/",
	"zero":         "0",
	"one":          "1",
	"two":          "2",
	"three":        "3",
	"four":         "4",
	"five":         "5",
	"six":          "6",
	"seven":        "7",
	"eight":        "8",
	"nine":         "9",
	"colon":        ":",
	"semicolon":    ";",
	"less":         "<",
	"equal":        "=",
	"greater":      ">",
	"question":     "?",
	"at":           "@",
	"bracketleft":  "[",
	"backslash":    "\\",
	"bracketright": "]",
	"asciicircum":  "^",
	"underscore":   "_",
	"quoteleft":    "`",
	"braceleft":    "{",
	"bar":          "|",
	"braceright":   "}",
	"asciitilde":   "~",
}

type BdfChar struct {
	name        string // translated name ("A", "B", "\033agrave;")
	data        []bool // character definition
	width       int    // character width
	height      int    // height (size of lines array)
	actualWidth int    // width with padding on left and right
	xoffset     int    // pixels to move to the right before drawing
	yoffset     int    // pixels to move to up before drawing
}

type BdfFont struct {
	name         string // font name
	foundry      string // font foundry ("Adobe")
	family       string // font family ("Times", "Helvetica")
	faceName     string // face name ("Times Roman", "Helvetica Bold")
	widthName    string // face name ("Normal")
	slant        string // font slant ("R", "i")
	weight       string // font weight ("Normal", "Bold")
	proportional bool   // yes or no? (fixed/proportional font)
	pixelSize    int
	fontAscent   int
	fontDescent  int
	chars        [256]BdfChar // list of chars
	otherChars   []BdfChar    // non ASCII chars
}

// Parse a BFF font file.
// Example usage:
// LoadFontFromFile("timR12.bdf", "timR12")
func LoadFontFromFile(path string, name string) (*Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadFontFromBytes(name, data)
}

// LoadFontFromBytes parses a BDF font from its raw bytes, for example a file
// embedded with go:embed. It splits the data into lines (handling both "\n"
// and "\r\n") and delegates to LoadFontFromData.
func LoadFontFromBytes(name string, data []byte) (*Font, error) {
	var lines []string = make([]string, 0)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return LoadFontFromData(name, lines)
}

func trimQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

// bdfValue returns the text following the first keyword token on a BDF line,
// trimmed of surrounding whitespace. For `FONT_ASCENT 12` it returns "12"; for
// `FACE_NAME "Helvetica Bold"` it returns `"Helvetica Bold"`. If the line has
// no value it returns "".
func bdfValue(line string) string {
	_, rest, found := strings.Cut(line, " ")
	if !found {
		return ""
	}
	return strings.TrimSpace(rest)
}

// atoiOrZero parses a base-10 integer, returning 0 when the token is not a
// valid number. BDF integer fields are well-defined, so a parse failure means
// a malformed file; treating it as zero keeps parsing panic-free.
func atoiOrZero(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

// hexDigit converts a single hexadecimal digit (upper or lower case) to its
// value, or returns -1 if r is not a hex digit. Used when unpacking BDF bitmap
// rows.
func hexDigit(r rune) int {
	switch {
	case r >= '0' && r <= '9':
		return int(r - '0')
	case r >= 'A' && r <= 'F':
		return int(r-'A') + 10
	case r >= 'a' && r <= 'f':
		return int(r-'a') + 10
	}
	return -1
}

// Parse the array of string lines that represent the BDF font.
// Example usage:
// LoadFontFromData("timR12", font_timR12)
func LoadFontFromData(name string, lines []string) (*Font, error) {
	var font BdfFont

	var char BdfChar
	inBitmap := false
	yPos := 0
	for lineNo, line := range lines {
		// Identify the line by its leading keyword token. Splitting on
		// whitespace (rather than slicing at fixed byte offsets) tolerates
		// extra spaces and short/truncated lines without panicking.
		fields := strings.Fields(line)
		keyword := ""
		if len(fields) > 0 {
			keyword = fields[0]
		}

		switch {
		// Examples: "STARTCHAR W" or "STARTCHAR space"
		case keyword == "STARTCHAR":
			char = BdfChar{}
			char.name = bdfValue(line)
			if val, contains := charTranslations[char.name]; contains {
				// translate; e.g. replace "space" with " "
				char.name = val
			}
		case keyword == "ENCODING":
			if len(fields) >= 2 {
				int1, err := strconv.ParseInt(trimQuotes(fields[1]), 10, 32)
				if err == nil && int1 > 0 {
					char.name = fmt.Sprintf("%c", int1)
				}
			}
		case keyword == "PIXEL_SIZE":
			if len(fields) >= 2 {
				if int1, err := strconv.ParseInt(fields[1], 10, 32); err == nil {
					font.pixelSize = int(int1)
				}
			}
		case keyword == "FONT_ASCENT":
			if len(fields) >= 2 {
				if int1, err := strconv.ParseInt(fields[1], 10, 32); err == nil {
					font.fontAscent = int(int1)
				}
			}
		case keyword == "FONT_DESCENT":
			if len(fields) >= 2 {
				if int1, err := strconv.ParseInt(fields[1], 10, 32); err == nil {
					font.fontDescent = int(int1)
				}
			}
		case keyword == "SPACING":
			if len(fields) >= 2 {
				font.proportional = trimQuotes(fields[1]) == "P"
			}
		case keyword == "FACE_NAME": // "Times Italic", "Helvetica"
			font.faceName = trimQuotes(bdfValue(line))
		case keyword == "SETWIDTH_NAME": // "Normal"
			font.widthName = trimQuotes(bdfValue(line))
		case keyword == "WEIGHT_NAME": // "Medium"
			font.weight = trimQuotes(bdfValue(line))
		case keyword == "SLANT": // "R", "I"
			font.slant = trimQuotes(bdfValue(line))
		case keyword == "BBX":
			// BBX bbw bbh bbxoff bbyoff
			if len(fields) >= 5 {
				char.width = atoiOrZero(fields[1])
				char.height = atoiOrZero(fields[2])
				char.xoffset = atoiOrZero(fields[3])
				char.yoffset = atoiOrZero(fields[4])
			}
			// Reject nonsensical or absurdly large glyph boxes. Real BDF
			// glyphs are tiny; bounding this prevents a malformed file from
			// triggering a huge allocation (a DoS via make([]bool, w*h)).
			if char.width < 0 || char.height < 0 || char.width > maxGlyphDim || char.height > maxGlyphDim {
				return nil, fmt.Errorf("ilibgo: BDF parse: invalid BBX dimensions %dx%d at line %d", char.width, char.height, lineNo+1)
			}
			char.data = make([]bool, char.width*char.height)
		case keyword == "DWIDTH":
			// BDF DWIDTH is "dwx0 dwy0"; we only use the horizontal advance.
			if len(fields) >= 2 {
				char.actualWidth = atoiOrZero(fields[1])
			}
		case keyword == "BITMAP":
			inBitmap = true
			yPos = 0
		case keyword == "ENDCHAR":
			if inBitmap {
				inBitmap = false
				// Route by rune value: single-rune names with a code in the
				// [0,255] range (ASCII and Latin-1) index the fixed array;
				// everything else (multi-rune names, codes >= 256) goes to
				// otherChars. The < 256 guard also prevents an out-of-range
				// index into the [256]BdfChar array.
				r := []rune(char.name)
				if len(r) == 1 && r[0] < 256 {
					font.chars[r[0]] = char
				} else {
					font.otherChars = append(font.otherChars, char)
				}
			} else {
				err := fmt.Errorf("found ENDCHAR without STARTCHAR at line %d", (lineNo + 1))
				return nil, err
			}
			char = BdfChar{}
		case inBitmap:
			// A bitmap row is one hex digit per 4 pixels. Guard every index so
			// a truncated or malformed row can't panic.
			row := []rune(strings.TrimSpace(line))
			for x := 0; x < char.width && yPos < char.height; x++ {
				whichChar := x / 4
				if whichChar >= len(row) {
					break
				}
				hexval := hexDigit(row[whichChar])
				if hexval < 0 {
					continue
				}
				whichBit := 3 - (x % 4)
				idx := yPos*char.width + x
				if idx >= 0 && idx < len(char.data) {
					char.data[idx] = hexval&(1<<whichBit) > 0
				}
			}
			yPos++
		}
	}

	ret := Font{name: name, font: &font}

	return &ret, nil
}

// Get the BdfChar value for the specified letter.
// Note: Only tested on ASCII with English so far.
func FontBDFGetRune(font *BdfFont, ch rune) (char *BdfChar) {
	var character *BdfChar
	// find the character
	val := int(ch)
	character = &font.chars[val]
	return character
}

// Get the BdfChar value for the specified letter.
// Note: Only tested on ASCII with English so far.
func FontBDFGetChar(font *BdfFont, ch string) (char *BdfChar) {
	var character *BdfChar
	// find the character
	if len(ch) == 1 {
		val := int(ch[0])
		character = &font.chars[val]
	} else {
		//if more than a single char, must not be in the ascii list
		for _, val := range font.otherChars {
			if ch == val.name {
				character = &val
				break
			}
		}
	}
	// if still not there, use space
	if character == nil {
		character = &font.chars[' ']
	}
	return character
}
