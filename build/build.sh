#!/bin/bash

cd ${0%/*}/..

rm -rf build/dist
mkdir build/dist

VERSION=$(cat VERSION)
COMMIT=$(git rev-parse --verify HEAD)
BRANCH=`eval echo $(git branch -r --contains ${COMMIT})`
NOW=$(TZ=UTC-8 date '+%Y-%m-%d %H:%M:%S')

echo "Version: ${VERSION}，Branch: ${BRANCH}，Commit: ${COMMIT}，Date: ${NOW}"

export CGO_ENABLED=0

echo "Compiling eim_api service..."
go build -o build/dist/eim_api -tags netgo -ldflags \
"-s -w
-X 'github.com/gzericlee/eim.ServiceName=EIM-Api'
-X 'github.com/gzericlee/eim.Version=${VERSION:-dev}'
-X 'github.com/gzericlee/eim.Branch=${BRANCH:-master}'
-X 'github.com/gzericlee/eim.Commit=${COMMIT:-dev}'
-X 'github.com/gzericlee/eim.Date=${NOW}'" \
./cmd/api

echo "Compiling eim_file_flex service..."
go build -o build/dist/eim_file_flex -tags netgo -ldflags \
"-s -w
-X 'github.com/gzericlee/eim.ServiceName=EIM-FileFlex'
-X 'github.com/gzericlee/eim.Version=${VERSION:-dev}'
-X 'github.com/gzericlee/eim.Branch=${BRANCH:-master}'
-X 'github.com/gzericlee/eim.Commit=${COMMIT:-dev}'
-X 'github.com/gzericlee/eim.Date=${NOW}'" \
./cmd/fileflex

echo "Compiling eim_gateway service..."
go build -o build/dist/eim_gateway -tags netgo -ldflags \
"-s -w
-X 'github.com/gzericlee/eim.ServiceName=EIM-Gateway'
-X 'github.com/gzericlee/eim.Version=${VERSION:-dev}'
-X 'github.com/gzericlee/eim.Branch=${BRANCH:-master}'
-X 'github.com/gzericlee/eim.Commit=${COMMIT:-dev}'
-X 'github.com/gzericlee/eim.Date=${NOW}'" \
./cmd/gateway

echo "Compiling eim_storage service..."
go build -o build/dist/eim_storage -tags netgo -ldflags \
"-s -w
-X 'github.com/gzericlee/eim.ServiceName=EIM-Storage'
-X 'github.com/gzericlee/eim.Version=${VERSION:-dev}'
-X 'github.com/gzericlee/eim.Branch=${BRANCH:-master}'
-X 'github.com/gzericlee/eim.Commit=${COMMIT:-dev}'
-X 'github.com/gzericlee/eim.Date=${NOW}'" \
./cmd/storage

echo "Compiling eim_dispatch service..."
go build -o build/dist/eim_dispatch -tags netgo -ldflags \
"-s -w
-X 'github.com/gzericlee/eim.ServiceName=EIM-Dispatch'
-X 'github.com/gzericlee/eim.Version=${VERSION:-dev}'
-X 'github.com/gzericlee/eim.Branch=${BRANCH:-master}'
-X 'github.com/gzericlee/eim.Commit=${COMMIT:-dev}'
-X 'github.com/gzericlee/eim.Date=${NOW}'" \
./cmd/dispatch

echo "Compiling eim_seq service..."
go build -o build/dist/eim_seq -tags netgo -ldflags \
"-s -w
-X 'github.com/gzericlee/eim.ServiceName=EIM-Seq'
-X 'github.com/gzericlee/eim.Version=${VERSION:-dev}'
-X 'github.com/gzericlee/eim.Branch=${BRANCH:-master}'
-X 'github.com/gzericlee/eim.Commit=${COMMIT:-dev}'
-X 'github.com/gzericlee/eim.Date=${NOW}'" \
./cmd/seq

echo "Compiling eim_auth service..."
go build -o build/dist/eim_auth -tags netgo -ldflags \
"-s -w
-X 'github.com/gzericlee/eim.ServiceName=EIM-Auth'
-X 'github.com/gzericlee/eim.Version=${VERSION:-dev}'
-X 'github.com/gzericlee/eim.Branch=${BRANCH:-master}'
-X 'github.com/gzericlee/eim.Commit=${COMMIT:-dev}'
-X 'github.com/gzericlee/eim.Date=${NOW}'" \
./cmd/auth

echo "Compiling eim_mock service..."
go build -o build/dist/eim_mock -tags netgo -ldflags \
"-s -w
-X 'github.com/gzericlee/eim.ServiceName=EIM-Mock'
-X 'github.com/gzericlee/eim.Version=${VERSION:-dev}'
-X 'github.com/gzericlee/eim.Branch=${BRANCH:-master}'
-X 'github.com/gzericlee/eim.Commit=${COMMIT:-dev}'
-X 'github.com/gzericlee/eim.Date=${NOW}'" \
./cmd/mock

echo "Compiled..."

if [[ -n $1 ]] && [[ $1 == 'build_images' ]]; then

  echo "Building images..."

  docker build -t eim/mock -f build/Dockerfile_mock build
  docker build -t eim/gateway -f build/Dockerfile_gateway build
  docker build -t eim/seq -f build/Dockerfile_seq build
  docker build -t eim/auth -f build/Dockerfile_auth build
  docker build -t eim/api -f build/Dockerfile_auth build
  docker build -t eim/dispatch -f build/Dockerfile_dispatch build
  docker build -t eim/file_flex -f build/Dockerfile_fileflex build
  docker build -t eim/storage -f build/Dockerfile_storage build

  echo "Compiled..."

fi