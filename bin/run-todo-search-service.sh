#!/usr/bin/env bash

REPO="benniekrijger"
IMAGE="todo-search-service-go"
VERSION="latest"

: ${HOST:=$(ipconfig getifaddr en0)}
: ${HOST:=$(ipconfig getifaddr en1)}
: ${HOST:=$(ipconfig getifaddr en2)}
: ${HOST:=$(ipconfig getifaddr en3)}
: ${HOST:=$(ipconfig getifaddr en4)}

echo "Stopping old instances with name: ${IMAGE}..."
docker ps -aq --filter name=${IMAGE} | xargs docker stop

echo "Removing old instances with name: ${IMAGE}..."
docker ps -aq --filter name=${IMAGE} | xargs docker rm -f

echo "Running docker image ${IMAGE} with tag ${VERSION}, reachable at 127.0.0.1:8012"

docker run \
  -d \
  --name ${IMAGE} \
  --publish 8012:8080 \
  -e ES_URL=http://$HOST:9200 \
  -e NATS_URL=nats://$HOST:4222 \
  ${REPO}/${IMAGE}:${VERSION}

