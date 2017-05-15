#!/usr/bin/env bash

echo "Removing old scylla instances..."

docker ps -aq --filter name=scylla | xargs docker rm -f

echo "Running new scylla instance..."

docker run \
  -d \
  -p 9042:9042 \
  --name scylla \
  scylladb/scylla

echo "Waiting for scylla to be started..."
sleep 15

docker exec -i -t scylla sh -c 'cqlsh "CREATE KEYSPACE todos WITH replication = { class: \"SimpleStrategy\", replication_factor : 1};"'