#!/bin/bash

source ../env

for no in `seq 1 6`; do \
  mkdir -p config/redis/redis-${no} \
  && NO=${no} envsubst < redis-cluster.tmpl > config/redis/redis-${no}/redis.conf \
  && mkdir -p data/redis/redis-${no};\
done

for no in `seq 1 3`; do \
  mkdir -p data/nsq/nsqd-${no}
done