
IMG ?= w6dio/init-tools:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

REF=$(shell git symbolic-ref --quiet HEAD 2> /dev/null)
VERSION=$(shell basename $(REF) )
VCS_REF=$(shell git rev-parse HEAD)
GOVERSION=$(shell go version | awk '{ print $3 }' | sed 's/go//')
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GOOS=$(shell uname -s | tr "[:upper:]" "[:lower:]")
GOARCH=$(shell uname -p)

all: init-tools

# Build ci-status binary
init-tools:
	@./build.sh

# Build the docker image
build:
	docker build  --build-arg=VERSION=${VERSION} --build-arg=VCS_REF=${VCS_REF} --build-arg=BUILD_DATE=${BUILD_DATE}  -t ${IMG} .

# Push the docker image
push:
	docker push ${IMG}

