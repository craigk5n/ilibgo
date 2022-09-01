This directory contains BDF fonts that can easily be used with libgo.
The original source of the BDF font files is:
  (https://gitlab.freedesktop.org/xorg/font)[https://gitlab.freedesktop.org/xorg/font]

To add more fonts, you will need to do the following:
- Create a new directory in this directory.  Note: Do not use dashes since
  they cannot be used in Golang package subdirectories.  Replace and dashes
  with an underscore.
- Copy the `Makefile` from one of the other examples like `adobe_100dpi`
  into your new directory.
- Edit your new `Makefile` and be sure to change the make variables at the top:
  - GITLABURL: This should be the base URL where your bdf fonts files can
    be downloaed
  - PACKAGE_NAME: This should the same as your new directory name
  - BDF_FILES: This should be a list of each BDF font filename that can be
    found in the GITLABURL repository.
- Run the following command to update the `Makefile`:
  ```
  make rebuild
  ```
- Then, run `make` (with no arguments).  This will download the BDF files locally,
  then create a Golang version of the font along with a PNG representation of
  the first 256 ASCII characters of the font.
- Create a pull request if you want to contribute the fonts back to the ilib project.

