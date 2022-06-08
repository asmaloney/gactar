GIT_COMMIT=$(shell git describe --always)

all:
	go build -tags urfave_cli_no_docs -ldflags "-X github.com/asmaloney/gactar/util/version.BuildVersion=${GIT_COMMIT}"