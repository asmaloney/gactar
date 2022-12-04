GIT_COMMIT=$(shell git describe --always)

.PHONY: all build clean test

default: build

all: build test

build:
	go build -ldflags "-X github.com/asmaloney/gactar/util/version.BuildVersion=${GIT_COMMIT}"

clean:
	rm ./gactar

test:
	go test ./...