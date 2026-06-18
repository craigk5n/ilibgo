// captcha generates a distorted-text CAPTCHA image and prints the text it
// encodes. Each character is drawn at a random angle, position and color over a
// background of random noise lines and speckle.
//
// Usage:
//
//	captcha [options]
//
//	  -out string    output PNG file (default "captcha.png")
//	  -text string   text to render (default: a random code)
//	  -len int       length of the random code when -text is unset (default 6)
//	  -w int         image width (default 240)
//	  -h int         image height (default 90)
//	  -seed int      RNG seed; 0 uses a time-based seed (default 0)
//
// The generated text is written to stdout so the caller can check answers.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/craigk5n/ilibgo"
	font "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)

// Unambiguous character set (no 0/O, 1/I/L, etc.).
const charset = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"

func main() {
	out := flag.String("out", "captcha.png", "output PNG file")
	text := flag.String("text", "", "text to render (default: random)")
	length := flag.Int("len", 6, "length of the random code when -text is unset")
	width := flag.Int("w", 240, "image width")
	height := flag.Int("h", 90, "image height")
	seed := flag.Int64("seed", 0, "RNG seed (0 = time-based)")
	flag.Parse()

	s := *seed
	if s == 0 {
		s = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(s))

	code := *text
	if code == "" {
		code = randomCode(rng, *length)
	}

	img, err := render(code, *width, *height, rng)
	if err != nil {
		fail(err)
	}

	f, err := os.Create(*out)
	if err != nil {
		fail(fmt.Errorf("create %s: %w", *out, err))
	}
	defer f.Close()
	if err := ilibgo.WriteImageFile(f, img, ilibgo.FormatPNG); err != nil {
		fail(fmt.Errorf("write %s: %w", *out, err))
	}
	fmt.Println(code)
}

func randomCode(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

func render(text string, w, h int, rng *rand.Rand) (*ilibgo.Image, error) {
	// Light, slightly random background.
	bg, _ := ilibgo.AllocColor(uint8(230+rng.Intn(26)), uint8(230+rng.Intn(26)), uint8(230+rng.Intn(26)))
	img := ilibgo.CreateImageWithBackground(w, h, bg)

	gc := ilibgo.CreateGraphicsContext()
	f, err := ilibgo.LoadFontFromData("helvB24", font.Font_helvB24())
	if err != nil {
		return nil, err
	}
	ilibgo.SetFont(&gc, f)
	fontHeight, _ := ilibgo.GetFontSize(f)

	// Noise lines behind the text.
	for i := 0; i < 6; i++ {
		ilibgo.SetForeground(&gc, midColor(rng))
		img.DrawLine(gc, rng.Intn(w), rng.Intn(h), rng.Intn(w), rng.Intn(h))
	}

	// Characters: each rotated and vertically jittered.
	margin := 10
	usable := w - 2*margin
	step := usable / len(text)
	x := margin
	for _, ch := range text {
		ilibgo.SetForeground(&gc, darkColor(rng))
		angle := rng.Float64()*40 - 20 // -20..+20 degrees
		yJitter := rng.Intn(fontHeight/3+1) - fontHeight/6
		y := h/2 + fontHeight/3 + yJitter
		img.DrawStringRotatedAngle(gc, x, y, string(ch), angle)
		x += step + rng.Intn(5) - 2
	}

	// Speckle noise.
	for i := 0; i < (w*h)/40; i++ {
		ilibgo.SetForeground(&gc, midColor(rng))
		img.SetPoint(gc, rng.Intn(w), rng.Intn(h))
	}

	return img, nil
}

// darkColor returns a dark, saturated-ish color for readable glyphs.
func darkColor(rng *rand.Rand) ilibgo.Color {
	c, _ := ilibgo.AllocColor(uint8(rng.Intn(110)), uint8(rng.Intn(110)), uint8(rng.Intn(110)))
	return c
}

// midColor returns a mid-tone color for noise.
func midColor(rng *rand.Rand) ilibgo.Color {
	c, _ := ilibgo.AllocColor(uint8(120+rng.Intn(80)), uint8(120+rng.Intn(80)), uint8(120+rng.Intn(80)))
	return c
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "captcha: %v\n", err)
	os.Exit(1)
}
