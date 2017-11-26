include .env

export

# TODO:
# quote the SERVERS variable in .env and split on spaces in the Makefile

# single service
ifeq ($(filter $(firstword $(MAKECMDGOALS)),"run" "config"),)
	RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  	$(eval $(RUN_ARGS):;@:)
endif


# Build the project
all: clean init config check build test

init: kubernetes monitoring network

configure:
	@./tools/grpc-generator/generate.sh --service $(RUN_ARGS)

all_config:
	@$(foreach VAR,$(SERVERS), \
	./tools/grpc-generator/generate.sh --service $(VAR);)

check:
	@./scripts/go_checker.sh

build:
	@./scripts/build.sh --service $(RUN_ARGS)

run: configure build

run_all: all_config
	$(foreach VAR,$(SERVERS), \
	./scripts/build.sh --service $(VAR);)
	@kubectl get pods --all-namespaces

test:

clean:

kubernetes:
	@./scripts/setup_kubernetes.sh

minikube:
	@./scripts/setup_minikube.sh

monitoring: prometheus grafana

prometheus:
	@./scripts/setup_prometheus.sh

grafana:
	@./scripts/setup_grafana.sh

network:
	@./scripts/setup_network.sh


.PHONY: run init configure check build test clean minikube monitoring \
	grafana prometheus test kubernetes run run_all all_config
