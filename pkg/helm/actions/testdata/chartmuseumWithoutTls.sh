#!/bin/bash
GOOS=${GOOS:-$(go env GOOS)}
GOARCH=${GOARCH:-$(go env GOARCH)}
./$GOOS-$GOARCH/chartmuseum --debug --port=8080 \
  --storage="local" \
  --storage-local-rootdir="./chartstore-8080"

