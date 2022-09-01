package ilibgo

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

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
	var lines []string = make([]string, 0)
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	fileScanner := bufio.NewScanner(fp)
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	ret, err := LoadFontFromData(name, lines)
	if err != nil {
		return nil, err
	}
	return ret, nil
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

// Parse the array of string lines that represent the BDF font.
// Example usage:
// LoadFontFromData("timR12", font_timR12)
func LoadFontFromData(name string, lines []string) (*Font, error) {
	var font BdfFont

	var char BdfChar
	inBitmap := false
	var temp int
	xPos := 0
	yPos := 0
	for lineNo, line := range lines {
		// Examples: "STARTCHAR W" or "STARTCHAR space"
		if strings.HasPrefix(line, "STARTCHAR") {
			char = BdfChar{}
			char.name = line[10:]
			if val, contains := charTranslations[char.name]; contains {
				// translate; e.g. replace "space" with " "
				char.name = val
			}
		} else if strings.HasPrefix(line, "ENCODING") {
			temp := trimQuotes(line[9:])
			int1, err := strconv.ParseInt(temp, 10, 32)
			if err == nil && int1 > 0 && int1 < 256 {
				char.name = fmt.Sprintf("%c", int1)
			} else if err == nil && int1 > 0 {
				char.name = fmt.Sprintf("%c", int1)
			}
		} else if strings.HasPrefix(line, "PIXEL_SIZE") {
			temp := line[11:]
			int1, err := strconv.ParseInt(temp, 10, 32)
			if err == nil {
				font.pixelSize = int(int1)
			}
		} else if strings.HasPrefix(line, "PIXEL_SIZE") {
			temp := line[11:]
			int1, err := strconv.ParseInt(temp, 10, 32)
			if err == nil {
				font.pixelSize = int(int1)
			}
		} else if strings.HasPrefix(line, "FONT_ASCENT") {
			temp := line[12:]
			int1, err := strconv.ParseInt(temp, 10, 32)
			if err == nil {
				font.fontAscent = int(int1)
			}
		} else if strings.HasPrefix(line, "FONT_DESCENT") {
			temp := line[12:]
			int1, err := strconv.ParseInt(temp, 10, 32)
			if err == nil {
				font.fontDescent = int(int1)
			}
		} else if strings.HasPrefix(line, "SPACING") {
			font.proportional = (line == "SPACING \"P\"")
		} else if strings.HasPrefix(line, "FACE_NAME") { // "Times Italic", "Helvetica"
			font.faceName = trimQuotes(line[10:])
		} else if strings.HasPrefix(line, "SETWIDTH_NAME") { // "Normal"
			font.widthName = trimQuotes(line[14:])
		} else if strings.HasPrefix(line, "WEIGHT_NAME") { // "Medium"
			font.weight = trimQuotes(line[12:])
		} else if strings.HasPrefix(line, "SLANT") { // "R", "I"
			font.weight = trimQuotes(line[6:])
		} else if strings.HasPrefix(line, "BBX") {
			fmt.Sscanf(line, "BBX %d %d %d %d", &char.width, &char.height, &char.xoffset, &char.yoffset)
			char.data = make([]bool, char.width*char.height)
		} else if strings.HasPrefix(line, "DWIDTH") {
			fmt.Sscanf(line, "DWIDTH %d %d", &char.actualWidth, &char.height, &char.xoffset, &temp)
		} else if strings.HasPrefix(line, "BITMAP") {
			inBitmap = true
			xPos = 0
			yPos = 0
		} else if line == "ENDCHAR" {
			if inBitmap {
				inBitmap = false
				if len(char.name) == 1 {
					r := []rune(char.name)
					font.chars[r[0]] = char
				} else {
					font.otherChars = append(font.otherChars, char)
				}
			} else {
				err := fmt.Errorf("found ENDCHAR without STARTCHAR at line %d", (lineNo + 1))
				return nil, err
			}
			char = BdfChar{}
		} else if inBitmap {
			xPos = 0
			for loop := 0; loop < char.width; loop++ {
				r := []rune(line)
				whichChar := loop / 4
				c := r[whichChar]
				var hexval int
				if unicode.IsDigit(c) {
					hexval = (int)(c - '0')
				} else {
					hexval = (int)(c-'A') + 10
				}
				whichBit := 3 - (loop % 4)
				if hexval&(1<<whichBit) > 0 {
					char.data[yPos*char.width+xPos] = true
					//fmt.Print("X")
				} else {
					char.data[yPos*char.width+xPos] = false
					//fmt.Print(" ")
				}
				xPos++
			}
			//fmt.Println("")
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
