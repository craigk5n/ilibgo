package main

import (
	"bytes"
	"image"
	_ "image/png" // register PNG decoder for image.DecodeConfig
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigk5n/ilibgo"
)

func TestParseFrequency(t *testing.T) {
	in := "  2 hello\n 15 world\n\n   \nbareword\nx y z\n3 multi token\n"
	specs, err := parseFrequency(strings.NewReader(in))
	if err != nil {
		t.Fatalf("parseFrequency: %v", err)
	}
	want := []WordSpec{
		{Text: "hello", Count: 2},
		{Text: "world", Count: 15},
		{Text: "bareword", Count: 1}, // no leading count -> count 1
		{Text: "x", Count: 1},        // "x y z": first field not numeric -> word "x"
		{Text: "multi", Count: 3},    // "3 multi token" -> count 3, word "multi"
	}
	if len(specs) != len(want) {
		t.Fatalf("got %d specs, want %d: %+v", len(specs), len(want), specs)
	}
	for i, w := range want {
		if specs[i] != w {
			t.Errorf("spec[%d] = %+v, want %+v", i, specs[i], w)
		}
	}
}

func TestBuildWordsExplicitSize(t *testing.T) {
	words := buildWords([]WordSpec{
		{Text: "a", Size: 50},
		{Text: "", Size: 99}, // empty text dropped
	}, 10, 100)
	if len(words) != 1 {
		t.Fatalf("got %d words, want 1", len(words))
	}
	if words[0].size != 50 {
		t.Errorf("explicit size = %g, want 50", words[0].size)
	}
}

func TestBuildWordsSqrtScaling(t *testing.T) {
	const min, max = 10.0, 100.0
	words := buildWords([]WordSpec{
		{Text: "lo", Count: 1},
		{Text: "mid", Count: 50},
		{Text: "hi", Count: 100},
	}, min, max)

	byText := map[string]float64{}
	for _, w := range words {
		byText[w.text] = w.size
	}
	if byText["lo"] != min {
		t.Errorf("lowest count size = %g, want min %g", byText["lo"], min)
	}
	if byText["hi"] != max {
		t.Errorf("highest count size = %g, want max %g", byText["hi"], max)
	}
	// sqrt scaling: midpoint count is above the linear midpoint.
	frac := math.Sqrt(float64(50-1) / float64(100-1))
	wantMid := min + (max-min)*frac
	if math.Abs(byText["mid"]-wantMid) > 1e-9 {
		t.Errorf("mid size = %g, want %g", byText["mid"], wantMid)
	}
	if byText["mid"] <= (min+max)/2 {
		t.Errorf("sqrt scaling should put mid (%g) above linear midpoint %g", byText["mid"], (min+max)/2)
	}
}

func TestBuildWordsDegenerateCount(t *testing.T) {
	// All identical counts -> maxSize (avoid an all-tiny cloud).
	words := buildWords([]WordSpec{
		{Text: "a", Count: 5},
		{Text: "b", Count: 5},
	}, 10, 100)
	for _, w := range words {
		if w.size != 100 {
			t.Errorf("%q size = %g, want max 100", w.text, w.size)
		}
	}
}

func TestParseColor(t *testing.T) {
	cases := []struct {
		in      string
		r, g, b uint8
		wantErr bool
	}{
		{in: "#ff0000", r: 255},
		{in: "00ff00", g: 255},
		{in: "0000FF", b: 255},
		{in: "black"},
		{in: "white", r: 255, g: 255, b: 255},
		{in: "", wantErr: true},
		{in: "#xyz", wantErr: true},
		{in: "notacolor", wantErr: true},
		{in: "#12345", wantErr: true}, // wrong length
	}
	for _, c := range cases {
		col, err := parseColor(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("parseColor(%q): expected error", c.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseColor(%q): %v", c.in, err)
			continue
		}
		r, g, b, _ := col.RGBA()
		gr, gg, gb := uint8(r>>8), uint8(g>>8), uint8(b>>8)
		if gr != c.r || gg != c.g || gb != c.b {
			t.Errorf("parseColor(%q) = (%d,%d,%d), want (%d,%d,%d)", c.in, gr, gg, gb, c.r, c.g, c.b)
		}
	}
}

func TestSplitPalette(t *testing.T) {
	got := splitPalette("red, green ,, blue ,")
	want := []string{"red", "green", "blue"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d]=%q, want %q", i, got[i], want[i])
		}
	}
}

func TestOverlapsAny(t *testing.T) {
	boxes := []image.Rectangle{image.Rect(0, 0, 10, 10)}
	if !overlapsAny(image.Rect(5, 5, 15, 15), boxes) {
		t.Error("expected overlap")
	}
	if overlapsAny(image.Rect(20, 20, 30, 30), boxes) {
		t.Error("expected no overlap")
	}
}

func TestSpiralFindCentersFirst(t *testing.T) {
	canvas := image.Rect(0, 0, 200, 200)
	pt, ok := spiralFind(20, 10, 100, 100, canvas, nil)
	if !ok {
		t.Fatal("expected placement")
	}
	// First candidate is the center: box centered on (100,100).
	if pt.X != 100-10 || pt.Y != 100-5 {
		t.Errorf("first placement = %v, want centered (90,95)", pt)
	}
}

func TestSpiralFindAvoidsOverlap(t *testing.T) {
	canvas := image.Rect(0, 0, 200, 200)
	// Occupy the center so the next word must spiral away.
	occupied := []image.Rectangle{image.Rect(80, 80, 120, 120)}
	pt, ok := spiralFind(20, 10, 100, 100, canvas, occupied)
	if !ok {
		t.Fatal("expected placement")
	}
	placed := image.Rect(pt.X, pt.Y, pt.X+20, pt.Y+10)
	if placed.Overlaps(occupied[0]) {
		t.Errorf("placement %v overlaps occupied %v", placed, occupied[0])
	}
}

func TestSpiralFindNoRoom(t *testing.T) {
	canvas := image.Rect(0, 0, 10, 10)
	if _, ok := spiralFind(20, 20, 5, 5, canvas, nil); ok {
		t.Error("box larger than canvas should not be placed")
	}
}

func TestConfigValidate(t *testing.T) {
	good := defaultConfig()
	if err := good.validate(); err != nil {
		t.Errorf("default config should be valid: %v", err)
	}
	cases := map[string]func(c *Config){
		"zero width":    func(c *Config) { c.Width = 0 },
		"neg size":      func(c *Config) { c.MinSize = -1 },
		"min>max":       func(c *Config) { c.MinSize, c.MaxSize = 100, 10 },
		"zero dpi":      func(c *Config) { c.DPI = 0 },
		"empty palette": func(c *Config) { c.Palette = nil },
	}
	for name, mut := range cases {
		c := defaultConfig()
		mut(&c)
		if err := c.validate(); err == nil {
			t.Errorf("%s: expected validation error", name)
		}
	}
}

func TestLoadConfigFileOverlay(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "c.json")
	if err := os.WriteFile(path, []byte(`{"background":"white","width":123}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := defaultConfig()
	if err := loadConfigFile(&cfg, path); err != nil {
		t.Fatalf("loadConfigFile: %v", err)
	}
	if cfg.Background != "white" || cfg.Width != 123 {
		t.Errorf("overlay failed: bg=%q width=%d", cfg.Background, cfg.Width)
	}
	// Unspecified keys keep defaults.
	if cfg.Height != defaultConfig().Height {
		t.Errorf("height should keep default, got %d", cfg.Height)
	}
	if err := loadConfigFile(&cfg, filepath.Join(dir, "missing.json")); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFontLoaderCaches(t *testing.T) {
	load := newFontLoader("", 72)
	f1, err := load(24)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	f2, _ := load(24)
	if f1 != f2 {
		t.Error("same size should return cached font")
	}
	f3, _ := load(48)
	if f1 == f3 {
		t.Error("different size should return a different font")
	}
	if _, err := newFontLoader("does-not-exist.ttf", 72)(24); err == nil {
		t.Error("expected error loading missing font file")
	}
}

func TestLayoutAndRender(t *testing.T) {
	cfg := defaultConfig()
	cfg.Width, cfg.Height = 600, 400
	words := buildWords([]WordSpec{
		{Text: "big", Size: 60},
		{Text: "small", Size: 16},
		{Text: "fixed", Color: "red", Size: 24},
	}, cfg.MinSize, cfg.MaxSize)

	rng := rand.New(rand.NewSource(1))
	placed, skipped, err := layout(words, cfg, newFontLoader("", cfg.DPI), rng)
	if err != nil {
		t.Fatalf("layout: %v", err)
	}
	if len(placed) != 3 || len(skipped) != 0 {
		t.Fatalf("placed=%d skipped=%d, want 3/0", len(placed), len(skipped))
	}
	// Placed boxes must stay within the canvas.
	for _, p := range placed {
		if p.x < 0 || p.y < 0 || p.x >= cfg.Width || p.y >= cfg.Height {
			t.Errorf("%q baseline (%d,%d) out of bounds", p.text, p.x, p.y)
		}
	}
	img, err := render(cfg, placed)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if ilibgo.ImageWidth(img) != cfg.Width || ilibgo.ImageHeight(img) != cfg.Height {
		t.Errorf("image %dx%d, want %dx%d", ilibgo.ImageWidth(img), ilibgo.ImageHeight(img), cfg.Width, cfg.Height)
	}
}

func TestLayoutSkipsOversized(t *testing.T) {
	cfg := defaultConfig()
	cfg.Width, cfg.Height = 40, 30
	words := buildWords([]WordSpec{{Text: "enormous", Size: 200}}, cfg.MinSize, cfg.MaxSize)
	_, skipped, err := layout(words, cfg, newFontLoader("", cfg.DPI), rand.New(rand.NewSource(1)))
	if err != nil {
		t.Fatalf("layout: %v", err)
	}
	if len(skipped) != 1 {
		t.Errorf("expected 1 skipped word, got %d", len(skipped))
	}
}

func TestRunFrequencyMode(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.png")
	in := strings.NewReader("3 alpha\n9 beta\n1 gamma\n")
	if err := run([]string{"-out", out, "-w", "300", "-h", "200"}, in); err != nil {
		t.Fatalf("run: %v", err)
	}
	assertPNG(t, out, 300, 200)
}

func TestRunConfigMode(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "c.json")
	out := filepath.Join(dir, "out.png")
	json := `{"width":320,"height":240,"background":"white","palette":["navy"],
	          "words":[{"text":"one","size":40},{"text":"two","count":5}]}`
	if err := os.WriteFile(cfgPath, []byte(json), 0o644); err != nil {
		t.Fatal(err)
	}
	// Empty reader => no piped words; words come from config.
	if err := run([]string{"-config", cfgPath, "-out", out}, bytes.NewReader(nil)); err != nil {
		t.Fatalf("run: %v", err)
	}
	assertPNG(t, out, 320, 240)
}

func TestRunFlagOverridesConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "c.json")
	out := filepath.Join(dir, "out.png")
	if err := os.WriteFile(cfgPath, []byte(`{"width":100,"height":100,"words":[{"text":"x","size":20}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// -w overrides the config's width.
	if err := run([]string{"-config", cfgPath, "-out", out, "-w", "256", "-h", "256"}, bytes.NewReader(nil)); err != nil {
		t.Fatalf("run: %v", err)
	}
	assertPNG(t, out, 256, 256)
}

func TestRunNoWords(t *testing.T) {
	err := run([]string{"-out", filepath.Join(t.TempDir(), "x.png")}, bytes.NewReader(nil))
	if err == nil || !strings.Contains(err.Error(), "no words") {
		t.Errorf("expected no-words error, got %v", err)
	}
}

func TestRunInvalidConfigValue(t *testing.T) {
	in := strings.NewReader("1 a\n")
	err := run([]string{"-min", "100", "-max", "10"}, in)
	if err == nil || !strings.Contains(err.Error(), "minSize") {
		t.Errorf("expected validation error, got %v", err)
	}
}

// assertPNG checks that path is a PNG of the expected dimensions.
func assertPNG(t *testing.T, path string, w, h int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open output: %v", err)
	}
	defer f.Close()
	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if format != "png" {
		t.Errorf("format = %q, want png", format)
	}
	if cfg.Width != w || cfg.Height != h {
		t.Errorf("dimensions = %dx%d, want %dx%d", cfg.Width, cfg.Height, w, h)
	}
}
