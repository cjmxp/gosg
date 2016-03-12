#!/bin/bash

files=`ls *.proto`

protoc \
 -I/usr/local/include \
 -I. \
 -I$GOPATH/src \
 -I$GOPATH/src/github.com/gengo/grpc-gateway/third_party/googleapis \
 --go_out=:. \
 $files
