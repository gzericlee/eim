#!/bin/bash

dir=$(dirname "$0")

source "${dir}"/env

if [ "$1" = 'all' ] || [ "$1" = '' ]; then
  #docker compose -p eim down --remove-orphans
  docker compose -p etcd down --remove-orphans
  docker compose -p nats down --remove-orphans
  docker compose -p redis down --remove-orphans
  docker compose -p mongodb down --remove-orphans
else
  docker compose -p "$1" down --remove-orphans
fi