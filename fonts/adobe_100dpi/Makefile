# Makefile for fonts
# This will build go files for each bdf file
# IF you update the bdf flles in this directory, run "make rebuild" to
# update this makefile.

BDFTOGO = ../../bdftogo/bdftogo
CURL = curl --silent
DISPLAYFONT = ../../displayfont/displayfont
GITLABURL = https://gitlab.freedesktop.org/xorg/font/adobe-100dpi/-/raw/master
PACKAGE_NAME = adobe_100dpi

BDF_FILES = courB08.bdf courB10.bdf courB12.bdf courB14.bdf courB18.bdf courB24.bdf \
	courBO08.bdf courBO10.bdf courBO12.bdf courBO14.bdf courBO18.bdf courBO24.bdf \
	courO08.bdf courO10.bdf courO12.bdf courO14.bdf courO18.bdf courO24.bdf courR08.bdf \
	courR10.bdf courR12.bdf courR14.bdf courR18.bdf courR24.bdf helvB08.bdf helvB10.bdf \
	helvB12.bdf helvB14.bdf helvB18.bdf helvB24.bdf helvBO08.bdf helvBO10.bdf helvBO12.bdf \
	helvBO14.bdf helvBO18.bdf helvBO24.bdf helvO08.bdf helvO10.bdf helvO12.bdf helvO14.bdf \
	helvO18.bdf helvO24.bdf helvR08.bdf helvR10.bdf helvR12.bdf helvR14.bdf helvR18.bdf \
	helvR24.bdf ncenB08.bdf ncenB10.bdf ncenB12.bdf ncenB14.bdf ncenB18.bdf ncenB24.bdf \
	ncenBI08.bdf ncenBI10.bdf ncenBI12.bdf ncenBI14.bdf ncenBI18.bdf ncenBI24.bdf ncenI08.bdf \
	ncenI10.bdf ncenI12.bdf ncenI14.bdf ncenI18.bdf ncenI24.bdf ncenR08.bdf ncenR10.bdf \
	ncenR12.bdf ncenR14.bdf ncenR18.bdf ncenR24.bdf symb08.bdf symb10.bdf symb12.bdf \
	symb14.bdf symb18.bdf symb24.bdf timB08.bdf timB10.bdf timB12.bdf timB14.bdf timB18.bdf \
	timB24.bdf timBI08.bdf timBI10.bdf timBI12.bdf timBI14.bdf timBI18.bdf timBI24.bdf \
	timI08.bdf timI10.bdf timI12.bdf timI14.bdf timI18.bdf timI24.bdf timR08.bdf timR10.bdf \
	timR12.bdf timR14.bdf timR18.bdf timR24.bdf

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
	courB08.go \
	courB10.go \
	courB12.go \
	courB14.go \
	courB18.go \
	courB24.go \
	courBO08.go \
	courBO10.go \
	courBO12.go \
	courBO14.go \
	courBO18.go \
	courBO24.go \
	courO08.go \
	courO10.go \
	courO12.go \
	courO14.go \
	courO18.go \
	courO24.go \
	courR08.go \
	courR10.go \
	courR12.go \
	courR14.go \
	courR18.go \
	courR24.go \
	helvB08.go \
	helvB10.go \
	helvB12.go \
	helvB14.go \
	helvB18.go \
	helvB24.go \
	helvBO08.go \
	helvBO10.go \
	helvBO12.go \
	helvBO14.go \
	helvBO18.go \
	helvBO24.go \
	helvO08.go \
	helvO10.go \
	helvO12.go \
	helvO14.go \
	helvO18.go \
	helvO24.go \
	helvR08.go \
	helvR10.go \
	helvR12.go \
	helvR14.go \
	helvR18.go \
	helvR24.go \
	ncenB08.go \
	ncenB10.go \
	ncenB12.go \
	ncenB14.go \
	ncenB18.go \
	ncenB24.go \
	ncenBI08.go \
	ncenBI10.go \
	ncenBI12.go \
	ncenBI14.go \
	ncenBI18.go \
	ncenBI24.go \
	ncenI08.go \
	ncenI10.go \
	ncenI12.go \
	ncenI14.go \
	ncenI18.go \
	ncenI24.go \
	ncenR08.go \
	ncenR10.go \
	ncenR12.go \
	ncenR14.go \
	ncenR18.go \
	ncenR24.go \
	symb08.go \
	symb10.go \
	symb12.go \
	symb14.go \
	symb18.go \
	symb24.go \
	timB08.go \
	timB10.go \
	timB12.go \
	timB14.go \
	timB18.go \
	timB24.go \
	timBI08.go \
	timBI10.go \
	timBI12.go \
	timBI14.go \
	timBI18.go \
	timBI24.go \
	timI08.go \
	timI10.go \
	timI12.go \
	timI14.go \
	timI18.go \
	timI24.go \
	timR08.go \
	timR10.go \
	timR12.go \
	timR14.go \
	timR18.go \
	timR24.go



_GO_FILES: $(GO_FILES)






courB08.go:
	$(CURL) $(GITLABURL)/courB08.bdf > courB08.bdf
	$(BDFTOGO) -infile courB08.bdf -package $(PACKAGE_NAME) -outfile courB08.go
	$(DISPLAYFONT) -infile courB08.bdf -outfile courB08.png


courB10.go:
	$(CURL) $(GITLABURL)/courB10.bdf > courB10.bdf
	$(BDFTOGO) -infile courB10.bdf -package $(PACKAGE_NAME) -outfile courB10.go
	$(DISPLAYFONT) -infile courB10.bdf -outfile courB10.png


courB12.go:
	$(CURL) $(GITLABURL)/courB12.bdf > courB12.bdf
	$(BDFTOGO) -infile courB12.bdf -package $(PACKAGE_NAME) -outfile courB12.go
	$(DISPLAYFONT) -infile courB12.bdf -outfile courB12.png


courB14.go:
	$(CURL) $(GITLABURL)/courB14.bdf > courB14.bdf
	$(BDFTOGO) -infile courB14.bdf -package $(PACKAGE_NAME) -outfile courB14.go
	$(DISPLAYFONT) -infile courB14.bdf -outfile courB14.png


courB18.go:
	$(CURL) $(GITLABURL)/courB18.bdf > courB18.bdf
	$(BDFTOGO) -infile courB18.bdf -package $(PACKAGE_NAME) -outfile courB18.go
	$(DISPLAYFONT) -infile courB18.bdf -outfile courB18.png


courB24.go:
	$(CURL) $(GITLABURL)/courB24.bdf > courB24.bdf
	$(BDFTOGO) -infile courB24.bdf -package $(PACKAGE_NAME) -outfile courB24.go
	$(DISPLAYFONT) -infile courB24.bdf -outfile courB24.png


courBO08.go:
	$(CURL) $(GITLABURL)/courBO08.bdf > courBO08.bdf
	$(BDFTOGO) -infile courBO08.bdf -package $(PACKAGE_NAME) -outfile courBO08.go
	$(DISPLAYFONT) -infile courBO08.bdf -outfile courBO08.png


courBO10.go:
	$(CURL) $(GITLABURL)/courBO10.bdf > courBO10.bdf
	$(BDFTOGO) -infile courBO10.bdf -package $(PACKAGE_NAME) -outfile courBO10.go
	$(DISPLAYFONT) -infile courBO10.bdf -outfile courBO10.png


courBO12.go:
	$(CURL) $(GITLABURL)/courBO12.bdf > courBO12.bdf
	$(BDFTOGO) -infile courBO12.bdf -package $(PACKAGE_NAME) -outfile courBO12.go
	$(DISPLAYFONT) -infile courBO12.bdf -outfile courBO12.png


courBO14.go:
	$(CURL) $(GITLABURL)/courBO14.bdf > courBO14.bdf
	$(BDFTOGO) -infile courBO14.bdf -package $(PACKAGE_NAME) -outfile courBO14.go
	$(DISPLAYFONT) -infile courBO14.bdf -outfile courBO14.png


courBO18.go:
	$(CURL) $(GITLABURL)/courBO18.bdf > courBO18.bdf
	$(BDFTOGO) -infile courBO18.bdf -package $(PACKAGE_NAME) -outfile courBO18.go
	$(DISPLAYFONT) -infile courBO18.bdf -outfile courBO18.png


courBO24.go:
	$(CURL) $(GITLABURL)/courBO24.bdf > courBO24.bdf
	$(BDFTOGO) -infile courBO24.bdf -package $(PACKAGE_NAME) -outfile courBO24.go
	$(DISPLAYFONT) -infile courBO24.bdf -outfile courBO24.png


courO08.go:
	$(CURL) $(GITLABURL)/courO08.bdf > courO08.bdf
	$(BDFTOGO) -infile courO08.bdf -package $(PACKAGE_NAME) -outfile courO08.go
	$(DISPLAYFONT) -infile courO08.bdf -outfile courO08.png


courO10.go:
	$(CURL) $(GITLABURL)/courO10.bdf > courO10.bdf
	$(BDFTOGO) -infile courO10.bdf -package $(PACKAGE_NAME) -outfile courO10.go
	$(DISPLAYFONT) -infile courO10.bdf -outfile courO10.png


courO12.go:
	$(CURL) $(GITLABURL)/courO12.bdf > courO12.bdf
	$(BDFTOGO) -infile courO12.bdf -package $(PACKAGE_NAME) -outfile courO12.go
	$(DISPLAYFONT) -infile courO12.bdf -outfile courO12.png


courO14.go:
	$(CURL) $(GITLABURL)/courO14.bdf > courO14.bdf
	$(BDFTOGO) -infile courO14.bdf -package $(PACKAGE_NAME) -outfile courO14.go
	$(DISPLAYFONT) -infile courO14.bdf -outfile courO14.png


courO18.go:
	$(CURL) $(GITLABURL)/courO18.bdf > courO18.bdf
	$(BDFTOGO) -infile courO18.bdf -package $(PACKAGE_NAME) -outfile courO18.go
	$(DISPLAYFONT) -infile courO18.bdf -outfile courO18.png


courO24.go:
	$(CURL) $(GITLABURL)/courO24.bdf > courO24.bdf
	$(BDFTOGO) -infile courO24.bdf -package $(PACKAGE_NAME) -outfile courO24.go
	$(DISPLAYFONT) -infile courO24.bdf -outfile courO24.png


courR08.go:
	$(CURL) $(GITLABURL)/courR08.bdf > courR08.bdf
	$(BDFTOGO) -infile courR08.bdf -package $(PACKAGE_NAME) -outfile courR08.go
	$(DISPLAYFONT) -infile courR08.bdf -outfile courR08.png


courR10.go:
	$(CURL) $(GITLABURL)/courR10.bdf > courR10.bdf
	$(BDFTOGO) -infile courR10.bdf -package $(PACKAGE_NAME) -outfile courR10.go
	$(DISPLAYFONT) -infile courR10.bdf -outfile courR10.png


courR12.go:
	$(CURL) $(GITLABURL)/courR12.bdf > courR12.bdf
	$(BDFTOGO) -infile courR12.bdf -package $(PACKAGE_NAME) -outfile courR12.go
	$(DISPLAYFONT) -infile courR12.bdf -outfile courR12.png


courR14.go:
	$(CURL) $(GITLABURL)/courR14.bdf > courR14.bdf
	$(BDFTOGO) -infile courR14.bdf -package $(PACKAGE_NAME) -outfile courR14.go
	$(DISPLAYFONT) -infile courR14.bdf -outfile courR14.png


courR18.go:
	$(CURL) $(GITLABURL)/courR18.bdf > courR18.bdf
	$(BDFTOGO) -infile courR18.bdf -package $(PACKAGE_NAME) -outfile courR18.go
	$(DISPLAYFONT) -infile courR18.bdf -outfile courR18.png


courR24.go:
	$(CURL) $(GITLABURL)/courR24.bdf > courR24.bdf
	$(BDFTOGO) -infile courR24.bdf -package $(PACKAGE_NAME) -outfile courR24.go
	$(DISPLAYFONT) -infile courR24.bdf -outfile courR24.png


helvB08.go:
	$(CURL) $(GITLABURL)/helvB08.bdf > helvB08.bdf
	$(BDFTOGO) -infile helvB08.bdf -package $(PACKAGE_NAME) -outfile helvB08.go
	$(DISPLAYFONT) -infile helvB08.bdf -outfile helvB08.png


helvB10.go:
	$(CURL) $(GITLABURL)/helvB10.bdf > helvB10.bdf
	$(BDFTOGO) -infile helvB10.bdf -package $(PACKAGE_NAME) -outfile helvB10.go
	$(DISPLAYFONT) -infile helvB10.bdf -outfile helvB10.png


helvB12.go:
	$(CURL) $(GITLABURL)/helvB12.bdf > helvB12.bdf
	$(BDFTOGO) -infile helvB12.bdf -package $(PACKAGE_NAME) -outfile helvB12.go
	$(DISPLAYFONT) -infile helvB12.bdf -outfile helvB12.png


helvB14.go:
	$(CURL) $(GITLABURL)/helvB14.bdf > helvB14.bdf
	$(BDFTOGO) -infile helvB14.bdf -package $(PACKAGE_NAME) -outfile helvB14.go
	$(DISPLAYFONT) -infile helvB14.bdf -outfile helvB14.png


helvB18.go:
	$(CURL) $(GITLABURL)/helvB18.bdf > helvB18.bdf
	$(BDFTOGO) -infile helvB18.bdf -package $(PACKAGE_NAME) -outfile helvB18.go
	$(DISPLAYFONT) -infile helvB18.bdf -outfile helvB18.png


helvB24.go:
	$(CURL) $(GITLABURL)/helvB24.bdf > helvB24.bdf
	$(BDFTOGO) -infile helvB24.bdf -package $(PACKAGE_NAME) -outfile helvB24.go
	$(DISPLAYFONT) -infile helvB24.bdf -outfile helvB24.png


helvBO08.go:
	$(CURL) $(GITLABURL)/helvBO08.bdf > helvBO08.bdf
	$(BDFTOGO) -infile helvBO08.bdf -package $(PACKAGE_NAME) -outfile helvBO08.go
	$(DISPLAYFONT) -infile helvBO08.bdf -outfile helvBO08.png


helvBO10.go:
	$(CURL) $(GITLABURL)/helvBO10.bdf > helvBO10.bdf
	$(BDFTOGO) -infile helvBO10.bdf -package $(PACKAGE_NAME) -outfile helvBO10.go
	$(DISPLAYFONT) -infile helvBO10.bdf -outfile helvBO10.png


helvBO12.go:
	$(CURL) $(GITLABURL)/helvBO12.bdf > helvBO12.bdf
	$(BDFTOGO) -infile helvBO12.bdf -package $(PACKAGE_NAME) -outfile helvBO12.go
	$(DISPLAYFONT) -infile helvBO12.bdf -outfile helvBO12.png


helvBO14.go:
	$(CURL) $(GITLABURL)/helvBO14.bdf > helvBO14.bdf
	$(BDFTOGO) -infile helvBO14.bdf -package $(PACKAGE_NAME) -outfile helvBO14.go
	$(DISPLAYFONT) -infile helvBO14.bdf -outfile helvBO14.png


helvBO18.go:
	$(CURL) $(GITLABURL)/helvBO18.bdf > helvBO18.bdf
	$(BDFTOGO) -infile helvBO18.bdf -package $(PACKAGE_NAME) -outfile helvBO18.go
	$(DISPLAYFONT) -infile helvBO18.bdf -outfile helvBO18.png


helvBO24.go:
	$(CURL) $(GITLABURL)/helvBO24.bdf > helvBO24.bdf
	$(BDFTOGO) -infile helvBO24.bdf -package $(PACKAGE_NAME) -outfile helvBO24.go
	$(DISPLAYFONT) -infile helvBO24.bdf -outfile helvBO24.png


helvO08.go:
	$(CURL) $(GITLABURL)/helvO08.bdf > helvO08.bdf
	$(BDFTOGO) -infile helvO08.bdf -package $(PACKAGE_NAME) -outfile helvO08.go
	$(DISPLAYFONT) -infile helvO08.bdf -outfile helvO08.png


helvO10.go:
	$(CURL) $(GITLABURL)/helvO10.bdf > helvO10.bdf
	$(BDFTOGO) -infile helvO10.bdf -package $(PACKAGE_NAME) -outfile helvO10.go
	$(DISPLAYFONT) -infile helvO10.bdf -outfile helvO10.png


helvO12.go:
	$(CURL) $(GITLABURL)/helvO12.bdf > helvO12.bdf
	$(BDFTOGO) -infile helvO12.bdf -package $(PACKAGE_NAME) -outfile helvO12.go
	$(DISPLAYFONT) -infile helvO12.bdf -outfile helvO12.png


helvO14.go:
	$(CURL) $(GITLABURL)/helvO14.bdf > helvO14.bdf
	$(BDFTOGO) -infile helvO14.bdf -package $(PACKAGE_NAME) -outfile helvO14.go
	$(DISPLAYFONT) -infile helvO14.bdf -outfile helvO14.png


helvO18.go:
	$(CURL) $(GITLABURL)/helvO18.bdf > helvO18.bdf
	$(BDFTOGO) -infile helvO18.bdf -package $(PACKAGE_NAME) -outfile helvO18.go
	$(DISPLAYFONT) -infile helvO18.bdf -outfile helvO18.png


helvO24.go:
	$(CURL) $(GITLABURL)/helvO24.bdf > helvO24.bdf
	$(BDFTOGO) -infile helvO24.bdf -package $(PACKAGE_NAME) -outfile helvO24.go
	$(DISPLAYFONT) -infile helvO24.bdf -outfile helvO24.png


helvR08.go:
	$(CURL) $(GITLABURL)/helvR08.bdf > helvR08.bdf
	$(BDFTOGO) -infile helvR08.bdf -package $(PACKAGE_NAME) -outfile helvR08.go
	$(DISPLAYFONT) -infile helvR08.bdf -outfile helvR08.png


helvR10.go:
	$(CURL) $(GITLABURL)/helvR10.bdf > helvR10.bdf
	$(BDFTOGO) -infile helvR10.bdf -package $(PACKAGE_NAME) -outfile helvR10.go
	$(DISPLAYFONT) -infile helvR10.bdf -outfile helvR10.png


helvR12.go:
	$(CURL) $(GITLABURL)/helvR12.bdf > helvR12.bdf
	$(BDFTOGO) -infile helvR12.bdf -package $(PACKAGE_NAME) -outfile helvR12.go
	$(DISPLAYFONT) -infile helvR12.bdf -outfile helvR12.png


helvR14.go:
	$(CURL) $(GITLABURL)/helvR14.bdf > helvR14.bdf
	$(BDFTOGO) -infile helvR14.bdf -package $(PACKAGE_NAME) -outfile helvR14.go
	$(DISPLAYFONT) -infile helvR14.bdf -outfile helvR14.png


helvR18.go:
	$(CURL) $(GITLABURL)/helvR18.bdf > helvR18.bdf
	$(BDFTOGO) -infile helvR18.bdf -package $(PACKAGE_NAME) -outfile helvR18.go
	$(DISPLAYFONT) -infile helvR18.bdf -outfile helvR18.png


helvR24.go:
	$(CURL) $(GITLABURL)/helvR24.bdf > helvR24.bdf
	$(BDFTOGO) -infile helvR24.bdf -package $(PACKAGE_NAME) -outfile helvR24.go
	$(DISPLAYFONT) -infile helvR24.bdf -outfile helvR24.png


ncenB08.go:
	$(CURL) $(GITLABURL)/ncenB08.bdf > ncenB08.bdf
	$(BDFTOGO) -infile ncenB08.bdf -package $(PACKAGE_NAME) -outfile ncenB08.go
	$(DISPLAYFONT) -infile ncenB08.bdf -outfile ncenB08.png


ncenB10.go:
	$(CURL) $(GITLABURL)/ncenB10.bdf > ncenB10.bdf
	$(BDFTOGO) -infile ncenB10.bdf -package $(PACKAGE_NAME) -outfile ncenB10.go
	$(DISPLAYFONT) -infile ncenB10.bdf -outfile ncenB10.png


ncenB12.go:
	$(CURL) $(GITLABURL)/ncenB12.bdf > ncenB12.bdf
	$(BDFTOGO) -infile ncenB12.bdf -package $(PACKAGE_NAME) -outfile ncenB12.go
	$(DISPLAYFONT) -infile ncenB12.bdf -outfile ncenB12.png


ncenB14.go:
	$(CURL) $(GITLABURL)/ncenB14.bdf > ncenB14.bdf
	$(BDFTOGO) -infile ncenB14.bdf -package $(PACKAGE_NAME) -outfile ncenB14.go
	$(DISPLAYFONT) -infile ncenB14.bdf -outfile ncenB14.png


ncenB18.go:
	$(CURL) $(GITLABURL)/ncenB18.bdf > ncenB18.bdf
	$(BDFTOGO) -infile ncenB18.bdf -package $(PACKAGE_NAME) -outfile ncenB18.go
	$(DISPLAYFONT) -infile ncenB18.bdf -outfile ncenB18.png


ncenB24.go:
	$(CURL) $(GITLABURL)/ncenB24.bdf > ncenB24.bdf
	$(BDFTOGO) -infile ncenB24.bdf -package $(PACKAGE_NAME) -outfile ncenB24.go
	$(DISPLAYFONT) -infile ncenB24.bdf -outfile ncenB24.png


ncenBI08.go:
	$(CURL) $(GITLABURL)/ncenBI08.bdf > ncenBI08.bdf
	$(BDFTOGO) -infile ncenBI08.bdf -package $(PACKAGE_NAME) -outfile ncenBI08.go
	$(DISPLAYFONT) -infile ncenBI08.bdf -outfile ncenBI08.png


ncenBI10.go:
	$(CURL) $(GITLABURL)/ncenBI10.bdf > ncenBI10.bdf
	$(BDFTOGO) -infile ncenBI10.bdf -package $(PACKAGE_NAME) -outfile ncenBI10.go
	$(DISPLAYFONT) -infile ncenBI10.bdf -outfile ncenBI10.png


ncenBI12.go:
	$(CURL) $(GITLABURL)/ncenBI12.bdf > ncenBI12.bdf
	$(BDFTOGO) -infile ncenBI12.bdf -package $(PACKAGE_NAME) -outfile ncenBI12.go
	$(DISPLAYFONT) -infile ncenBI12.bdf -outfile ncenBI12.png


ncenBI14.go:
	$(CURL) $(GITLABURL)/ncenBI14.bdf > ncenBI14.bdf
	$(BDFTOGO) -infile ncenBI14.bdf -package $(PACKAGE_NAME) -outfile ncenBI14.go
	$(DISPLAYFONT) -infile ncenBI14.bdf -outfile ncenBI14.png


ncenBI18.go:
	$(CURL) $(GITLABURL)/ncenBI18.bdf > ncenBI18.bdf
	$(BDFTOGO) -infile ncenBI18.bdf -package $(PACKAGE_NAME) -outfile ncenBI18.go
	$(DISPLAYFONT) -infile ncenBI18.bdf -outfile ncenBI18.png


ncenBI24.go:
	$(CURL) $(GITLABURL)/ncenBI24.bdf > ncenBI24.bdf
	$(BDFTOGO) -infile ncenBI24.bdf -package $(PACKAGE_NAME) -outfile ncenBI24.go
	$(DISPLAYFONT) -infile ncenBI24.bdf -outfile ncenBI24.png


ncenI08.go:
	$(CURL) $(GITLABURL)/ncenI08.bdf > ncenI08.bdf
	$(BDFTOGO) -infile ncenI08.bdf -package $(PACKAGE_NAME) -outfile ncenI08.go
	$(DISPLAYFONT) -infile ncenI08.bdf -outfile ncenI08.png


ncenI10.go:
	$(CURL) $(GITLABURL)/ncenI10.bdf > ncenI10.bdf
	$(BDFTOGO) -infile ncenI10.bdf -package $(PACKAGE_NAME) -outfile ncenI10.go
	$(DISPLAYFONT) -infile ncenI10.bdf -outfile ncenI10.png


ncenI12.go:
	$(CURL) $(GITLABURL)/ncenI12.bdf > ncenI12.bdf
	$(BDFTOGO) -infile ncenI12.bdf -package $(PACKAGE_NAME) -outfile ncenI12.go
	$(DISPLAYFONT) -infile ncenI12.bdf -outfile ncenI12.png


ncenI14.go:
	$(CURL) $(GITLABURL)/ncenI14.bdf > ncenI14.bdf
	$(BDFTOGO) -infile ncenI14.bdf -package $(PACKAGE_NAME) -outfile ncenI14.go
	$(DISPLAYFONT) -infile ncenI14.bdf -outfile ncenI14.png


ncenI18.go:
	$(CURL) $(GITLABURL)/ncenI18.bdf > ncenI18.bdf
	$(BDFTOGO) -infile ncenI18.bdf -package $(PACKAGE_NAME) -outfile ncenI18.go
	$(DISPLAYFONT) -infile ncenI18.bdf -outfile ncenI18.png


ncenI24.go:
	$(CURL) $(GITLABURL)/ncenI24.bdf > ncenI24.bdf
	$(BDFTOGO) -infile ncenI24.bdf -package $(PACKAGE_NAME) -outfile ncenI24.go
	$(DISPLAYFONT) -infile ncenI24.bdf -outfile ncenI24.png


ncenR08.go:
	$(CURL) $(GITLABURL)/ncenR08.bdf > ncenR08.bdf
	$(BDFTOGO) -infile ncenR08.bdf -package $(PACKAGE_NAME) -outfile ncenR08.go
	$(DISPLAYFONT) -infile ncenR08.bdf -outfile ncenR08.png


ncenR10.go:
	$(CURL) $(GITLABURL)/ncenR10.bdf > ncenR10.bdf
	$(BDFTOGO) -infile ncenR10.bdf -package $(PACKAGE_NAME) -outfile ncenR10.go
	$(DISPLAYFONT) -infile ncenR10.bdf -outfile ncenR10.png


ncenR12.go:
	$(CURL) $(GITLABURL)/ncenR12.bdf > ncenR12.bdf
	$(BDFTOGO) -infile ncenR12.bdf -package $(PACKAGE_NAME) -outfile ncenR12.go
	$(DISPLAYFONT) -infile ncenR12.bdf -outfile ncenR12.png


ncenR14.go:
	$(CURL) $(GITLABURL)/ncenR14.bdf > ncenR14.bdf
	$(BDFTOGO) -infile ncenR14.bdf -package $(PACKAGE_NAME) -outfile ncenR14.go
	$(DISPLAYFONT) -infile ncenR14.bdf -outfile ncenR14.png


ncenR18.go:
	$(CURL) $(GITLABURL)/ncenR18.bdf > ncenR18.bdf
	$(BDFTOGO) -infile ncenR18.bdf -package $(PACKAGE_NAME) -outfile ncenR18.go
	$(DISPLAYFONT) -infile ncenR18.bdf -outfile ncenR18.png


ncenR24.go:
	$(CURL) $(GITLABURL)/ncenR24.bdf > ncenR24.bdf
	$(BDFTOGO) -infile ncenR24.bdf -package $(PACKAGE_NAME) -outfile ncenR24.go
	$(DISPLAYFONT) -infile ncenR24.bdf -outfile ncenR24.png


symb08.go:
	$(CURL) $(GITLABURL)/symb08.bdf > symb08.bdf
	$(BDFTOGO) -infile symb08.bdf -package $(PACKAGE_NAME) -outfile symb08.go
	$(DISPLAYFONT) -infile symb08.bdf -outfile symb08.png


symb10.go:
	$(CURL) $(GITLABURL)/symb10.bdf > symb10.bdf
	$(BDFTOGO) -infile symb10.bdf -package $(PACKAGE_NAME) -outfile symb10.go
	$(DISPLAYFONT) -infile symb10.bdf -outfile symb10.png


symb12.go:
	$(CURL) $(GITLABURL)/symb12.bdf > symb12.bdf
	$(BDFTOGO) -infile symb12.bdf -package $(PACKAGE_NAME) -outfile symb12.go
	$(DISPLAYFONT) -infile symb12.bdf -outfile symb12.png


symb14.go:
	$(CURL) $(GITLABURL)/symb14.bdf > symb14.bdf
	$(BDFTOGO) -infile symb14.bdf -package $(PACKAGE_NAME) -outfile symb14.go
	$(DISPLAYFONT) -infile symb14.bdf -outfile symb14.png


symb18.go:
	$(CURL) $(GITLABURL)/symb18.bdf > symb18.bdf
	$(BDFTOGO) -infile symb18.bdf -package $(PACKAGE_NAME) -outfile symb18.go
	$(DISPLAYFONT) -infile symb18.bdf -outfile symb18.png


symb24.go:
	$(CURL) $(GITLABURL)/symb24.bdf > symb24.bdf
	$(BDFTOGO) -infile symb24.bdf -package $(PACKAGE_NAME) -outfile symb24.go
	$(DISPLAYFONT) -infile symb24.bdf -outfile symb24.png


timB08.go:
	$(CURL) $(GITLABURL)/timB08.bdf > timB08.bdf
	$(BDFTOGO) -infile timB08.bdf -package $(PACKAGE_NAME) -outfile timB08.go
	$(DISPLAYFONT) -infile timB08.bdf -outfile timB08.png


timB10.go:
	$(CURL) $(GITLABURL)/timB10.bdf > timB10.bdf
	$(BDFTOGO) -infile timB10.bdf -package $(PACKAGE_NAME) -outfile timB10.go
	$(DISPLAYFONT) -infile timB10.bdf -outfile timB10.png


timB12.go:
	$(CURL) $(GITLABURL)/timB12.bdf > timB12.bdf
	$(BDFTOGO) -infile timB12.bdf -package $(PACKAGE_NAME) -outfile timB12.go
	$(DISPLAYFONT) -infile timB12.bdf -outfile timB12.png


timB14.go:
	$(CURL) $(GITLABURL)/timB14.bdf > timB14.bdf
	$(BDFTOGO) -infile timB14.bdf -package $(PACKAGE_NAME) -outfile timB14.go
	$(DISPLAYFONT) -infile timB14.bdf -outfile timB14.png


timB18.go:
	$(CURL) $(GITLABURL)/timB18.bdf > timB18.bdf
	$(BDFTOGO) -infile timB18.bdf -package $(PACKAGE_NAME) -outfile timB18.go
	$(DISPLAYFONT) -infile timB18.bdf -outfile timB18.png


timB24.go:
	$(CURL) $(GITLABURL)/timB24.bdf > timB24.bdf
	$(BDFTOGO) -infile timB24.bdf -package $(PACKAGE_NAME) -outfile timB24.go
	$(DISPLAYFONT) -infile timB24.bdf -outfile timB24.png


timBI08.go:
	$(CURL) $(GITLABURL)/timBI08.bdf > timBI08.bdf
	$(BDFTOGO) -infile timBI08.bdf -package $(PACKAGE_NAME) -outfile timBI08.go
	$(DISPLAYFONT) -infile timBI08.bdf -outfile timBI08.png


timBI10.go:
	$(CURL) $(GITLABURL)/timBI10.bdf > timBI10.bdf
	$(BDFTOGO) -infile timBI10.bdf -package $(PACKAGE_NAME) -outfile timBI10.go
	$(DISPLAYFONT) -infile timBI10.bdf -outfile timBI10.png


timBI12.go:
	$(CURL) $(GITLABURL)/timBI12.bdf > timBI12.bdf
	$(BDFTOGO) -infile timBI12.bdf -package $(PACKAGE_NAME) -outfile timBI12.go
	$(DISPLAYFONT) -infile timBI12.bdf -outfile timBI12.png


timBI14.go:
	$(CURL) $(GITLABURL)/timBI14.bdf > timBI14.bdf
	$(BDFTOGO) -infile timBI14.bdf -package $(PACKAGE_NAME) -outfile timBI14.go
	$(DISPLAYFONT) -infile timBI14.bdf -outfile timBI14.png


timBI18.go:
	$(CURL) $(GITLABURL)/timBI18.bdf > timBI18.bdf
	$(BDFTOGO) -infile timBI18.bdf -package $(PACKAGE_NAME) -outfile timBI18.go
	$(DISPLAYFONT) -infile timBI18.bdf -outfile timBI18.png


timBI24.go:
	$(CURL) $(GITLABURL)/timBI24.bdf > timBI24.bdf
	$(BDFTOGO) -infile timBI24.bdf -package $(PACKAGE_NAME) -outfile timBI24.go
	$(DISPLAYFONT) -infile timBI24.bdf -outfile timBI24.png


timI08.go:
	$(CURL) $(GITLABURL)/timI08.bdf > timI08.bdf
	$(BDFTOGO) -infile timI08.bdf -package $(PACKAGE_NAME) -outfile timI08.go
	$(DISPLAYFONT) -infile timI08.bdf -outfile timI08.png


timI10.go:
	$(CURL) $(GITLABURL)/timI10.bdf > timI10.bdf
	$(BDFTOGO) -infile timI10.bdf -package $(PACKAGE_NAME) -outfile timI10.go
	$(DISPLAYFONT) -infile timI10.bdf -outfile timI10.png


timI12.go:
	$(CURL) $(GITLABURL)/timI12.bdf > timI12.bdf
	$(BDFTOGO) -infile timI12.bdf -package $(PACKAGE_NAME) -outfile timI12.go
	$(DISPLAYFONT) -infile timI12.bdf -outfile timI12.png


timI14.go:
	$(CURL) $(GITLABURL)/timI14.bdf > timI14.bdf
	$(BDFTOGO) -infile timI14.bdf -package $(PACKAGE_NAME) -outfile timI14.go
	$(DISPLAYFONT) -infile timI14.bdf -outfile timI14.png


timI18.go:
	$(CURL) $(GITLABURL)/timI18.bdf > timI18.bdf
	$(BDFTOGO) -infile timI18.bdf -package $(PACKAGE_NAME) -outfile timI18.go
	$(DISPLAYFONT) -infile timI18.bdf -outfile timI18.png


timI24.go:
	$(CURL) $(GITLABURL)/timI24.bdf > timI24.bdf
	$(BDFTOGO) -infile timI24.bdf -package $(PACKAGE_NAME) -outfile timI24.go
	$(DISPLAYFONT) -infile timI24.bdf -outfile timI24.png


timR08.go:
	$(CURL) $(GITLABURL)/timR08.bdf > timR08.bdf
	$(BDFTOGO) -infile timR08.bdf -package $(PACKAGE_NAME) -outfile timR08.go
	$(DISPLAYFONT) -infile timR08.bdf -outfile timR08.png


timR10.go:
	$(CURL) $(GITLABURL)/timR10.bdf > timR10.bdf
	$(BDFTOGO) -infile timR10.bdf -package $(PACKAGE_NAME) -outfile timR10.go
	$(DISPLAYFONT) -infile timR10.bdf -outfile timR10.png


timR12.go:
	$(CURL) $(GITLABURL)/timR12.bdf > timR12.bdf
	$(BDFTOGO) -infile timR12.bdf -package $(PACKAGE_NAME) -outfile timR12.go
	$(DISPLAYFONT) -infile timR12.bdf -outfile timR12.png


timR14.go:
	$(CURL) $(GITLABURL)/timR14.bdf > timR14.bdf
	$(BDFTOGO) -infile timR14.bdf -package $(PACKAGE_NAME) -outfile timR14.go
	$(DISPLAYFONT) -infile timR14.bdf -outfile timR14.png


timR18.go:
	$(CURL) $(GITLABURL)/timR18.bdf > timR18.bdf
	$(BDFTOGO) -infile timR18.bdf -package $(PACKAGE_NAME) -outfile timR18.go
	$(DISPLAYFONT) -infile timR18.bdf -outfile timR18.png


timR24.go:
	$(CURL) $(GITLABURL)/timR24.bdf > timR24.bdf
	$(BDFTOGO) -infile timR24.bdf -package $(PACKAGE_NAME) -outfile timR24.go
	$(DISPLAYFONT) -infile timR24.bdf -outfile timR24.png

