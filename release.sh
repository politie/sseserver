#!/usr/bin/env bash
set -e

CWD="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [ "$#" -ne 1 ]; then
    echo "Please provide a release version number"
    exit 1;
fi

VERSION=$1
BUILD_OS="darwin linux windows"
# BUILD_OS="darwin linux"

function build_binaries {

    # save current GOOS and GOARCH
    CURRENT_GOOS=${GOOS}
    CURRENT_GOARCH=${GOARCH}

    echo "Building binaries for version ${VERSION}"
    for GOOS in ${BUILD_OS}; do
        for GOARCH in 386 amd64; do
            FILE=sseserver-${VERSION}-${GOOS}-${GOARCH}
            if [ -f ${FILE} ]; then
                rm ${FILE}
            fi
            echo "Building ${FILE}"
            export GOOS=${GOOS}
            export GOARCH=${GOARCH}
            go build -o ${CWD}/bin/${FILE}
            if [ "${GOOS}" == "windows" ]; then
               mv ${CWD}/bin/${FILE} ${CWD}/bin/${FILE}.exe
            else
               chmod +x ${CWD}/bin/${FILE}
            fi
        done
    done

    # reset GOOS and GOARCH
    export GOOS=${CURRENT_GOOS}
    export GOARCH=${CURRENT_GOARCH}
    echo "Done"
}

function test_app {
  go test -v
}

function main {
    build_binaries
}

main
