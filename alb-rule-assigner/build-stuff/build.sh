#!/bin/bash

cd "$(dirname "$0")/.."

BUILD_BUILDNUMBER=temp

DOCKER_REPO="listener-stuff"

IMAGE_TAG="$DOCKER_REPO":"$BUILD_BUILDNUMBER"

docker build -t $IMAGE_TAG . 

id=$(docker create "$IMAGE_TAG" )

mkdir -p artifacts

docker cp "$id":/app/main.zip ./artifacts
