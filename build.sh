#!/usr/bin/env bash

set -x
set -e

BINDIR=$(pwd)/bin
mkdir -p $BINDIR
for tool in tools/*
do
name=$(basename $tool)
echo "building $name"
pushd $tool
go build \
    -ldflags="-X 'main.Version=${VERSION}' -X 'main.Revision=${VCS_REF}' -X 'main.GoVersion=go${GOVERSION}' -X 'main.Built=${BUILD_DATE}' -X 'main.OsArch=${GOOS}/${GOARCH}'" \
    -mod=vendor \
    -o ${BINDIR}/$name
popd
done

