# qr — a dependency-free QR Code encoder

Package `qr` encodes data as a [QR Code](https://en.wikipedia.org/wiki/QR_code)
(ISO/IEC 18004). It has **no third-party dependencies** — only the Go standard
library.

```
import "github.com/craigk5n/ilibgo/qr"
```

## Features

- **Byte mode** — encodes any input (ASCII, UTF-8, arbitrary bytes).
- **Automatic version selection** — picks the smallest of the 40 versions that fits.
- **All four error-correction levels** — `Low`, `Medium`, `Quartile`, `High`.
- **Reed–Solomon** error correction over GF(256), with block interleaving.
- **Automatic mask selection** — evaluates all 8 masks and keeps the lowest-penalty one.

The output is a square grid of dark/light modules; rendering to an image (or
anything else) is up to the caller.

## Library usage

```go
code, err := qr.EncodeText("https://github.com/craigk5n/ilibgo", qr.Medium)
if err != nil {
    log.Fatal(err)
}
fmt.Println(code.Version(), code.Size()) // e.g. 3 29

// Walk the module grid (true = dark).
for y := 0; y < code.Size(); y++ {
    for x := 0; x < code.Size(); x++ {
        if code.Module(x, y) {
            // ... draw a dark module at (x, y) ...
        }
    }
}
```

`Encode([]byte, Ecc)` is the same but takes raw bytes.

### Render to an `ilibgo` image

This is exactly what the `qrgen` tool does — draw each dark module as a filled
square, leaving a light "quiet zone" border:

```go
const scale, border = 8, 4 // pixels per module, quiet-zone modules
dim := (code.Size() + 2*border) * scale

white, _ := ilibgo.AllocNamedColor("white")
img := ilibgo.CreateImageWithBackground(dim, dim, white)

gc := ilibgo.CreateGraphicsContext()
black, _ := ilibgo.AllocNamedColor("black")
ilibgo.SetForeground(&gc, black)

for y := 0; y < code.Size(); y++ {
    for x := 0; x < code.Size(); x++ {
        if code.Module(x, y) {
            img.FillRectangle(gc, (border+x)*scale, (border+y)*scale, scale, scale)
        }
    }
}
// img.WriteImageFile(...) etc.
```

### Render to the terminal

```go
for y := 0; y < code.Size(); y++ {
    for x := 0; x < code.Size(); x++ {
        if code.Module(x, y) {
            fmt.Print("##")
        } else {
            fmt.Print("  ")
        }
    }
    fmt.Println()
}
```

## Command-line tool: `qrgen`

The `qrgen` program (in `../qrgen`) renders a QR code straight to a PNG.

```
qrgen [options] text...

  -out string    output PNG file (default "qr.png")
  -ecc string    error-correction level: L, M, Q, or H (default "M")
  -scale int     pixels per module (default 8)
  -border int    quiet-zone width in modules (default 4)
```

Examples:

```sh
qrgen -out site.png https://github.com/craigk5n/ilibgo
echo "hello world" | qrgen -ecc H -scale 12 -out hello.png
```

The text is taken from the arguments, or from standard input when no arguments
are given.

## Notes & limitations

- **Byte mode only.** Numeric and alphanumeric modes (which pack digits/uppercase
  more tightly) are not implemented; byte mode encodes everything correctly, just
  slightly less compactly for those inputs.
- **No ECI / structured append / Kanji mode.**
- **Verification.** The encoder's output has been round-trip verified against an
  independent QR decoder across versions 1–17 and all error-correction levels; a
  golden-fingerprint test guards against regressions. That decoder is not a
  dependency of this project.

## License

Same license as the parent `ilibgo` project.
