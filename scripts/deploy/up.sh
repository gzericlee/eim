#!/bin/bash

source ./env

docker-compose -p eim -f eim-docker-compose/docker-compose.yaml up -d
docker-compose -p tidb -f tidb-docker-compose/docker-compose.yml up -d
docker-compose -p redis -f eim-docker-compose/redis-cluster.yaml up -d
docker-compose -p nsq -f eim-docker-compose/nsq-cluster.yaml up -d