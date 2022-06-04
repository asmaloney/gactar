GIT_COMMIT=$(shell git describe --always)

all:
	go build -ldflags "-X github.com/asmaloney/gactar/util/version.BuildVersion=${GIT_COMMIT}"