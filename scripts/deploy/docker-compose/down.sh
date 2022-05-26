#!/bin/bash

source ./env

if [ $1 == 'all' ]; then
  #docker-compose -p eim down --remove-orphans
  docker-compose -p etcd down --remove-orphans
  docker-compose -p nsq down --remove-orphans
  docker-compose -p redis down --remove-orphans
  docker-compose -p tidb down --remove-orphans
else
  docker-compose -p $1 down --remove-orphans
fi