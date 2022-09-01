# Change Log

### 31 Aug 2022
  - Release 2.0.0
    + Converted from C to Go.  This includes almost all of the original C code.
      except the ifraggraph tool which is not really useful in 2022.  The Perl
      module was also dropped.  Most of the API dropped the leading "I", so
      "IAllocColor" is now "AllocColor" because Go supports namespaces better than
      K&R C.
### 25 Oct 2004
  - Release 1.1.9
    + Added IFloodFill
### 15 Aug 2001
  - Release 1.1.8
    + The perl module is now included with the main distribution.
    + Added missing Makefile from fonts directory.
    + Fixed bug in ICopyImageScaled().
    + Don't allow libjpeg to call exit() on errors reading input data,
      return an error number instead.
    + Fixed compile problems that may be encountered if not using either
      libjpeg, giflib or libpng.
### 24 May 2000
  - Release 1.1.7
    + Added (finally) "install" target.  Just do "make install" as root
      after "make" to install everything
    + Now builds both static and shared library
    + Added support for reading BMP images
      (contributed by Jim Winstead <jimw@trainedmonkey.com>)
    + Added IDrawStringRotatedAngle() (but does not yet support text
      style from ISetTextStyle yet)
      (contributed by Geovan Rodriguez <geovan@cigb.edu.cu>)
    + PNG output improvements
### 29 Nov 1999
  - Release 1.1.6
    + Added new functions:
        IAllocNamedColor()
        IDrawArc()
        IDrawEllipse()
        IDrawCircle()
        IDrawPolygon()
        IFillPolygon()
        IFillArc()
        IFillCircle()
        IFillEllipse()
        IDrawEnclosedArc()
        IArcProperties()
    + Fixed bug in iindex client
### 25 Aug 1999
  - Release 1.1.5
    + Added IDrawStringRotated() function for drawing text vertically
      (both at 90 and 270 degrees).
    + Added support for styled text using ISetTextStyle().  Current
      styles include ITEXT_ETCHED_IN, ITEXT_ETCHED_OUT and
      ITEXT_SHADOWED.  See the isample example application for
      an example.
    + Added ISetBackground() function (required to use the ITEXT_ETCHED_IN
      and ITEXT_ETCHED_OUT text drawings styles).
### 20 Aug 1999
  - Release 1.1.4
    + Added ISetComment() function
    + Added IGetTransparent() function
    + Fixed reading of interlaced GIF images (would die on an error)
    + Fixed reading of transparent GIF images (would die on an error)
    + Added support for writing interlaced GIF images
    + Added support for writing transparent GIF images
    + Updated names of examples and clients to start with 'i'
    + Updated iindex client to generated 3D-style imagemap with
      corresponding HTML for client-side imagemap.
### 23 Jul 1999
  - Release 1.1.3
    + Added support for reading/writing JPEG (color and greyscale)
      currently always uses 75 for quality when writing
    + Added support for reading/writing PGM (raw only)
    + Added new ICopyImageScaled() function.
    + Added new client application "index" that creates a single
      image index of mini-images from other images.
### 19 Jul 1999
  - Release 1.1.2
    + Added initial support for reading and writing PNG files.  Currently,
      all images are written as 16-bit images (no 8-bit colormapped
      images) and transparency, alpha channels and interlacing are
      not yet implemented when writing PNG files.
### 12 Apr 1999
  - Release 1.1.1
    + Fixed bug that would not allow high ascii values to be drawn.
    + Fixed bug where BDF fonts that specified a space as an empty
      bitmap were reported as corrupted font files.
    + Added new client application displayfont.
    + Fixed ICopyImage() to ignore transparent bits when copying
      an image.
### 18 May 1998
  - Release 1.1.0
    + Pulled out GIF code that was probably in violation of the
      Unisys copyright.  GIFLIB can be used instead (which is not
      in violation of the copyright because of the grandfather clause.)
    + Small fix in IDrawLine.
    + Added new IErrorString() function for turning IError values into
      strings suitable for error messages.
    + Added IFileType() function for guessing image type (PPM, GIF, etc.)
      by filename extension.
    + Added ITextHeight() and ITextDimensions() for calculating pixel
      dimensions of text.
### 20 May 1996
  - Release 1.0

