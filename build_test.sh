#!/bin/bash

apt-get update -qq -y
apt-get install uuid-runtime

go test -v ./...