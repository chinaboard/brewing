#!/usr/bin/env bash

#set -x

GitStatus=`git status -s`
BuildTime=`date +'%Y.%m.%d.%H%M%S'`
BuildGoVersion=`go version`

LDFlags=" \
    -X 'github.com/chinaboard/brewing/pkg/bininfo.GitStatus=${GitStatus}' \
    -X 'github.com/chinaboard/brewing/pkg/bininfo.BuildTime=${BuildTime}' \
    -X 'github.com/chinaboard/brewing/pkg/bininfo.BuildGoVersion=${BuildGoVersion}' \
"

ROOT_DIR=`pwd`

if [ ! -d ${ROOT_DIR}/bin ]; then
  mkdir bin
fi

cd ${ROOT_DIR} && CGO_ENABLED=0 go build -ldflags "$LDFlags" -o ${ROOT_DIR}/bin/brewing-worker ${ROOT_DIR}/cmd/worker

ls -lrt ${ROOT_DIR}/bin &&
echo 'build done.'