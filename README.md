# Ilib - Image Library for Go

Copyright (C) 2001-2022 Craig Knudsen, craig@k5n.us
http://github.com/craigk5n/ilibgo

Ilib is a library (and some tools and examples) written in Go
that can read, create, manipulate and save images.  It is capable
of using X11 BDF fonts for drawing text.  That means you get
lots of fonts to use.  You can even create your
own if you know how to create an X11 BDF font.  It should be able
to read any image file that the base Go image package supports.


## History

This could was ported to Go from the original [version in C](https://github.com/craigk5n/ilib).
The API remains as close to the original as it made sense to do.
All the functions dropped the "I" at the start since Go supports
namespaces.  The original C library was modeled after the
X11 functions, documented
[here](https://www.x.org/releases/X11R7.6/doc/man/man3/) as man pages.


## BDF Fonts

Additional BDF Fonts can be found at:
  [https://gitlab.freedesktop.org/xorg/font](https://gitlab.freedesktop.org/xorg/font)

Note that fonts can be loaded as external fonts at run-time, or
fonts can be embedded in the binary by using the [`bdftogo`](clients/bdftogo)
tool included in this package.


Note that some of the BDF fonts are bundled with this package.  Please
see the [COPYING](https://gitlab.freedesktop.org/xorg/font/adobe-100dpi/-/blob/master/COPYING)
notice

> Copyright 1984-1989, 1994 Adobe Systems Incorporated.
> Copyright 1988, 1994 Digital Equipment Corporation.
> 
> Adobe is a trademark of Adobe Systems Incorporated which may be
> registered in certain jurisdictions.
> Permission to use these trademarks is hereby granted only in
> association with the images described in this file.
> 
> Permission to use, copy, modify, distribute and sell this software
> and its documentation for any purpose and without fee is hereby
> granted, provided that the above copyright notices appear in all
> copies and that both those copyright notices and this permission
> notice appear in supporting documentation, and that the names of
> Adobe Systems and Digital Equipment Corporation not be used in
> advertising or publicity pertaining to distribution of the software
> without specific, written prior permission.  Adobe Systems and
> Digital Equipment Corporation make no representations about the
> suitability of this software for any purpose.  It is provided "as
> is" without express or implied warranty.


Edit/Import/Create BDF fonts with
[FontForge](https://github.com/fontforge/fontforge).

## Building
`go build`


## History
See [ChangeLog](ChangeLog.md)
