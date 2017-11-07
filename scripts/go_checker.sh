#!/usr/bin/env bash

# Script for Go code error checking.


# TODO: check if with many/big files
# is faster to run in background or parallelise the jobs

set -xe


CHECK_FOLDERS="./... ."
PREV_DIR="$(pwd)"
PROJECT_DIR=${GOPATH}/src/github.com/luigi-riefolo/eGO


cd $PROJECT_DIR

#gofmt -w

errcheck $CHECK_FOLDERS
structcheck $CHECK_FOLDERS
varcheck $CHECK_FOLDERS
staticcheck $CHECK_FOLDERS
goconst -ignore vendor $CHECK_FOLDERS

go vet ./...
gosimple $CHECK_FOLDERS
goimports -d $(find . -type f -name '*.go')
golint -set_exit_status ./...
govendor service +local

#tomlv some-toml-file.toml

#	go-bindata -o data/bindata.go -pkg data data/*.json
#go fmt $$(go list ./... | grep -v /vendor/) ;


cd $PREV_DIR
