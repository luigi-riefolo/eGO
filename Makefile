# Credits:
# https://gist.github.com/turtlemonvh/38bd3d73e61769767c35931d8c70ccb4

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
	#-rm -f ${TEST_REPORT}
	#-rm -f ${VET_REPORT}
	#-rm -f ${BINARY}-*

.PHONY: init config check build test clean
