#!/bin/bash
ENTRY="cmd/builder/main.go"
UPX="9"
GIT_COMMIT=$(git rev-list -1 HEAD) 
VERSION=$(git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///')

if [ -z ${VERSION}]; then
    VERSION="unknown"
fi

# If CGO is needed add -cgo
# ./go-builder ${ENTRY} -cgo -git ${GIT_COMMIT} -version ${VERSION} -upx ${UPX}
./go-builder ${ENTRY} -git ${GIT_COMMIT} -version ${VERSION} -upx ${UPX}
