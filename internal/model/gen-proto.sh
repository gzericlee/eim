#!/bin/sh

dir=$(dirname "$0")

protoc --go_out=${dir} ${dir}/*.proto
protoc --go-grpc_out=${dir} ${dir}/*.proto