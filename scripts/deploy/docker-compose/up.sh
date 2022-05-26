#!/bin/bash

source ./env

if [ $1 == 'all' ]; then
  #docker-compose -p eim -f eim/docker-compose.yml up -d
  docker-compose -p tidb -f tidb/docker-compose.yml up -d
  docker-compose -p redis -f redis/docker-compose.yml up -d
  docker-compose -p nsq -f nsq/docker-compose.yml up -d
  docker-compose -p etcd -f etcd/docker-compose.yml up -d
else
  docker-compose -p $1 -f $1/docker-compose.yml up -d
fi