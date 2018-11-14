NAME := misoapp
APPPATH := github.com/mamemomonga/miso-model1

SOUNDS    := $(shell find sounds -type f -name '*.wav')
SRCS_MISOAPP := $(shell find go/hardware go/cmd/misoapp -type f -name '*.go')
SRCS_PWRCTRL := $(shell find go/hardware go/cmd/power-controller -type f -name '*.go')
VERSION   := v$(shell cat version)
REVISION  := $(shell git rev-parse --short HEAD)

# LDFLAGS   := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
# BUILDARGS := -a -tags netgo -installsuffix netgo $(LDFLAGS)

LDFLAGS   := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\""
BUILDARGS := $(LDFLAGS)

SOUNDS_DIST := $(patsubst sounds/%, dist/sounds/%, $(SOUNDS))

export GOBIN := $(shell if [ -z "$$GOBIN" ]; then echo "$$GOPATH/bin"; else echo "$$GOBIN"; fi)

all: \
	dist/power-controller \
	dist/misoapp \
	dist/setup.sh \
	dist/config.yaml \
	$(SOUNDS_DIST)

clean:
	rm -rf dist

dist/power-controller: $(SRCS_PWRCTRL)
	mkdir -p dist
	cd go/cmd/power-controller; go build $(BUILDARGS) -o ../../../$@ .

dist/misoapp: $(SRCS_MISOAPP)
	mkdir -p dist
	cd go/cmd/misoapp; go build $(BUILDARGS) -o ../../../$@ .

dist/setup.sh: script/setup.sh
	cp $< $@

dist/config.yaml: etc/config.yaml
	cp $< $@

$(SOUNDS_DIST): $(SOUNDS)
	@if [ ! -e "dist/sounds" ]; then mkdir -p dist/sounds; fi
	cp -f $(patsubst  dist/sounds/%, sounds/%, $@) $@

# ---------------------

include remote
REMOTE_DIR := /usr/local/misoapp

remote-send: all
	ssh $(REMOTE_SSH) sudo mkdir -p $(REMOTE_DIR)
	rsync -av -e 'ssh' --rsync-path='sudo rsync' --exclude='*.swp' \
		$(CURDIR)/dist/ $(REMOTE_SSH):$(REMOTE_DIR)/

remote-run-misoapp: remote-send
	-ssh -t $(REMOTE_SSH) sudo systemctl stop misoapp
	ssh -t $(REMOTE_SSH) sudo $(REMOTE_DIR)/misoapp

remote-run-power-controller: remote-send
	-ssh -t $(REMOTE_SSH) sudo systemctl stop misoapp 
	ssh -t $(REMOTE_SSH) sudo $(REMOTE_DIR)/power-controller


