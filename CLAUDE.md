# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`ilibgo` is a Go image library (module `github.com/craigk5n/ilibgo`, package `ilibgo`) for reading, creating, manipulating, and saving images, with a focus on drawing primitives (lines, arcs, circles, rectangles, polygons, flood fill) and rendering text using X11 BDF bitmap fonts. It was ported from the original [C library `ilib`](https://github.com/craigk5n/ilib); the API stays close to the C version, which was itself modeled after X11 graphics functions.

## Commands

```bash
go build ./...          # build the library and all tools
go vet ./...            # static analysis
go test ./...           # run tests (note: no test files exist yet)
```

Build/run the example and tools (each is its own `package main`):

```bash
go run ./sample         # generates out.png demonstrating the API
go run ./displayfont -infile fonts/.../X.bdf -outfile out.png
go run ./thumbnails -w 120 -h 120 out.png in1.jpg in2.jpg
go run ./webreport [-tod|-dom|-dow|-moy] access.log
go run ./bdftogo -infile X.bdf -package mypkg   # convert BDF font -> embeddable .go
```

## Architecture

The root directory is the `ilibgo` library package. Each source file holds one concern:

- **`ilib.go`** — core types and file I/O. Defines `Image`, `Color`, `Font`, `GraphicsContext`, `Point`, `ImageOption`; `CreateImage`/`CreateImageWithBackground`; `ReadImageFile`/`WriteImageFile`; format detection (`FormatStringToType`, `FileType`, `IsSupportedFormat`). Reading uses Go's `image.Decode` (any registered format); writing supports PNG, JPEG, GIF, BMP, TIFF, PPM (PGM/PBM/XPM are stubbed and return errors).
- **`const.go`** — enums/constants: `ImageFormat`, `LineStyle`, `FillStyle`, `TextStyle`, `TextDirection`, and version info.
- **`igc.go`** — `GraphicsContext` (GC) construction and setters (`CreateGraphicsContext`, `SetForeground`, `SetBackground`, `SetFont`, `SetLineWidth`, `SetLineStyle`, `SetTextStyle`).
- **`color.go`** — `AllocColor(r,g,b)` and `AllocNamedColor(name)` backed by an embedded X11 `rgb.txt` color map.
- **Drawing primitives** — `point.go`, `line.go`, `arc.go` (arcs/ellipses), `circle.go`, `rectangle.go`, `polygon.go`, `flood.go`, `copy.go` (`CopyImage`, `CopyImageScaled`).
- **`text.go`** — text rendering: `DrawString`, `DrawStringRotated` (left-to-right / top-to-bottom / bottom-to-top), `DrawStringRotatedAngle`, and `TextDimensions`/`TextWidth`/`TextHeight`. Implements the `TextStyle` effects (etched/shadowed).
- **`fontbdf.go`** — BDF font parser. `LoadFontFromFile(path, name)` reads a `.bdf` file; `LoadFontFromBytes(name, data)` parses raw bytes (used by the embedded fonts); `LoadFontFromData(name, lines)` parses already-split lines (the shared core). ASCII chars are indexed into a `[256]BdfChar` array; non-ASCII go into `otherChars`.

### Key API conventions

- **Method-style drawing API.** Drawing operations are methods on `*Image`, e.g. `img.FillRectangle(gc, x, y, w, h)`. The `GraphicsContext` (passed by value) carries foreground/background color, font, line/text style — set it up, then pass it to the draw methods. The old free-function forms (`FillRectangle(img, gc, …)`, including `Copy*` taking `(source, dest, …)`) still exist as `// Deprecated:` forwarders in `deprecated.go` for backward compatibility; new code should use the methods. Note `Copy*` are methods on the **destination** image: `dst.CopyImage(src, gc, …)`. Text *measurement* (`TextWidth`/`TextHeight`/`TextDimensions`) and GC setters (`SetForeground`, etc.) remain free functions since they don't act on an `*Image`.
- **No more "I" prefix.** The port dropped the leading `I` from the C API (`IAllocColor` → `AllocColor`). The arc/circle functions that still had it are now `DrawArc`, `DrawEnclosedArc`, `FillArc`, `DrawCircle`, `FillCircle`, `DrawEllipse`. The old `IXxx` names remain as `// Deprecated:` aliases in `deprecated.go` (forwarding to the new names) for backward compatibility — don't use them in new code.
- **Mostly opaque types.** `Image`, `Font`, and `GraphicsContext` have private fields; manipulate them only through the provided functions (e.g. `ImageWidth(img)`, not `img.width`). Exceptions callers construct directly: `Point` has exported `X, Y` fields (plus the `Pt(x, y)` helper), and `Color` is built with `NewColor(r, g, b, a)` / `AllocColor` / `AllocNamedColor` and read back via its `RGBA()` method (it satisfies `image/color.Color`).
- **Standard-library interop.** `*Image` implements `image.Image` and `draw.Image` (`image_std.go`), so it can be passed straight to `image/draw`, the stdlib encoders, and `golang.org/x/image` scalers alongside the package's own drawing functions.

### Fonts

`fonts/` holds the bundled BDF fonts as canonical `.bdf` files, grouped by foundry (`adobe_100dpi`, `adobe_utopia_100dpi`, `bh_lucidatypewriter_100dpi`). Each foundry is its own Go package: an `embed.go` embeds that directory's `*.bdf` via `//go:embed` and exposes a private `lines()` helper, and one small wrapper file per font keeps the existing accessor API — `font.Font_helvR24()` still returns `[]string`, now read from the embedded file. Pass that to `LoadFontFromData`, or load any `.bdf` at runtime with `LoadFontFromFile` / from bytes with `LoadFontFromBytes`. The `bdftogo` tool remains available for users who prefer to bake their own fonts into Go `[]string` source, but it is no longer part of this repo's build.

## Notes

- There are currently **no `_test.go` files** in this repo.
- Several features are stubbed/TODO: PGM/PBM/XPM writing, line widths > 3, dashed line styles, tiled/stippled fills, escape-sequence (non-ASCII named) character handling in text rendering.
- `.gitignore` excludes generated output (`out.png`, `test.png`, `*.gif`, `*.ppm`).
