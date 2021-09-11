GIT_COMMIT=$(shell git describe --always)

all:
	go build -ldflags "-X github.com/asmaloney/gactar/version.BuildVersion=${GIT_COMMIT}"