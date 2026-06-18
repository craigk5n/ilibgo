# IDEAS.md — Improvement Recommendations for `ilibgo`

Recommendations for bringing `ilibgo` in line with Go best practices for an image library. Grounded in a read of the current source (ported from the original C `ilib`). Ordered roughly by impact. Each item notes concrete locations.

---

## 1. Correctness bugs (fix first)

These are real defects found while reading the port. Most stem from C → Go translation.

### 1.1 `WriteImageFile` silently discards encoder errors
`ilib.go:152` — only the PGM/PBM/XPM branches assign `err`; the real encoders ignore their return value:
```go
png.Encode(f, img.data)        // error dropped
jpeg.Encode(f, img.data, &options)  // error dropped
// ...gif, bmp, tiff, ppm all drop the error
```
A failed write (disk full, broken pipe) returns `nil`. Capture and return each encoder's error, and join it with the `f.Close()` error.

### 1.2 `WriteImageFile` closes a file the doc says it leaves open
`ilib.go:151-154`: the comment says *"The file is left open for the caller to close,"* but the body does `defer f.Close()`. `ReadImageFile` (`ilib.go:186`) does the opposite — it never closes. Pick one contract. **Better: don't take `*os.File` at all** (see §3.4) — take `io.Writer`/`io.Reader` and let the caller own the file. A function should not close a handle it didn't open.

### 1.3 `IDrawEnclosedArc` truncates trig to integers — broken geometry
`arc.go:54-55`:
```go
myx = x + int(r1*int(math.Cos(a+float64(loop)*da)))   // int(cos) ∈ {-1,0,1}
myy = y + int(r2*int(math.Sin(a+float64(loop)*da)))
```
`int(math.Cos(...))` collapses to -1/0/1 *before* multiplying by the radius. Compare the correct form in `IDrawArc` (`arc.go:30`): `int(float64(r1)*math.Cos(...))`. Fix to multiply in float64, then cast once.

### 1.4 `DrawStringRotatedAngle` indexes empty substrings and walks bytes
`text.go:279,284,285,298`: `char := text[loop:loop]` is always `""` (`loop:loop`), so no character is ever matched. Should be `text[loop:loop+1]`, and the loop should range over `[]rune` (like `drawStringRotated90` does) to be UTF-8-safe. As written this function renders nothing.

### 1.5 BDF parser: `SLANT` overwrites the wrong field
`fontbdf.go:176-177`: the `SLANT` branch assigns `font.weight` instead of `font.slant`. The `slant` field is declared but never populated.

### 1.6 BDF parser: `DWIDTH` `Sscanf` arg/verb mismatch
`fontbdf.go:181-182`: format `"DWIDTH %d %d"` has two verbs but four target pointers, and the second `%d` writes into `char.height` — clobbering the height already parsed from `BBX`. Verify against the BDF spec (`DWIDTH dwx0 dwy0`) and parse only what the format declares.

### 1.7 BDF parser: duplicated `PIXEL_SIZE` branch (dead code)
`fontbdf.go:144-155`: the `PIXEL_SIZE` `else if` appears twice back-to-back; the second is unreachable.

### 1.8 `FillPolygon` returns a leftover debug error
`polygon.go:158`: `return fmt.Errorf("ilib bug no. 34534345")`. Either handle the single-intersection case properly or return a descriptive error — don't ship a debug sentinel.

### 1.9 `FloodFill` is unbounded recursion
`flood.go:54,58`: recurses per scanline-neighbor. A large fill region overflows the goroutine stack. Reimplement with an explicit stack/queue (scanline flood fill). Also the bottom-edge check `y < image.height` then reads `y+1` — off by one at the last row (harmless only because `image.RGBA.RGBAAt` clips out-of-bounds to zero, which is itself a latent logic bug masking).

---

## 2. Make the public API actually usable and idiomatic

### 2.1 `Point` and `Color` cannot be constructed by external callers — blocker
`Point{x, y}` (`ilib.go:33`) and `Color{color}` (`ilib.go:38`) have **only unexported fields**. External code can't build a `Point`, so `DrawPolygon`/`FillPolygon([]Point)` are unusable outside the package, and `Color` can only come from `AllocColor`/`AllocNamedColor`. Options:
- Export fields (`Point{X, Y int}`), or
- Provide constructors (`func Pt(x, y int) Point`), or
- Reuse the standard library: accept `image.Point` and `color.Color` directly.

**Strong recommendation: lean on `image`/`color` standard types** where possible rather than wrapping them — interoperates with the entire Go image ecosystem.

### 2.2 Prefer methods over free functions taking `*Image` — ✅ DONE
**Implemented.** All drawing operations are now methods on `*Image` (e.g. `img.FillRectangle(gc, …)`; `Copy*` are methods on the destination image). The old free functions remain as `// Deprecated:` forwarders in `deprecated.go`. The `GraphicsContext` stayed a by-value parameter rather than the pointer sketched below. Original note below.

Every primitive is `func DrawX(image *Image, gc GraphicsContext, ...)`. Idiomatic Go would hang these off a type:
```go
func (img *Image) FillRectangle(gc *GraphicsContext, r image.Rectangle) error
```
This reads better (`img.FillRectangle(...)`), enables a fluent style, and lets `Image` implement `image.Image`/`draw.Image` (see §4.1). The current shape is a literal C port.

### 2.3 Remove the inconsistent `I` prefix — ✅ DONE
**Implemented.** Renamed to `DrawArc`, `DrawEnclosedArc`, `FillArc`, `DrawCircle`, `FillCircle`, `DrawEllipse`; the old `IXxx` names are kept as `// Deprecated:` forwarding aliases in `deprecated.go`. Original note below.

Most functions dropped the C `I` prefix, but these kept it: `IDrawArc`, `IDrawEnclosedArc`, `IFillArc`, `IDrawCircle`, `IFillCircle`, `IDrawEllipse` (`arc.go`, `circle.go`). Rename for consistency (`DrawArc`, `FillArc`, …). Keep deprecated aliases for one release if backward-compat matters.

### 2.4 Don't return `error` you never produce
`DrawLine`, `DrawRectangle`, `FillRectangle`, `DrawPolygon`, `SetPoint`, `DrawPoint`, the arc/circle funcs all `return nil` unconditionally. Either:
- drop the return (cleaner), or
- actually validate inputs (negative radius, nil font, out-of-range coords) and return meaningful errors.
Returning a permanently-nil `error` trains callers to ignore it — which then hides the cases that *do* fail.

### 2.5 Decide `GraphicsContext` value-vs-pointer semantics
`CreateGraphicsContext()` returns a value (`igc.go:6`), but setters take `*GraphicsContext`, and draw functions take it by value. Passing by value is fine (cheap, naturally concurrency-safe) but then `SetForeground(&gc, …)` mutating a local is easy to get wrong. Pick one model and document it. Functional options or a builder would be cleaner:
```go
gc := NewGC(WithForeground(red), WithFont(f))
```

### 2.6 Polish error strings
`color.go` `AllocNamedColor`: `"no such color'%s'"` is missing a space. Go convention: lowercase, no trailing punctuation, wrap with `%w` when propagating.

---

## 3. Tooling, modules, and project hygiene

### 3.1 Bump the Go version
`go.mod` says `go 1.13` (2019). Move to a currently-supported release (1.22/1.23+). Benefits: built-in `min`/`max` (lets you delete the hand-rolled ones in `polygon.go:16-30`), `//go:embed`, generics, `errors.Join`, modern `golang.org/x/image`.

### 3.2 Add tests — there are currently none
No `*_test.go` files exist. High-value, low-effort targets:
- **Table-driven unit tests**: `FormatStringToType`/`FileType`/`IsSupportedFormat`, `AllocNamedColor`, `colorsMatch`, `TextDimensions`.
- **Golden-image tests** for drawing primitives: render to an `*image.RGBA`, compare against a committed PNG (or its SHA-256). The `sample` tool is effectively a manual golden test already — formalize it.
- **BDF round-trip / parse tests** using a small fixture font.
- **Fuzz test** the BDF parser (`LoadFontFromData`) — it does a lot of hand `strconv`/`Sscanf` parsing on untrusted text (`go test -fuzz`).
- **Benchmarks** for `FillRectangle`, `CopyImageScaled`, `FloodFill`.
Target the repo's 80% coverage bar (`go test -cover -race ./...`).

### 3.3 Add CI + linters
No CI config present. Add a GitHub Actions workflow running `go vet`, `go test -race -cover`, `staticcheck`, `golangci-lint`, and `gosec`. `go vet` alone would have flagged the `Sscanf` arg mismatch in §1.6 and likely the empty-slice in §1.4.

### 3.4 Use `io.Reader`/`io.Writer`, not `*os.File`
`ReadImageFile(f *os.File)` / `WriteImageFile(f *os.File, …)` force a file. Accepting `io.Reader`/`io.Writer` makes them testable in-memory (`bytes.Buffer`), usable with HTTP responses, gzip streams, etc. Keep thin `…File(path string)` convenience wrappers that open/close.

### 3.5 snake_case parameter names
`copy.go`: `src_x`, `src_y`, `dest_x`, `dest_width`… Go uses `srcX`, `destWidth`. Minor but it's the first thing reviewers notice.

### 3.6 Add package docs and runnable examples
There's a good package comment in `ilib.go`, but no `Example_*` functions (which appear in godoc and run under `go test`). Add `ExampleCreateImage`, `ExampleDrawString`, etc. Consider a `doc.go`.

### 3.7 Dead/placeholder fields
`Image.comments` (`ilib.go:59`) is parsed nowhere and written nowhere (TODO). Either implement metadata pass-through to PNG/JPEG or remove it.

---

## 4. Leverage the standard library (`image`, `draw`)

### 4.1 Implement the standard interfaces
`Image` already wraps `*image.RGBA` (`ilib.go:60`). Have it satisfy `image.Image` (`ColorModel`, `Bounds`, `At`) and `draw.Image` (`Set`). Then it drops straight into `image/draw`, `image/png`, scaling packages, etc., and users can mix `ilibgo` with the broader ecosystem for free.

### 4.2 Use `draw.Draw` for fills and copies (also a perf win — see §5)
`CreateImageWithBackground` (`ilib.go:137`) fills via a per-pixel `SetPoint` double loop. Replace with one call:
```go
draw.Draw(img.data, img.data.Bounds(), &image.Uniform{C: background.color}, image.Point{}, draw.Src)
```
Same for `FillRectangle` (uniform over a sub-rect) and the unscaled path of `CopyImage` (`draw.Draw(dest, r, src, sp, draw.Src)`).

---

## 5. Performance

### 5.1 Avoid per-pixel `Set`/`At` in hot loops
`image.RGBA.Set`/`RGBAAt` do bounds checks and color conversion on every call. `CreateImageWithBackground`, `FillRectangle`, `CopyImage`, `FloodFill` all loop pixel-by-pixel through these. For fills, prefer `draw.Draw` (§4.2). For custom inner loops, index `img.data.Pix` directly with the stride:
```go
i := img.data.PixOffset(x, y)
img.data.Pix[i+0], img.data.Pix[i+1], img.data.Pix[i+2], img.data.Pix[i+3] = r, g, b, a
```

### 5.2 `CopyImage` reassigns the GC foreground every pixel
`copy.go:28-31` calls `GetPoint`+`SetForeground`+`SetPoint` per pixel. For a straight copy this is `draw.Draw` (one call). Keep the manual loop only for the transparent-color-skip feature, and even then write `Pix` directly.

### 5.3 `CopyImageScaled` uses nearest-neighbor only
`copy.go:44` acknowledges this. For quality scaling, delegate to `golang.org/x/image/draw` (`draw.CatmullRom`, `draw.ApproxBiLinear`) — already an indirect dependency via `x/image`.

### 5.4 Iterative flood fill
Covered in §1.9 — also a performance issue: recursion + per-pixel `GetPoint` is slow and stack-heavy. Scanline flood fill with a `Pix`-backed visited check is far faster.

### 5.5 Precompute trig in arc loops
`arc.go` recomputes `math.Cos/Sin` and the `2*Pi/360` conversion inside the loop. Hoist the constant and, if needed, step with incremental rotation. Minor compared to the above.

---

## 6. Fonts & packaging

### 6.1 Replace generated `.go` font files with `//go:embed` — ✅ DONE
**Implemented.** The bundled fonts are now canonical `.bdf` files embedded per-foundry via `//go:embed *.bdf`; the `Font_xxx() []string` accessors are preserved as thin wrappers (non-breaking), and `LoadFontFromBytes` was added. `bdftogo` is retained as an optional tool but is no longer part of the build. Original notes below for reference.

`bdftogo` converts each `.bdf` into a Go file holding a giant `[]string`. With Go ≥1.16 you can instead embed the raw `.bdf` bytes:
```go
//go:embed fonts/adobe_100dpi/helvR24.bdf
var helvR24BDF []byte
```
and parse at init/first-use. This removes the codegen step entirely, shrinks the source tree, compiles faster, and keeps fonts in their canonical format. `bdftogo` can stay as an optional tool for users who want pre-baked fonts.

### 6.2 Consider splitting bundled fonts into a sub-package or separate module
The embedded font tables bloat every binary that imports the root package even if it never draws text. A `ilibgo/fonts` sub-module (or build-tagged inclusion) lets callers opt in.

### 6.3 Parse BDF more robustly
The parser uses fixed substring offsets (`line[10:]`, `line[12:]`, etc.) and assumes field positions. Switch to `strings.Fields`/`strings.Cut` keyed on the keyword — resilient to extra whitespace and safer against panics on short lines. Pair with the fuzz test in §3.2.

---

## 7. Concurrency & safety notes (document, don't necessarily change)

- `GraphicsContext` passed by value is good — it's effectively immutable per call.
- An `*Image` is **not** safe for concurrent writes (shared `*image.RGBA`). Document this. If parallel tile rendering is desired later, render into separate sub-images and composite with `draw.Draw`.
- No use of `context.Context` anywhere; for a pure CPU image library that's fine — don't add it speculatively.

---

## 8. Example apps & utilities (ideas)

Small programs that showcase the library and/or are useful on their own. Each is a `package main` under its own directory, like the existing tools.

**Done:**
- **`iconvert`** ✅ — format converter (`iconvert in.png out.tiff`), thin wrapper over `ReadImageFile`/`FileType`/`WriteImageFile`.
- **`chart`** ✅ — bar-chart generator from `label=value` args or `label,value` stdin; a general successor to the hardcoded `webreport` grapher.
- **`qr` package + `qrgen` tool** ✅ — dependency-free QR encoder (byte mode, all versions/ECC levels, Reed-Solomon, mask selection) plus a renderer that draws it with `FillRectangle`. See §10.
- **`captcha`** ✅ — distorted-text CAPTCHA: per-char `DrawStringRotatedAngle` with random angle/jitter/color, noise lines (`DrawLine`) and speckle (`SetPoint`); prints the answer and is seedable for reproducibility.

**Examples (showcase value):**
- **`sparkline`** — tiny inline trend image from a number list; minimal "hello world" for the draw API.
- **`watermark`** — overlay semi-transparent text/logo onto a loaded image; demonstrates the alpha `NewColor` + image round-trip.
- **`fontsheet`** — render a full glyph specimen sheet for a bundled font by iterating the embedded `fonts.FS`; doubles as visual font QA.
- **`mandelbrot`/fractal** — per-pixel `SetPoint`/`NewColor`; also a natural benchmark target.
- **`montage`** — compose N images into a labeled grid (a configurable generalization of `thumbnails`).
- **`iresize`** — scale an image via `CopyImageScaled` (good demo target for §5.3 quality scaling).
- **`barcode`** — Code 39 / simple grid barcode with `FillRectangle`; pure geometry, no font.

**Utilities / infra:**
- **`bdfinfo`** — print a `.bdf`'s foundry/family/slant/ascent/descent/glyph-count; gives the parser a non-rendering consumer and exercises the metadata fields.
- **Golden-image CI check** — render a known scene, diff against a committed PNG hash to catch visual regressions.
- **Benchmarks** (`go test -bench`) — fills, `CopyImageScaled`, flood fill, fractal; anchor the perf claims from the `draw.Draw` work.

## 9. Scalable / TrueType font support

Feasible, medium effort for basic support; high effort for full feature parity. Notes:

- **Don't write a rasterizer** — wrap the stdlib/ecosystem: `golang.org/x/image/font`, `.../font/opentype`, `.../font/sfnt` (and `golang.org/x/image/math/fixed`). These parse and rasterize TrueType/OpenType.
- **The type is already designed for it** — `Font` wraps `*BdfFont` with a `// Add additional support font types (truetype, etc.) here` comment. Add a parallel TrueType-backed face.
- **Easy win, now unlocked:** because `*Image` implements `draw.Image` (§4.1), a `font.Drawer{Dst: img.data, Face: face, ...}.DrawString(s)` can render anti-aliased TrueType text straight onto an image in a few lines. Good for a basic `DrawTextTTF`-style method.
- **The expensive part is parity** with the existing BDF pipeline: the current renderer is bitmap-specific (per-pixel `SetPoint` from `bdfChar.data` booleans) and supports `TextStyle` (etched/shadowed) and arbitrary-angle rotation. TrueType glyphs are anti-aliased coverage masks that must be **alpha-blended**, not set/unset — and rotation needs transformed rasterization. Unifying TrueType with etched/shadowed/rotated effects is the real work.
- **Suggested path:** introduce a small `face` abstraction (glyph mask + advance), add a TrueType implementation, expose a basic anti-aliased `DrawString`-style method first, and treat style/rotation parity as follow-ups.

## 10. QR codes

- **Generator — ✅ DONE (from scratch, no deps).** Implemented as the `qr` package (byte mode, automatic smallest-version selection for all 40 versions, the four ECC levels, Reed–Solomon over GF(256), and lowest-penalty mask selection) with a `qrgen` tool that renders the module matrix via `FillRectangle`. Output was verified to scan with an independent ZXing-port decoder across versions 1–17 and all ECC levels; a golden-fingerprint test locks the layout against regressions. (Original note: a mature encoder like `rsc.io/qr` or `skip2/go-qrcode` would also work, but the from-scratch route keeps the project dependency-free.)
- **Decoder — possible but largely out of scope.** Decoding is a computer-vision task (binarization, finder-pattern detection, perspective correction, grid sampling, then RS error correction). Mature ports exist (`github.com/makiuchi-d/gozxing`, `github.com/liyue201/goqr`). `ilibgo`'s role would only be loading the image (`ReadImageFile` → `image.Image`); the analysis belongs in a dedicated CV library. Reasonable as a `qrdecode` example that wraps `gozxing`, but not a core library feature — `ilibgo` is about image *creation/manipulation*, not *analysis*.

---

## Suggested order of work

1. **§1 correctness bugs** (especially 1.1–1.6) — these are shipping defects.
2. **§3.1–3.3** — bump Go, add tests + CI so regressions are caught.
3. **§2.1** — make `Point`/`Color` constructible (unblocks real external use).
4. **§4 + §5** — stdlib interfaces and `draw.Draw` (correctness + perf together).
5. **§2.2–2.6** API polish, then **§6** font modernization.
