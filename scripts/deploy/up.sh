#!/bin/bash

source ./env

docker-compose -p eim -f eim-docker-compose/docker-compose.yaml up -d
docker-compose -p tidb -f tidb-docker-compose/docker-compose.yml up -d
docker-compose -p redis -f emb-docker-compose/redis-cluster.yaml up -d
docker-compose -p nsq -f emb-docker-compose/nsq-cluster.yaml up -d