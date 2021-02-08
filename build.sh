#!/bin/sh
set -e
set -x

make clean
make gen
make test
make build
docker build . -t $IMAGE
