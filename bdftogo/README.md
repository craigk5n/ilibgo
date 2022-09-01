# bdftogo - Image Library for Go

This tool will create a Go file from a BDF font file so that the BDF font can be bundled
with the Go file.  This allows an application to distribute the font
within the binary rather than having to also include the BDF font file
separately and then load the BDF font file.

## Example usage

The following command creates a Go version of the Times Roman 8 point font
using "main" as the package name for it.
```
./bdftogo -infile ../fonts/adobe_100dpi/timR08.bdf -outfile /tmp/timR08.go -package main
```

## BDF Fonts

Additional BDF Fonts can be found at:
  [https://gitlab.freedesktop.org/xorg/font](https://gitlab.freedesktop.org/xorg/font)

Note that fonts can be loaded as external fonts at run-time, or
fonts can be embedded in the binary by using the [`bdftogo`](clients/bdftogo)
tool included in this package.

