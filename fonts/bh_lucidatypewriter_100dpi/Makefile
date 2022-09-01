# Makefile for fonts
# This will build go files for each bdf file
# IF you update the bdf flles in this directory, run "make rebuild" to
# update this makefile.

BDFTOGO = ../../bdftogo/bdftogo
CURL = curl --silent
DISPLAYFONT = ../../displayfont/displayfont
GITLABURL = https://gitlab.freedesktop.org/xorg/font/bh-lucidatypewriter-100dpi/-/raw/master
PACKAGE_NAME = bh_lucidatypewriter_100dpi

BDF_FILES = lutBS08.bdf lutBS10.bdf lutBS12.bdf lutBS14.bdf lutBS18.bdf lutBS19.bdf \
	lutBS24.bdf lutRS08.bdf lutRS10.bdf lutRS12.bdf lutRS14.bdf lutRS18.bdf lutRS19.bdf \
	lutRS24.bdf


_ALL: _GO_FILES

# Rebuild this Makefile based on the *.bdf files in this directory
rebuild:
	@echo Rebuilding makefile...
	@cat Makefile | sed -n '1,/^# BEGIN GENERATED CONTENT/p' > Makefile.new
	@echo "" >> Makefile.new
	@echo "" >> Makefile.new
	@echo "GO_FILES = " | tr -d '\012' >> Makefile.new
	@for f in $(BDF_FILES); do \
	  echo " \\" >> Makefile.new; \
	  echo "	$$f" | sed 's/.bdf/.go/' | tr -d '\012' >> Makefile.new; \
	done
	@echo "XXX" | tr 'X' '\012' >> Makefile.new
	@echo '_GO_FILES: $$(GO_FILES)' >> Makefile.new; echo "" >> Makefile.new
	@echo "XXX" | tr 'X' '\012' >> Makefile.new
	@for f in $(BDF_FILES); do \
	  echo "" >> Makefile.new; \
	  g=`echo $$f | sed 's/bdf/go/'`; \
	  p=`echo $$f | sed 's/bdf/png/'`; \
	  echo "$$g": >> Makefile.new; \
	  echo "	=(CURL) =(GITLABURL)/$$f > $$f" | tr = '\044' >> Makefile.new; \
	  echo "	=(BDFTOGO) -infile $$f -package =(PACKAGE_NAME) -outfile $$g" | tr = '\044' >> Makefile.new; \
	  echo "	=(DISPLAYFONT) -infile $$f -outfile $$p" | tr = '\044' >> Makefile.new; \
	  echo "" >> Makefile.new; \
	done
	@mv Makefile Makefile.old && mv Makefile.new Makefile

clean:
	rm -f *.go *.png

# BEGIN GENERATED CONTENT


GO_FILES =  \
	lutBS08.go \
	lutBS10.go \
	lutBS12.go \
	lutBS14.go \
	lutBS18.go \
	lutBS19.go \
	lutBS24.go \
	lutRS08.go \
	lutRS10.go \
	lutRS12.go \
	lutRS14.go \
	lutRS18.go \
	lutRS19.go \
	lutRS24.go



_GO_FILES: $(GO_FILES)






lutBS08.go:
	$(CURL) $(GITLABURL)/lutBS08.bdf > lutBS08.bdf
	$(BDFTOGO) -infile lutBS08.bdf -package $(PACKAGE_NAME) -outfile lutBS08.go
	$(DISPLAYFONT) -infile lutBS08.bdf -outfile lutBS08.png


lutBS10.go:
	$(CURL) $(GITLABURL)/lutBS10.bdf > lutBS10.bdf
	$(BDFTOGO) -infile lutBS10.bdf -package $(PACKAGE_NAME) -outfile lutBS10.go
	$(DISPLAYFONT) -infile lutBS10.bdf -outfile lutBS10.png


lutBS12.go:
	$(CURL) $(GITLABURL)/lutBS12.bdf > lutBS12.bdf
	$(BDFTOGO) -infile lutBS12.bdf -package $(PACKAGE_NAME) -outfile lutBS12.go
	$(DISPLAYFONT) -infile lutBS12.bdf -outfile lutBS12.png


lutBS14.go:
	$(CURL) $(GITLABURL)/lutBS14.bdf > lutBS14.bdf
	$(BDFTOGO) -infile lutBS14.bdf -package $(PACKAGE_NAME) -outfile lutBS14.go
	$(DISPLAYFONT) -infile lutBS14.bdf -outfile lutBS14.png


lutBS18.go:
	$(CURL) $(GITLABURL)/lutBS18.bdf > lutBS18.bdf
	$(BDFTOGO) -infile lutBS18.bdf -package $(PACKAGE_NAME) -outfile lutBS18.go
	$(DISPLAYFONT) -infile lutBS18.bdf -outfile lutBS18.png


lutBS19.go:
	$(CURL) $(GITLABURL)/lutBS19.bdf > lutBS19.bdf
	$(BDFTOGO) -infile lutBS19.bdf -package $(PACKAGE_NAME) -outfile lutBS19.go
	$(DISPLAYFONT) -infile lutBS19.bdf -outfile lutBS19.png


lutBS24.go:
	$(CURL) $(GITLABURL)/lutBS24.bdf > lutBS24.bdf
	$(BDFTOGO) -infile lutBS24.bdf -package $(PACKAGE_NAME) -outfile lutBS24.go
	$(DISPLAYFONT) -infile lutBS24.bdf -outfile lutBS24.png


lutRS08.go:
	$(CURL) $(GITLABURL)/lutRS08.bdf > lutRS08.bdf
	$(BDFTOGO) -infile lutRS08.bdf -package $(PACKAGE_NAME) -outfile lutRS08.go
	$(DISPLAYFONT) -infile lutRS08.bdf -outfile lutRS08.png


lutRS10.go:
	$(CURL) $(GITLABURL)/lutRS10.bdf > lutRS10.bdf
	$(BDFTOGO) -infile lutRS10.bdf -package $(PACKAGE_NAME) -outfile lutRS10.go
	$(DISPLAYFONT) -infile lutRS10.bdf -outfile lutRS10.png


lutRS12.go:
	$(CURL) $(GITLABURL)/lutRS12.bdf > lutRS12.bdf
	$(BDFTOGO) -infile lutRS12.bdf -package $(PACKAGE_NAME) -outfile lutRS12.go
	$(DISPLAYFONT) -infile lutRS12.bdf -outfile lutRS12.png


lutRS14.go:
	$(CURL) $(GITLABURL)/lutRS14.bdf > lutRS14.bdf
	$(BDFTOGO) -infile lutRS14.bdf -package $(PACKAGE_NAME) -outfile lutRS14.go
	$(DISPLAYFONT) -infile lutRS14.bdf -outfile lutRS14.png


lutRS18.go:
	$(CURL) $(GITLABURL)/lutRS18.bdf > lutRS18.bdf
	$(BDFTOGO) -infile lutRS18.bdf -package $(PACKAGE_NAME) -outfile lutRS18.go
	$(DISPLAYFONT) -infile lutRS18.bdf -outfile lutRS18.png


lutRS19.go:
	$(CURL) $(GITLABURL)/lutRS19.bdf > lutRS19.bdf
	$(BDFTOGO) -infile lutRS19.bdf -package $(PACKAGE_NAME) -outfile lutRS19.go
	$(DISPLAYFONT) -infile lutRS19.bdf -outfile lutRS19.png


lutRS24.go:
	$(CURL) $(GITLABURL)/lutRS24.bdf > lutRS24.bdf
	$(BDFTOGO) -infile lutRS24.bdf -package $(PACKAGE_NAME) -outfile lutRS24.go
	$(DISPLAYFONT) -infile lutRS24.bdf -outfile lutRS24.png

