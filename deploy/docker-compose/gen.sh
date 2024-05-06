#!/bin/bash

dir=$(dirname "$0")

source "${dir}/env"

for no in $(seq 1 6); do \
  NO=${no} envsubst < "$dir/redis/cluster.tmpl" > "$dir/redis/redis-${no}.conf"
done