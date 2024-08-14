NAME = hole
VERSION=$(shell cat VERSION 2>/dev/null || echo "unknown version")
BUILDTIME=$(shell date -u)
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown commit")
GOBUILDFLAGS = CGO_ENABLED=0 go build -trimpath -ldflags '-X "main.Version=$(VERSION)" -X "main.Commit=$(COMMIT)" -X "main.BuildTime=$(BUILDTIME)" -w -s'

PLATFORM_LIST = \
	linux-amd64 \
	linux-arm64 \
	darwin-amd64 \
	darwin-arm64

WINDOWS_LIST = \
	windows-amd64 \
	windows-arm64

.DEFAULT_GOAL := help

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build:  ## build the binary
	$(GOBUILDFLAGS) -o bin/$(NAME) main.go

.PHONY: linux-amd64
linux-amd64:  ## build the binary for linux/amd64
	GOOS=linux GOARCH=amd64 $(GOBUILDFLAGS) -o bin/$(NAME)-$@ main.go

.PHONY: linux-arm64
linux-arm64:  ## build the binary for linux/arm64
	GOOS=linux GOARCH=arm64 $(GOBUILDFLAGS) -o bin/$(NAME)-$@ main.go

.PHONY: windows-amd64
windows-amd64:  ## build the binary for windows/amd64
	GOOS=windows GOARCH=amd64 $(GOBUILDFLAGS) -o bin/$(NAME)-$@.exe main.go

.PHONY: windows-arm64
windows-arm64:  ## build the binary for windows/arm64
	GOOS=windows GOARCH=arm64 $(GOBUILDFLAGS) -o bin/$(NAME)-$@.exe main.go

.PHONY: darwin-amd64
darwin-amd64:  ## build the binary for darwin/amd64
	GOOS=darwin GOARCH=amd64 $(GOBUILDFLAGS) -o bin/$(NAME)-$@ main.go

.PHONY: darwin-arm64
darwin-arm64:  ## build the binary for darwin/arm64
	GOOS=darwin GOARCH=arm64 $(GOBUILDFLAGS) -o bin/$(NAME)-$@ main.go

tar_release=$(addsuffix .tgz, $(PLATFORM_LIST))
zip_release=$(addsuffix .zip, $(WINDOWS_LIST))

$(tar_release): %.tgz: %
	chmod +x bin/$(NAME)-$(basename $@)
	tar -zcvf bin/$(NAME)-$(basename $@).tgz -C bin $(NAME)-$(basename $@)
	rm -f bin/$(NAME)-$(basename $@)

$(zip_release): %.zip: %
	zip -m -j bin/$(NAME)-$(basename $@).zip bin/$(NAME)-$(basename $@).exe
	rm -f bin/$(NAME)-$(basename $@).exe

.PHONY: release
release: clean $(tar_release) $(zip_release)

.PHONY: clean
clean:  ## clean the binary
	rm -rf bin