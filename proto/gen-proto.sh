#!/bin/sh

protoc --go_out=./proto/ ./proto/*.proto
protoc --go-grpc_out=./proto/ ./proto/*.proto