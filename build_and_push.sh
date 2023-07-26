#!/bin/sh

VERSION="0.1.0"
IMAGE_URL="registry.git.leon.wtf/leon/one-time-service/amd64"

docker build -t $IMAGE_URL:$VERSION . && docker push $IMAGE_URL:$VERSION