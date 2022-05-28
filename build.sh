#!/bin/bash

BRANCH=$(git symbolic-ref --short -q HEAD)
COMMIT=$(git rev-parse --verify HEAD)
NOW=$(date '+%Y-%m-%dT%H:%M:%S')

echo "Branch: ${BRANCH}，Commit: ${COMMIT}，Date: ${NOW}"

echo "Compiling eim_gateway service..."
go build -o dist/eim_gateway -tags netgo -ldflags \
"-s -w
-X eim/build.ServiceName=EIM-Gateway
-X eim/build.Branch=${BRANCH:-master}
-X eim/build.Commit=${COMMIT:-dev}
-X eim/build.Date=${NOW}" \
./cmd/gateway

echo "Compiling eim_storage service..."
go build -o dist/eim_storage -tags netgo -ldflags \
"-s -w
-X eim/build.ServiceName=EIM-Storage
-X eim/build.Branch=${BRANCH:-master}
-X eim/build.Commit=${COMMIT:-dev}
-X eim/build.Date=${NOW}" \
./cmd/storage

echo "Compiling eim_dispatch service..."
go build -o dist/eim_dispatch -tags netgo -ldflags \
"-s -w
-X eim/build.ServiceName=EIM-Dispatch
-X eim/build.Branch=${BRANCH:-master}
-X eim/build.Commit=${COMMIT:-dev}
-X eim/build.Date=${NOW}" \
./cmd/dispatch

echo "Compiling eim_seq service..."
go build -o dist/eim_seq -tags netgo -ldflags \
"-s -w
-X eim/build.ServiceName=EIM-Seq
-X eim/build.Branch=${BRANCH:-master}
-X eim/build.Commit=${COMMIT:-dev}
-X eim/build.Date=${NOW}" \
./cmd/seq

echo "Compiling eim_auth service..."
go build -o dist/eim_auth -tags netgo -ldflags \
"-s -w
-X eim/build.ServiceName=EIM-Auth
-X eim/build.Branch=${BRANCH:-master}
-X eim/build.Commit=${COMMIT:-dev}
-X eim/build.Date=${NOW}" \
./cmd/auth

echo "Compiling eim_mock service..."
go build -o dist/eim_mock -tags netgo -ldflags \
"-s -w
-X eim/build.ServiceName=EIM-Mock
-X eim/build.Branch=${BRANCH:-master}
-X eim/build.Commit=${COMMIT:-dev}
-X eim/build.Date=${NOW}" \
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