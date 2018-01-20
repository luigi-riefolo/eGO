include conf/develop.env
export

.PHONY: run init config check build test clean minikube monitoring \
	grafana prometheus test kubernetes run run_all all_config registry

# TODO:
# make include conditional on type of environment, but don't know hot to
# set/pass the type of environment. Which options do we have?
#

BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
COMMIT_HASH=$(shell git rev-parse --short HEAD)


SERVERS=$(shell ls src)

# single service
ifeq ($(filter $(firstword $(MAKECMDGOALS)),"run" "config"),)
	RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  	$(eval $(RUN_ARGS):;@:)
endif


# Build the project
all: clean grpc_service check build test

PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
GO_BINDATA := $(GOPATH)/bin/go-bindata


############
#    Go    #
############

## Install the dependencies
dep:
	$(info installing Go dependencies)
	@go get -v -d ./...

## Generate a gRPC service
grpc_service:
	$(info generating gRPC service $(RUN_ARGS))
	#@./tools/grpc-generator/generate.sh --service $(RUN_ARGS)

## Generate all the gRPC services
all_config:
	@$(foreach VAR,$(SERVERS), \
	./tools/grpc-generator/generate.sh --service $(VAR);)

## Check all Go code
code_check:
	$(info checking Go code)
	@./scripts/go_checker.sh

## Create a Docker image for a service
build:
	@./scripts/build.sh --service $(RUN_ARGS)

run: config build

run_all: all_config
	$(foreach VAR,$(SERVERS), \
	./scripts/build.sh --service $(VAR);)
	@kubectl get pods --all-namespaces

## Run memory sanitizer
msan: dep
	@go test -msan -short ${PKG_LIST}

## Test the Go code
test:

## Generate global code coverage report in HTML
coverhtml:
	#@./tools/coverage.sh html;

clean:

## Display this help screen
help:
	@gawk 'match($$0, /^## (.*)/, a) \
	{getline x; printf "\033[36m%-30s\033[0m %s\n", x, a[1];}' $(MAKEFILE_LIST)


#########################
#    Infrastructure     #
#########################

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



# INSTALL PROTOC
PROTOC := $(shell which protoc)
# If protoc isn't on the path, set it to a target that's never up to date, so
# the install command always runs.
ifeq ($(PROTOC),)
    PROTOC = must-rebuild
endif

# Figure out which machine we're running on.
UNAME := $(shell uname)

$(PROTOC):
# Run the right installation command for the operating system.
ifeq ($(UNAME), Darwin)
	brew install protobuf
endif
ifeq ($(UNAME), Linux)
	sudo apt-get install protobuf-compiler
endif
# You can add instructions for other operating systems here, or use different
# branching logic as appropriate.

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install
# it.
$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go


$(GO_BINDATA):
	go get github.com/kevinburke/go-bindata

#
#
#assets: $(shell find static) | $(GO_BINDATA)
#	$(GO_BINDATA) -o=assets/bindata.go --nocompress --nometadata --pkg=assets static/...
