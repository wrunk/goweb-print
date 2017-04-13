#!/bin/sh


export GOBIN=$GOPATH/src/github.com/wrunk/goweb-print

cd "$(dirname "$0")"

go install -v -ldflags=-s && ./goweb-print
