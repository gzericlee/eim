#!/bin/bash

cd ${0%/*}/..

rm -rf build/dist
mkdir build/dist

VERSION=$(cat VERSION)
COMMIT=$(git rev-parse --verify HEAD)
BRANCH=`eval echo $(git branch -r --contains ${COMMIT})`
NOW=$(TZ=UTC-8 date '+%Y-%m-%d %H:%M:%S')

echo "Version: ${VERSION}，Branch: ${BRANCH}，Commit: ${COMMIT}，Date: ${NOW}"

echo "Compiling eim_api service..."
go build -o build/dist/eim_api -tags netgo -ldflags \
"-s -w
-X 'eim.ServiceName=EIM-Api'
-X 'eim.Version=${VERSION:-dev}'
-X 'eim.Branch=${BRANCH:-master}'
-X 'eim.Commit=${COMMIT:-dev}'
-X 'eim.Date=${NOW}'" \
./cmd/api

echo "Compiling eim_gateway service..."
go build -o build/dist/eim_gateway -tags netgo -ldflags \
"-s -w
-X 'eim.ServiceName=EIM-Gateway'
-X 'eim.Version=${VERSION:-dev}'
-X 'eim.Branch=${BRANCH:-master}'
-X 'eim.Commit=${COMMIT:-dev}'
-X 'eim.Date=${NOW}'" \
./cmd/gateway

echo "Compiling eim_storage service..."
go build -o build/dist/eim_storage -tags netgo -ldflags \
"-s -w
-X 'eim.ServiceName=EIM-Storage'
-X 'eim.Version=${VERSION:-dev}'
-X 'eim.Branch=${BRANCH:-master}'
-X 'eim.Commit=${COMMIT:-dev}'
-X 'eim.Date=${NOW}'" \
./cmd/storage

echo "Compiling eim_dispatch service..."
go build -o build/dist/eim_dispatch -tags netgo -ldflags \
"-s -w
-X 'eim.ServiceName=EIM-Dispatch'
-X 'eim.Version=${VERSION:-dev}'
-X 'eim.Branch=${BRANCH:-master}'
-X 'eim.Commit=${COMMIT:-dev}'
-X 'eim.Date=${NOW}'" \
./cmd/dispatch

echo "Compiling eim_seq service..."
go build -o build/dist/eim_seq -tags netgo -ldflags \
"-s -w
-X 'eim.ServiceName=EIM-Seq'
-X 'eim.Version=${VERSION:-dev}'
-X 'eim.Branch=${BRANCH:-master}'
-X 'eim.Commit=${COMMIT:-dev}'
-X 'eim.Date=${NOW}'" \
./cmd/seq

echo "Compiling eim_auth service..."
go build -o build/dist/eim_auth -tags netgo -ldflags \
"-s -w
-X 'eim.ServiceName=EIM-Auth'
-X 'eim.Version=${VERSION:-dev}'
-X 'eim.Branch=${BRANCH:-master}'
-X 'eim.Commit=${COMMIT:-dev}'
-X 'eim.Date=${NOW}'" \
./cmd/auth

echo "Compiling eim_mock service..."
go build -o build/dist/eim_mock -tags netgo -ldflags \
"-s -w
-X 'eim.ServiceName=EIM-Mock'
-X 'eim.Version=${VERSION:-dev}'
-X 'eim.Branch=${BRANCH:-master}'
-X 'eim.Commit=${COMMIT:-dev}'
-X 'eim.Date=${NOW}'" \
./cmd/mock

echo "Compiled..."

if [[ -n $1 ]] && [[ $1 == 'build_images' ]]; then

  echo "Building images..."

  docker build -t eim/mock -f Dockerfile_mock ./dist
  docker build -t eim/gateway -f Dockerfile_gateway ./dist
  docker build -t eim/seq -f Dockerfile_seq ./dist
  docker build -t eim/auth -f Dockerfile_auth ./dist
  docker build -t eim/dispatch -f Dockerfile_dispatch ./dist
  docker build -t eim/storage -f Dockerfile_storage ./dist

  echo "Compiled..."

fi