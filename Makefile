
EPOCH_TEST_COMMIT	:= 78e6667ae2d67aad100b28ee9580b41b7a24e667
OUTPUT_DIRNAME		?= output
DOC_FILENAME		?= oci-runtime-spec
A2X			?= $(shell command -v a2x 2>/dev/null)
DBLATEX			?= $(shell command -v dblatex 2>/dev/null)

# These docs are in an order that determines how they show up in the PDF/HTML docs.
DOC_FILES := $(wildcard *.asc)

default: docs

.PHONY: docs
docs: $(OUTPUT_DIRNAME)/$(DOC_FILENAME).pdf $(OUTPUT_DIRNAME)/$(DOC_FILENAME).html

ifeq "$(strip $(A2X) $(DBLATEX))" ''
$(OUTPUT_DIRNAME)/$(DOC_FILENAME).pdf:
	$(error cannot build $@ without a2x and dblatex)
else
$(OUTPUT_DIRNAME)/$(DOC_FILENAME).pdf: $(DOC_FILES)
	mkdir -p $(OUTPUT_DIRNAME)
	$(A2X) -f pdf --no-xmllint $(ASCIIDOC_SRC)$(DOC_FILENAME).asc
	mv $(ASCIIDOC_SRC)$(DOC_FILENAME).pdf $@
endif

ifeq "$(strip $(A2X))" ''
$(OUTPUT_DIRNAME)/$(DOC_FILENAME).html:
	$(error cannot build $@ without a2x)
else
$(OUTPUT_DIRNAME)/$(DOC_FILENAME).html: $(DOC_FILES)
	mkdir -p $(OUTPUT_DIRNAME)
	$(A2X) -f xhtml --no-xmllint -D $(OUTPUT_DIRNAME) $(ASCIIDOC_SRC)$(DOC_FILENAME).asc
endif

code-of-conduct.md:
	curl -o $@ https://raw.githubusercontent.com/opencontainers/tob/d2f9d68c1332870e40693fe077d311e0742bc73d/code-of-conduct.md

version.md: ./specs-go/version.go
	go run ./.tool/version-doc.go > $@

HOST_GOLANG_VERSION	= $(shell go version | cut -d ' ' -f3 | cut -c 3-)
# this variable is used like a function. First arg is the minimum version, Second arg is the version to be checked.
ALLOWED_GO_VERSION	= $(shell test '$(shell /bin/echo -e "$(1)\n$(2)" | sort -V | head -n1)' = '$(1)' && echo 'true')

.PHONY: test .govet .golint .gitvalidation

test: .govet .golint .gitvalidation

.govet:
	go vet -x ./...

# `go get github.com/golang/lint/golint`
.golint:
ifeq ($(call ALLOWED_GO_VERSION,1.6,$(HOST_GOLANG_VERSION)),true)
	@which golint > /dev/null 2>/dev/null || (echo "ERROR: golint not found. Consider 'make install.tools' target" && false)
	golint ./...
endif


# When this is running in travis, it will only check the travis commit range
.gitvalidation:
	@which git-validation > /dev/null 2>/dev/null || (echo "ERROR: git-validation not found. Consider 'make install.tools' target" && false)
ifeq ($(TRAVIS),true)
	git-validation -q -run DCO,short-subject,dangling-whitespace
else
	git-validation -v -run DCO,short-subject,dangling-whitespace -range $(EPOCH_TEST_COMMIT)..HEAD
endif


.PHONY: install.tools
install.tools: .install.golint .install.gitvalidation

# golint does not even build for <go1.6
.install.golint:
ifeq ($(call ALLOWED_GO_VERSION,1.6,$(HOST_GOLANG_VERSION)),true)
	go get -u github.com/golang/lint/golint
endif

.install.gitvalidation:
	go get -u github.com/vbatts/git-validation


.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIRNAME) *~
	rm -f code-of-conduct.md version.md

