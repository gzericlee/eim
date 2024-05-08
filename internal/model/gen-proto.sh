#!/bin/sh

dir=$(dirname "$0")

name=$1

protoc --go_out=${dir} ${dir}/${name}.proto
protoc --go-grpc_out=${dir} ${dir}/${name}.proto