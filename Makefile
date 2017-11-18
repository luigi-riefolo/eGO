# TODO:
# support multiple services


# Build the project
all: clean init config check build test

init:
	./scripts/setup_kubernetes.sh

config:

check:
	./scripts/go_checker.sh

build:

test:

clean:

kubernetes:
	./scripts/setup_kubernetes.sh

prometheus:
	./scripts/setup_prometheus.sh



.PHONY: init config check build test clean
