#!/bin/bash
GOOS=${GOOS:-$(go env GOOS)}
GOARCH=${GOARCH:-$(go env GOARCH)}
rm -rf ./chartstore
rm -rf ./temporary
rm -rf ./$GOOS-$GOARCH
rm -rf ./chartmuseum.tar.gz