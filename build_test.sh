#!/bin/bash

go test -v ./...
go build -o build/fli-docker $HOME/gopath/src/github.com/ClusterHQ/fli-docker/cmd/fli-docker