#!/bin/bash
GOOS=${GOOS:-$(go env GOOS)}
GOARCH=${GOARCH:-$(go env GOARCH)}
./$GOOS-$GOARCH/chartmuseum --debug --port=8443 \
  --storage="local" \
  --storage-local-rootdir="./chartstore-8443" \
  --tls-cert=./server.crt --tls-key=./server.key 
