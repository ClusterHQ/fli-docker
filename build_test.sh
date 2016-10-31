#!/bin/bash

BUCKET_NAME="clusterhq-fli-docker"
VERSION="0.0.2-dev"

function PreflightUbuntu
{
  apt-get update -y
  apt-get install -y python3
  pip3 install awscli

}

function SetupAWSCredentials
{
  if ! test -d ${HOME}/.aws; then
    mkdir ${HOME}/.aws
  fi
  printf "$(cat .aws-cred-template)" "$AWS_SECRET_ID" "$AWS_ACCESS_KEY" > ${HOME}/.aws/credentials
}

### Upload HTML to Amazon S3
function UploadToS3
{
  SetupAWSCredentials
  aws s3 sync build/bin/ s3://$BUCKET_NAME/$VERSION
}


function BuildAndTest
{
  go test -v ./...
  go build -v ./...
  go install -v ./...
  mkdir build/bin
  mv $GOPATH/bin/fli-docker build/bin/
}

function Main
{
  BuildAndTest

  if [[ "$OSTYPE" == "linux-gnu" ]]; then
    PreflightUbuntu
  else
    echo "Unrecognized operating system"
    exit 1
  fi

  if ! $TRAVIS_PULL_REQUEST && [ $TRAVIS_BRANCH == "s3" ]; then
	  UploadToS3
  else
	  echo "Skipping push of version $VERSION, not master branch"
  fi

Main