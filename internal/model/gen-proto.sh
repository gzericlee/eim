#!/bin/sh

dir=$(dirname "$0")

cd ${dir}

name=$1

protoc --go_out=. ${name}.proto
protoc --go-grpc_out=. ${name}.proto