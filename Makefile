include .env

export


#VERSION=$(shell git rev-parse --abbrev-ref HEAD)
#git rev-parse HEAD
#git rev-parse --abbrev-ref HEAD

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

config:
	@./tools/grpc-generator/generate.sh --service $(RUN_ARGS)

all_config:
	@$(foreach VAR,$(SERVERS), \
	./tools/grpc-generator/generate.sh --service $(VAR);)

check:
	@./scripts/go_checker.sh

build:
	@./scripts/build.sh --service $(RUN_ARGS)

run: config build

run_all: all_config
	$(foreach VAR,$(SERVERS), \
	./scripts/build.sh --service $(VAR);)
	@kubectl get pods --all-namespaces

test:

clean:

registry:
	@./tools/registry/create_registry.sh

swagger:
	echo $(BASE)
	@./scripts/deploy.sh swagger

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


.PHONY: run init config check build test clean minikube monitoring \
	grafana prometheus test kubernetes run run_all all_config registry
