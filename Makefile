# TODO:
# support multiple services


# Build the project
all: clean init config check build test

init:
	./scripts/setup_kubernetes.sh

config:
#tomlv some-toml-file.toml
#go-bindata -o data/bindata.go -pkg data data/*.json

check:
	./scripts/go_checker.sh

build:

test:

clean:

.PHONY: init config check build test clean
