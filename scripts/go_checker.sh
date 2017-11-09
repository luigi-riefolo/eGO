#!/usr/bin/env bash

# Script for Go code error checking.

set -xe

DIR="./..."
FOLDERS="pkg src"
SUB_FOLDERS="pkg/... src/..."

PREV_DIR="$(pwd)"
PROJECT_DIR=${GOPATH}/src/github.com/luigi-riefolo/eGO

cd $PROJECT_DIR

goimports -w $FOLDERS

errcheck $DIR
structcheck $DIR
varcheck $DIR
staticcheck $DIR
goconst -ignore vendor $FOLDERS

go vet $DIR
gosimple $DIR
golint -set_exit_status $SUB_FOLDERS

cd $PREV_DIR

echo "Go checker completed"
