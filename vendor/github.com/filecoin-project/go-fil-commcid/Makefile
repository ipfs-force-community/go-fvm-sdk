all: build
.PHONY: all

SUBMODULES=

commcid:
	go build ./...
.PHONY: filestore
SUBMODULES+=commcid

build: $(SUBMODULES)

clean: