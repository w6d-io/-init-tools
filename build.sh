#!/bin/bash -x

for tool in tools/*
do
name=$(basename $tool)
echo "building $name"
pushd $tool
go build \
    -ldflags="-X 'main.Version=${VERSION}' -X 'main.Revision=${VCS_REF}' -X 'main.GoVersion=go${GOVERSION}' -X 'main.Built=${BUILD_DATE}' -X 'main.OsArch=${GOOS}/${GOARCH}'" \
    -mod=vendor \
    -o /build/bin/$name
popd
done

