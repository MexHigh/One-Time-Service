#!/bin/bash

VERSION="${1:-dev}"
IMAGE_URL="registry.git.leon.wtf/leon/one-time-service/amd64"

echo "[*] Building image with tag: $IMAGE_URL:$VERSION"

docker build --build-arg ADDON_VERSION=$VERSION -t $IMAGE_URL:$VERSION .
docker push $IMAGE_URL:$VERSION