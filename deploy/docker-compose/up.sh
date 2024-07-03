#!/bin/bash

dir=$(dirname "$0")

source "${dir}"/env

for no in $(seq 1 6); do \
  NO=${no} envsubst < "$dir/redis/cluster.tmpl" > "$dir/redis/redis-${no}.conf"
done

if [ "$1" = 'all' ] || [ "$1" = '' ]; then
  #docker-compose -p eim -f ${dir}/eim/docker-compose.yml up -d
  docker compose -p redis -f "${dir}"/redis/docker-compose.yml up -d
  docker compose -p nats -f "${dir}"/nats/docker-compose.yml up -d
  docker compose -p etcd -f "${dir}"/etcd/docker-compose.yml up -d
  docker compose -p mongodb -f "${dir}"/mongodb/docker-compose.yml up -d
else
  docker compose -p "$1" -f "${dir}"/"$1"/docker-compose.yml up -d
fi