// Command primitives renders a showcase of ilibgo's drawing primitives —
// anti-aliased shapes, smooth curves, and alpha compositing — to a PNG. It is
// the source of the "primitives" gallery image in the README.
//
//	go run ./primitives -outfile primitives.png
package main

import (
	"flag"
	"log"
	"math"

	"github.com/craigk5n/ilibgo"
	font "github.com/craigk5n/ilibgo/fonts/adobe_100dpi"
)

func mustColor(name string) ilibgo.Color {
	c, err := ilibgo.AllocNamedColor(name)
	if err != nil {
		log.Fatalf("color %q: %v", name, err)
	}
	return c
}

func main() {
	outfile := flag.String("outfile", "primitives.png", "output PNG file")
	flag.Parse()

	const w, h = 720, 320
	white := mustColor("white")
	img := ilibgo.CreateImageWithBackground(w, h, white)

	gc := ilibgo.CreateGraphicsContext()
	ilibgo.SetAntiAlias(&gc, true)

	drawRosette(img, gc, 120, 150, 95)
	drawCurves(img, gc, 240)
	drawCompositing(img, gc, 480)
	drawLabels(img)

	if err := ilibgo.SaveImageFile(*outfile, img, ilibgo.FormatPNG); err != nil {
		log.Fatalf("save %s: %v", *outfile, err)
	}
	log.Printf("wrote %s", *outfile)
}

// drawRosette draws a "mystic rose": every pair of points equally spaced around
// a circle is joined by an anti-aliased line, producing a smooth curved
// envelope from straight segments alone.
func drawRosette(img *ilibgo.Image, gc ilibgo.GraphicsContext, cx, cy, r int) {
	const n = 24
	pt := func(i int) (int, int) {
		a := 2 * math.Pi * float64(i) / n
		return cx + int(float64(r)*math.Cos(a)), cy + int(float64(r)*math.Sin(a))
	}
	rose, _ := ilibgo.AllocColorAlpha(70, 130, 180, 150) // translucent steelblue
	ilibgo.SetForeground(&gc, rose)
	ilibgo.SetBlendMode(&gc, ilibgo.BlendOver)
	for i := 0; i < n; i++ {
		x1, y1 := pt(i)
		for j := i + 1; j < n; j++ {
			x2, y2 := pt(j)
			img.DrawLine(gc, x1, y1, x2, y2)
		}
	}
	ilibgo.SetBlendMode(&gc, ilibgo.BlendReplace)
	ilibgo.SetForeground(&gc, mustColor("navy"))
	img.DrawCircle(gc, cx, cy, r)
}

// drawCurves draws a cubic Bezier and a Catmull-Rom spline, both anti-aliased.
func drawCurves(img *ilibgo.Image, gc ilibgo.GraphicsContext, x0 int) {
	ilibgo.SetForeground(&gc, mustColor("crimson"))
	img.DrawBezier(gc, []ilibgo.Point{
		ilibgo.Pt(x0+20, 230), ilibgo.Pt(x0+60, 40),
		ilibgo.Pt(x0+140, 250), ilibgo.Pt(x0+200, 70),
	})

	ilibgo.SetForeground(&gc, mustColor("darkgreen"))
	img.DrawSpline(gc, []ilibgo.Point{
		ilibgo.Pt(x0+20, 120), ilibgo.Pt(x0+60, 200), ilibgo.Pt(x0+100, 110),
		ilibgo.Pt(x0+150, 210), ilibgo.Pt(x0+200, 130),
	})

	// A filled pie wedge (anti-aliased arc fill).
	ilibgo.SetForeground(&gc, mustColor("goldenrod"))
	img.FillArc(gc, x0+120, 150, 36, 36, 20, 160)
}

// drawCompositing draws three translucent circles with source-over blending so
// the overlaps mix colors.
func drawCompositing(img *ilibgo.Image, gc ilibgo.GraphicsContext, x0 int) {
	ilibgo.SetBlendMode(&gc, ilibgo.BlendOver)
	red, _ := ilibgo.AllocColorAlpha(220, 50, 50, 140)
	green, _ := ilibgo.AllocColorAlpha(50, 180, 80, 140)
	blue, _ := ilibgo.AllocColorAlpha(60, 90, 220, 140)

	cx, cy, r := x0+120, 150, 60
	ilibgo.SetForeground(&gc, red)
	img.FillCircle(gc, cx, cy-32, r)
	ilibgo.SetForeground(&gc, green)
	img.FillCircle(gc, cx-36, cy+24, r)
	ilibgo.SetForeground(&gc, blue)
	img.FillCircle(gc, cx+36, cy+24, r)
}

func drawLabels(img *ilibgo.Image) {
	gc := ilibgo.CreateGraphicsContext()
	ilibgo.SetForeground(&gc, mustColor("black"))
	f, err := ilibgo.LoadFontFromData("helvB12", font.Font_helvB12())
	if err != nil {
		log.Fatalf("font: %v", err)
	}
	ilibgo.SetFont(&gc, f)
	img.DrawString(gc, 48, 300, "anti-aliased lines")
	img.DrawString(gc, 300, 300, "Bezier & spline curves")
	img.DrawString(gc, 540, 300, "alpha compositing")
}
