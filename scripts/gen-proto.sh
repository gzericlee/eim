#!/bin/sh

cd internal/model/proto

name=$1

protoc --go_out=../ *.proto

protoc-go-inject-tag -input=../*.pb.go