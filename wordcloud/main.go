// wordcloud renders a word cloud to a PNG image.
//
// Words come from two sources, which may be combined:
//
//   - a JSON config file (-config) that carries render settings and may list
//     words with explicit per-word sizes and colors;
//
//   - "count word" frequency lines read from stdin or file arguments — the
//     output of a pipeline such as:
//
//     cat *.txt | tr ' ' '\012' | sort | uniq -c | sort -n
//
// A word's point size is taken from its explicit "size" when given; otherwise
// it is derived from its count using sqrt scaling between -min and -max. Words
// are placed largest-first with Archimedean-spiral packing so they do not
// overlap, colored randomly from -palette (unless a per-word color is set), on
// a -bg background.
//
// Usage:
//
//	wordcloud [options] [freqfile ...]
//
//	  -config string    JSON config file
//	  -out string       output PNG file (default "wordcloud.png")
//	  -font string      TrueType/OpenType font path (default: built-in goregular)
//	  -w, -h int        image dimensions (default 800x600)
//	  -bg string        background color name or #rrggbb (default "black")
//	  -palette string   comma-separated word colors
//	  -min, -max float  min/max point size (default 14, 96)
//	  -dpi float         rendering DPI (default 72)
//	  -seed int         random seed for color/placement reproducibility (default 1)
//
// Examples:
//
//	wordcloud -config cloud.json
//	cat *.txt | tr ' ' '\012' | sort | uniq -c | sort -n | wordcloud -out cloud.png
//	wordcloud -bg white -palette "navy,teal,crimson" freq.txt
package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"

	"github.com/craigk5n/ilibgo"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	if err := run(os.Args[1:], os.Stdin); err != nil {
		fmt.Fprintln(os.Stderr, "wordcloud:", err)
		os.Exit(1)
	}
}

// run wires up flags, config, and rendering. It is separated from main so tests
// can exercise it with explicit args and input.
func run(argv []string, stdin io.Reader) error {
	fs := flag.NewFlagSet("wordcloud", flag.ContinueOnError)
	configPath := fs.String("config", "", "JSON config file")
	out := fs.String("out", "wordcloud.png", "output PNG file")
	font := fs.String("font", "", "TrueType/OpenType font path (default: built-in goregular)")
	width := fs.Int("w", 0, "image width")
	height := fs.Int("h", 0, "image height")
	bg := fs.String("bg", "", "background color name or #rrggbb")
	palette := fs.String("palette", "", "comma-separated word colors")
	minSize := fs.Float64("min", 0, "minimum font point size")
	maxSize := fs.Float64("max", 0, "maximum font point size")
	dpi := fs.Float64("dpi", 0, "rendering DPI")
	seed := fs.Int64("seed", 0, "random seed for reproducible output")
	if err := fs.Parse(argv); err != nil {
		return err
	}

	cfg := defaultConfig()
	if *configPath != "" {
		if err := loadConfigFile(&cfg, *configPath); err != nil {
			return err
		}
	}
	applyFlags(&cfg, fs, flagValues{
		out: out, font: font, width: width, height: height, bg: bg,
		palette: palette, minSize: minSize, maxSize: maxSize, dpi: dpi, seed: seed,
	})
	if err := cfg.validate(); err != nil {
		return err
	}

	specs, err := gatherWords(cfg.Words, fs.Args(), stdin)
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return fmt.Errorf("no words: provide words in -config or pipe \"count word\" lines on stdin")
	}

	words := buildWords(specs, cfg.MinSize, cfg.MaxSize)
	load := newFontLoader(cfg.Font, cfg.DPI)
	rng := rand.New(rand.NewSource(cfg.Seed))

	placed, skipped, err := layout(words, cfg, load, rng)
	if err != nil {
		return err
	}
	if len(skipped) > 0 {
		fmt.Fprintf(os.Stderr, "wordcloud: %d word(s) did not fit and were skipped: %s\n",
			len(skipped), strings.Join(skipped, ", "))
	}

	img, err := render(cfg, placed)
	if err != nil {
		return err
	}
	if err := ilibgo.SaveImageFile(*out, img, ilibgo.FormatPNG); err != nil {
		return fmt.Errorf("save %s: %w", *out, err)
	}
	fmt.Printf("wrote %s (%d words placed)\n", *out, len(placed))
	return nil
}

// flagValues bundles the flag pointers so applyFlags can override config fields
// only for flags the user actually set.
type flagValues struct {
	out, font, bg, palette *string
	width, height          *int
	minSize, maxSize, dpi  *float64
	seed                   *int64
}

// applyFlags overlays explicitly-set flags onto cfg (flags win over JSON).
func applyFlags(cfg *Config, fs *flag.FlagSet, v flagValues) {
	set := map[string]bool{}
	fs.Visit(func(f *flag.Flag) { set[f.Name] = true })

	if set["font"] {
		cfg.Font = *v.font
	}
	if set["w"] {
		cfg.Width = *v.width
	}
	if set["h"] {
		cfg.Height = *v.height
	}
	if set["bg"] {
		cfg.Background = *v.bg
	}
	if set["palette"] {
		cfg.Palette = splitPalette(*v.palette)
	}
	if set["min"] {
		cfg.MinSize = *v.minSize
	}
	if set["max"] {
		cfg.MaxSize = *v.maxSize
	}
	if set["dpi"] {
		cfg.DPI = *v.dpi
	}
	if set["seed"] {
		cfg.Seed = *v.seed
	}
}

// splitPalette parses a comma-separated color list, trimming blanks.
func splitPalette(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// gatherWords merges words from the config with words parsed from the frequency
// input (file arguments if present, otherwise stdin when piped).
func gatherWords(fromConfig []WordSpec, fileArgs []string, stdin io.Reader) ([]WordSpec, error) {
	specs := append([]WordSpec(nil), fromConfig...)

	if len(fileArgs) > 0 {
		for _, path := range fileArgs {
			f, err := os.Open(path)
			if err != nil {
				return nil, fmt.Errorf("open %s: %w", path, err)
			}
			parsed, perr := parseFrequency(f)
			f.Close()
			if perr != nil {
				return nil, perr
			}
			specs = append(specs, parsed...)
		}
		return specs, nil
	}

	if hasPipedInput(stdin) {
		parsed, err := parseFrequency(stdin)
		if err != nil {
			return nil, err
		}
		specs = append(specs, parsed...)
	}
	return specs, nil
}

// hasPipedInput reports whether stdin is a pipe/redirect with data, so an
// interactive run without input does not block on the terminal.
func hasPipedInput(stdin io.Reader) bool {
	f, ok := stdin.(*os.File)
	if !ok {
		return true // non-*File readers in tests always carry data
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) == 0
}

// newFontLoader returns a size-keyed font loader. When path is empty it uses the
// built-in goregular TTF; otherwise it loads the file once and re-rasterizes it
// per size. Results are cached by rounded point size.
func newFontLoader(path string, dpi float64) fontLoader {
	cache := map[int]*ilibgo.Font{}
	var fileData []byte
	return func(size float64) (*ilibgo.Font, error) {
		key := int(size + 0.5)
		if f, ok := cache[key]; ok {
			return f, nil
		}
		var (
			f   *ilibgo.Font
			err error
		)
		if path == "" {
			f, err = ilibgo.LoadTrueTypeFromBytes(goregular.TTF, "goregular", size, dpi)
		} else {
			if fileData == nil {
				fileData, err = os.ReadFile(path)
				if err != nil {
					return nil, fmt.Errorf("read font %s: %w", path, err)
				}
			}
			f, err = ilibgo.LoadTrueTypeFromBytes(fileData, path, size, dpi)
		}
		if err != nil {
			return nil, fmt.Errorf("load font at %gpt: %w", size, err)
		}
		cache[key] = f
		return f, nil
	}
}
