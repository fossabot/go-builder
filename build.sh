#!/bin/bash
GIT_COMMIT=$(git rev-list -1 HEAD) 
VERSION=$(git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///')

if [ -z ${VERSION}]; then
    VERSION="unknown"
fi

echo ${VERSION}

go run cmd/builder/main.go -- cmd/builder/main.go -git ${GIT_COMMIT} -version ${VERSION} -upx 9