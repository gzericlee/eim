#!/bin/bash

dir=$(dirname "$0")

source "${dir}"/env

mkdir -p "${EIM_CONFIG_DIR}"
mkdir -p "${EIM_DATA_DIR}"
mkdir -p "${EIM_DATA_DIR}"/tidb
mkdir -p "${EIM_DATA_DIR}"/tidb/logs

for no in $(seq 1 6); do \
  mkdir -p "${EIM_CONFIG_DIR}"/redis/redis-"${no}" \
  && NO=${no} envsubst < "${dir}"/redis/cluster.tmpl > "${EIM_CONFIG_DIR}"/redis/redis-"${no}"/redis.conf \
  && mkdir -p "${EIM_DATA_DIR}"/redis/redis-"${no}";\
done

for no in $(seq 1 3); do \
  mkdir -p "${EIM_DATA_DIR}"/nsq/nsqd-"${no}"
done

for no in $(seq 1 3); do \
  mkdir -p "${EIM_DATA_DIR}"/etcd/etcd-"${no}"
done