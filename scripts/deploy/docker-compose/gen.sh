#!/bin/bash

source ./env

mkdir -p ${EIM_CONFIG_DIR}
mkdir -p ${EIM_DATA_DIR}

for no in `seq 1 6`; do \
  mkdir -p ${EIM_CONFIG_DIR}/redis/redis-${no} \
  && NO=${no} envsubst < redis/cluster.tmpl > ${EIM_CONFIG_DIR}/redis/redis-${no}/redis.conf \
  && mkdir -p ${EIM_DATA_DIR}/redis/redis-${no};\
done

for no in `seq 1 3`; do \
  mkdir -p ${EIM_DATA_DIR}/nsq/nsqd-${no}
done

for no in `seq 1 3`; do \
  mkdir -p ${EIM_DATA_DIR}/etcd/etcd-${no}
done