package ilibgo_test

import (
	"bytes"
	"fmt"

	"github.com/craigk5n/ilibgo"
	font "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)

// Create a blank canvas with a solid background color.
func ExampleCreateImageWithBackground() {
	white, _ := ilibgo.AllocNamedColor("white")
	img := ilibgo.CreateImageWithBackground(320, 240, white)
	fmt.Printf("%dx%d\n", ilibgo.ImageWidth(img), ilibgo.ImageHeight(img))
	// Output: 320x240
}

// Look up a color by its X11 name.
func ExampleAllocNamedColor() {
	red, err := ilibgo.AllocNamedColor("red")
	if err != nil {
		panic(err)
	}
	r, g, b, a := red.RGBA()
	fmt.Println(r>>8, g>>8, b>>8, a>>8)
	// Output: 255 0 0 255
}

// Fill a rectangle using the method-style drawing API. The GraphicsContext
// carries the foreground color; pass it to the draw method.
func ExampleImage_FillRectangle() {
	white, _ := ilibgo.AllocNamedColor("white")
	blue, _ := ilibgo.AllocNamedColor("blue")
	img := ilibgo.CreateImageWithBackground(100, 100, white)

	gc := ilibgo.CreateGraphicsContext()
	ilibgo.SetForeground(&gc, blue)
	if err := img.FillRectangle(gc, 10, 10, 50, 40); err != nil {
		panic(err)
	}

	// The pixel at (30, 30) lies inside the filled rectangle.
	r, g, b, _ := img.GetPoint(30, 30).RGBA()
	fmt.Println(r>>8, g>>8, b>>8)
	// Output: 0 0 255
}

// Encode an image to any io.Writer. A PNG begins with an 8-byte signature.
func ExampleEncode() {
	white, _ := ilibgo.AllocNamedColor("white")
	img := ilibgo.CreateImageWithBackground(16, 16, white)

	var buf bytes.Buffer
	if err := ilibgo.Encode(&buf, img, ilibgo.FormatPNG); err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", buf.Bytes()[:8])
	// Output: 89504e470d0a1a0a
}

// Render text with a bundled X11 BDF font. The font is loaded from the
// embedded accessor, set on the GraphicsContext, then drawn.
func ExampleImage_DrawString() {
	white, _ := ilibgo.AllocNamedColor("white")
	black, _ := ilibgo.AllocNamedColor("black")
	img := ilibgo.CreateImageWithBackground(200, 40, white)

	f, err := ilibgo.LoadFontFromData("helvB12", font.Font_helvB12())
	if err != nil {
		panic(err)
	}
	gc := ilibgo.CreateGraphicsContext()
	ilibgo.SetForeground(&gc, black)
	ilibgo.SetFont(&gc, f)
	img.DrawString(gc, 10, 25, "Hello, ilibgo")

	w, _ := ilibgo.TextWidth(gc, f, "Hello, ilibgo")
	fmt.Println(w > 0)
	// Output: true
}
