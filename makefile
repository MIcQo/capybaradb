ARGS = CGO_ENABLED=0 GOGC=off
BINARY_NAME = testtool

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

TAG_NAME := $(shell git describe --abbrev=0 --tags --exact-match)
SHA := $(shell git rev-parse HEAD)
VERSION_GIT := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
VERSION := $(if $(VERSION),$(VERSION),$(VERSION_GIT))
CODENAME ?= cheddar
DATE := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

GOOS ?= darwin
GOARCH ?= arm64

build-binary:
	@$(ARGS) GOOS=${GOOS} GOARCH=${GOARCH}  go build -ldflags "-s -w \
		-X capybaradb/cmd.Version=$(VERSION) \
		-X capybaradb/cmd.Codename=$(CODENAME) \
		-X capybaradb/cmd.BuildDate=$(DATE)" \
		-o "./build/${GOOS}/${GOARCH}/"

	@zip -r release/${BINARY_NAME}-${VERSION}-${GOOS}-${GOARCH}.zip build/${GOOS}/${GOARCH}/*

create-build-dir:
	rm -rdf ./build
	rm -rdf ./release
	mkdir -p ./release ./build

build-all: create-build-dir build-darwin build-linux build-windows

build-darwin: build-darwin-amd64 build-darwin-arm64
build-linux: build-linux-amd64 build-linux-arm64
build-windows: build-windows-amd64 build-windows-arm64

build-windows-amd64:
	echo "Building windows/amd64"
	@GOOS=windows GOARCH=amd64 $(MAKE) build-binary
	echo "windows/amd64 done\n"

build-windows-arm64:
	echo "Building windows/arm64"
	@GOOS=windows GOARCH=arm64 $(MAKE) build-binary
	echo "windows/arm64 done\n"

build-linux-amd64:
	echo "Building linux/amd64"
	@GOOS=linux GOARCH=amd64 $(MAKE) build-binary
	echo "linux/amd64 done\n"

build-linux-arm64:
	echo "Building linux/arm64"
	@GOOS=linux GOARCH=arm64 $(MAKE) build-binary
	echo "linux/arm64 done\n"

build-darwin-amd64:
	echo "Building darwin/amd64"
	@GOOS=darwin GOARCH=amd64 $(MAKE) build-binary
	echo "darwin/amd64 done\n"

build-darwin-arm64:
	echo "Building darwin/arm64"
	@GOOS=darwin GOARCH=arm64 $(MAKE) build-binary
	echo "darwin/arm64 done\n"