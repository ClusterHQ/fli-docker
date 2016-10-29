#!/bin/bash

PACKAGE_NAME="fli-docker"
VERSION="0.0.0-dev"
FILE_TARGET_PATH="fli-docker"

go test -v ./...
go build -v ./...
go install -v ./...
mv $GOPATH/bin/fli-docker build/

if ! $TRAVIS_PULL_REQUEST
	then
		curl -T build/fli-docker -u$BINTRAY_USER:$BINTRAY_API_KEY https://api.bintray.com/content/chqtest/fli-docker/$PACKAGE_NAME/$VERSION/$FILE_TARGET_PATH
fi