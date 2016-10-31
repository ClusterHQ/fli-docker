#!/bin/bash

BUCKET_NAME="clusterhq-fli-docker"
VERSION="0.0.2-dev"
UPLOAD_ON_BRANCH="s3"

function PreflightUbuntu
{
  pip install awscli --user travis
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

  if ! $TRAVIS_PULL_REQUEST && [ $TRAVIS_BRANCH == "$UPLOAD_ON_BRANCH" ]; then
      # make sure UPLOAD_ON_BRANCH and VERSION are set above.
	  UploadToS3
  else
	  echo "Skipping push of version $VERSION, not master branch"
  fi
}

Main