#!/bin/bash

go test -v ./...
go build -o build/fli-docker cmd/fli-docker