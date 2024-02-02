GIT_COMMIT=$(shell git describe --always)

.PHONY: all build clean test install-venv update-venv

default: build

all: build test

build:
	go build -ldflags "-X github.com/asmaloney/gactar/util/version.BuildVersion=${GIT_COMMIT}"

clean:
	go clean

test:
	go test ./...

####
# The following are convenience targets for managing the virtual environment

# Create the gactar venv (with dev packages)
install-venv:
	./gactar env setup --dev

# Update the python version in the venv and
# any versions of packages from the requirements files
update-venv:
	./gactar env update --all --dev
