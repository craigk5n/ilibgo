package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// WordSpec is a single word in the cloud. A word's rendered point size comes
// from Size when it is > 0; otherwise it is derived from Count via sqrt scaling
// (see buildWords). Color, when set, overrides the random palette pick.
type WordSpec struct {
	Text  string  `json:"text"`
	Size  float64 `json:"size,omitempty"`
	Count int     `json:"count,omitempty"`
	Color string  `json:"color,omitempty"`
}

// Config holds every setting for a render. It is populated from defaults, then
// (optionally) a JSON file, then explicitly-set command-line flags, in that
// order of increasing precedence.
type Config struct {
	Width      int        `json:"width,omitempty"`
	Height     int        `json:"height,omitempty"`
	Background string     `json:"background,omitempty"`
	Font       string     `json:"font,omitempty"` // path to a .ttf/.otf; empty = built-in goregular
	DPI        float64    `json:"dpi,omitempty"`
	MinSize    float64    `json:"minSize,omitempty"`
	MaxSize    float64    `json:"maxSize,omitempty"`
	Palette    []string   `json:"palette,omitempty"`
	Seed       int64      `json:"seed,omitempty"`
	Words      []WordSpec `json:"words,omitempty"`
}

// defaultConfig returns a Config with sensible built-in values. The palette is a
// readable set on a dark background.
func defaultConfig() Config {
	return Config{
		Width:      800,
		Height:     600,
		Background: "black",
		Font:       "",
		DPI:        72,
		MinSize:    14,
		MaxSize:    96,
		Palette:    []string{"#e6194b", "#3cb44b", "#ffe119", "#4363d8", "#f58231", "#46f0f0", "#f032e6", "#fabebe"},
		Seed:       1,
		Words:      nil,
	}
}

// loadConfigFile reads a JSON config, overlaying it onto cfg. Unspecified JSON
// fields leave cfg's values untouched (json.Unmarshal only writes present keys).
func loadConfigFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config %s: %w", path, err)
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config %s: %w", path, err)
	}
	return nil
}

// validate checks the config for usable values before rendering.
func (c Config) validate() error {
	if c.Width <= 0 || c.Height <= 0 {
		return fmt.Errorf("image dimensions must be positive (got %dx%d)", c.Width, c.Height)
	}
	if c.MinSize <= 0 || c.MaxSize <= 0 {
		return fmt.Errorf("font sizes must be positive (min=%g max=%g)", c.MinSize, c.MaxSize)
	}
	if c.MinSize > c.MaxSize {
		return fmt.Errorf("minSize (%g) must not exceed maxSize (%g)", c.MinSize, c.MaxSize)
	}
	if c.DPI <= 0 {
		return fmt.Errorf("dpi must be positive (got %g)", c.DPI)
	}
	if len(c.Palette) == 0 {
		return fmt.Errorf("palette must contain at least one color")
	}
	return nil
}
