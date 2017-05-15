#!/usr/bin/env bash

echo "Removing old elasticsearch instances..."

docker ps -aq --filter name=elasticsearch | xargs docker rm -f

echo "Running new elasticsearch instance..."

docker run \
  -d \
  -p 9200:9200 \
  -e "http.host=0.0.0.0" \
  -e "transport.host=0.0.0.0" \
  --name elasticsearch \
  elasticsearch