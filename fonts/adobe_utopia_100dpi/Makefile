# Makefile for fonts
# This will build go files for each bdf file
# IF you update the bdf flles in this directory, run "make rebuild" to
# update this makefile.

BDFTOGO = ../../bdftogo/bdftogo
CURL = curl --silent
DISPLAYFONT = ../../displayfont/displayfont
GITLABURL = https://gitlab.freedesktop.org/xorg/font/adobe-utopia-100dpi/-/raw/master
PACKAGE_NAME = adobe_utopia_100dpi

BDF_FILES = UTBI__10.bdf UTBI__12.bdf UTBI__14.bdf UTBI__18.bdf UTBI__24.bdf \
	UTB___10.bdf UTB___12.bdf UTB___14.bdf UTB___18.bdf UTB___24.bdf UTI___10.bdf \
	UTI___12.bdf UTI___14.bdf UTI___18.bdf UTI___24.bdf UTRG__10.bdf UTRG__12.bdf \
	UTRG__14.bdf UTRG__18.bdf UTRG__24.bdf

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
	UTBI__10.go \
	UTBI__12.go \
	UTBI__14.go \
	UTBI__18.go \
	UTBI__24.go \
	UTB___10.go \
	UTB___12.go \
	UTB___14.go \
	UTB___18.go \
	UTB___24.go \
	UTI___10.go \
	UTI___12.go \
	UTI___14.go \
	UTI___18.go \
	UTI___24.go \
	UTRG__10.go \
	UTRG__12.go \
	UTRG__14.go \
	UTRG__18.go \
	UTRG__24.go



_GO_FILES: $(GO_FILES)






UTBI__10.go:
	$(CURL) $(GITLABURL)/UTBI__10.bdf > UTBI__10.bdf
	$(BDFTOGO) -infile UTBI__10.bdf -package $(PACKAGE_NAME) -outfile UTBI__10.go
	$(DISPLAYFONT) -infile UTBI__10.bdf -outfile UTBI__10.png


UTBI__12.go:
	$(CURL) $(GITLABURL)/UTBI__12.bdf > UTBI__12.bdf
	$(BDFTOGO) -infile UTBI__12.bdf -package $(PACKAGE_NAME) -outfile UTBI__12.go
	$(DISPLAYFONT) -infile UTBI__12.bdf -outfile UTBI__12.png


UTBI__14.go:
	$(CURL) $(GITLABURL)/UTBI__14.bdf > UTBI__14.bdf
	$(BDFTOGO) -infile UTBI__14.bdf -package $(PACKAGE_NAME) -outfile UTBI__14.go
	$(DISPLAYFONT) -infile UTBI__14.bdf -outfile UTBI__14.png


UTBI__18.go:
	$(CURL) $(GITLABURL)/UTBI__18.bdf > UTBI__18.bdf
	$(BDFTOGO) -infile UTBI__18.bdf -package $(PACKAGE_NAME) -outfile UTBI__18.go
	$(DISPLAYFONT) -infile UTBI__18.bdf -outfile UTBI__18.png


UTBI__24.go:
	$(CURL) $(GITLABURL)/UTBI__24.bdf > UTBI__24.bdf
	$(BDFTOGO) -infile UTBI__24.bdf -package $(PACKAGE_NAME) -outfile UTBI__24.go
	$(DISPLAYFONT) -infile UTBI__24.bdf -outfile UTBI__24.png


UTB___10.go:
	$(CURL) $(GITLABURL)/UTB___10.bdf > UTB___10.bdf
	$(BDFTOGO) -infile UTB___10.bdf -package $(PACKAGE_NAME) -outfile UTB___10.go
	$(DISPLAYFONT) -infile UTB___10.bdf -outfile UTB___10.png


UTB___12.go:
	$(CURL) $(GITLABURL)/UTB___12.bdf > UTB___12.bdf
	$(BDFTOGO) -infile UTB___12.bdf -package $(PACKAGE_NAME) -outfile UTB___12.go
	$(DISPLAYFONT) -infile UTB___12.bdf -outfile UTB___12.png


UTB___14.go:
	$(CURL) $(GITLABURL)/UTB___14.bdf > UTB___14.bdf
	$(BDFTOGO) -infile UTB___14.bdf -package $(PACKAGE_NAME) -outfile UTB___14.go
	$(DISPLAYFONT) -infile UTB___14.bdf -outfile UTB___14.png


UTB___18.go:
	$(CURL) $(GITLABURL)/UTB___18.bdf > UTB___18.bdf
	$(BDFTOGO) -infile UTB___18.bdf -package $(PACKAGE_NAME) -outfile UTB___18.go
	$(DISPLAYFONT) -infile UTB___18.bdf -outfile UTB___18.png


UTB___24.go:
	$(CURL) $(GITLABURL)/UTB___24.bdf > UTB___24.bdf
	$(BDFTOGO) -infile UTB___24.bdf -package $(PACKAGE_NAME) -outfile UTB___24.go
	$(DISPLAYFONT) -infile UTB___24.bdf -outfile UTB___24.png


UTI___10.go:
	$(CURL) $(GITLABURL)/UTI___10.bdf > UTI___10.bdf
	$(BDFTOGO) -infile UTI___10.bdf -package $(PACKAGE_NAME) -outfile UTI___10.go
	$(DISPLAYFONT) -infile UTI___10.bdf -outfile UTI___10.png


UTI___12.go:
	$(CURL) $(GITLABURL)/UTI___12.bdf > UTI___12.bdf
	$(BDFTOGO) -infile UTI___12.bdf -package $(PACKAGE_NAME) -outfile UTI___12.go
	$(DISPLAYFONT) -infile UTI___12.bdf -outfile UTI___12.png


UTI___14.go:
	$(CURL) $(GITLABURL)/UTI___14.bdf > UTI___14.bdf
	$(BDFTOGO) -infile UTI___14.bdf -package $(PACKAGE_NAME) -outfile UTI___14.go
	$(DISPLAYFONT) -infile UTI___14.bdf -outfile UTI___14.png


UTI___18.go:
	$(CURL) $(GITLABURL)/UTI___18.bdf > UTI___18.bdf
	$(BDFTOGO) -infile UTI___18.bdf -package $(PACKAGE_NAME) -outfile UTI___18.go
	$(DISPLAYFONT) -infile UTI___18.bdf -outfile UTI___18.png


UTI___24.go:
	$(CURL) $(GITLABURL)/UTI___24.bdf > UTI___24.bdf
	$(BDFTOGO) -infile UTI___24.bdf -package $(PACKAGE_NAME) -outfile UTI___24.go
	$(DISPLAYFONT) -infile UTI___24.bdf -outfile UTI___24.png


UTRG__10.go:
	$(CURL) $(GITLABURL)/UTRG__10.bdf > UTRG__10.bdf
	$(BDFTOGO) -infile UTRG__10.bdf -package $(PACKAGE_NAME) -outfile UTRG__10.go
	$(DISPLAYFONT) -infile UTRG__10.bdf -outfile UTRG__10.png


UTRG__12.go:
	$(CURL) $(GITLABURL)/UTRG__12.bdf > UTRG__12.bdf
	$(BDFTOGO) -infile UTRG__12.bdf -package $(PACKAGE_NAME) -outfile UTRG__12.go
	$(DISPLAYFONT) -infile UTRG__12.bdf -outfile UTRG__12.png


UTRG__14.go:
	$(CURL) $(GITLABURL)/UTRG__14.bdf > UTRG__14.bdf
	$(BDFTOGO) -infile UTRG__14.bdf -package $(PACKAGE_NAME) -outfile UTRG__14.go
	$(DISPLAYFONT) -infile UTRG__14.bdf -outfile UTRG__14.png


UTRG__18.go:
	$(CURL) $(GITLABURL)/UTRG__18.bdf > UTRG__18.bdf
	$(BDFTOGO) -infile UTRG__18.bdf -package $(PACKAGE_NAME) -outfile UTRG__18.go
	$(DISPLAYFONT) -infile UTRG__18.bdf -outfile UTRG__18.png


UTRG__24.go:
	$(CURL) $(GITLABURL)/UTRG__24.bdf > UTRG__24.bdf
	$(BDFTOGO) -infile UTRG__24.bdf -package $(PACKAGE_NAME) -outfile UTRG__24.go
	$(DISPLAYFONT) -infile UTRG__24.bdf -outfile UTRG__24.png

