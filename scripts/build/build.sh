#!/bin/bash
BRANCH=$(git symbolic-ref --short -q HEAD)
COMMIT=$(git rev-parse --verify HEAD)
NOW=$(date '+%Y-%m-%dT%H:%M:%S')

go build -o dist/eim_gateway -tags netgo -ldflags \
"-s -w
-X eim/build.ServiceName=EIM-Gateway
-X eim/build.Branch=${BRANCH:-master}
-X eim/build.Commit=${COMMIT:-dev}
-X eim/build.Date=${NOW}" \
./cmd/gateway

go build -o dist/eim_storage -tags netgo -ldflags \
"-s -w
-X eim/build.ServiceName=EIM-Storage
-X eim/build.Branch=${BRANCH:-master}
-X eim/build.Commit=${COMMIT:-dev}
-X eim/build.Date=${NOW}" \
./cmd/storage

go build -o dist/eim_dispatch -tags netgo -ldflags \
"-s -w
-X eim/build.ServiceName=EIM-Dispatch
-X eim/build.Branch=${BRANCH:-master}
-X eim/build.Commit=${COMMIT:-dev}
-X eim/build.Date=${NOW}" \
./cmd/dispatch

go build -o dist/eim_seq -tags netgo -ldflags \
"-s -w
-X eim/build.ServiceName=EIM-Seq
-X eim/build.Branch=${BRANCH:-master}
-X eim/build.Commit=${COMMIT:-dev}
-X eim/build.Date=${NOW}" \
./cmd/seq

go build -o dist/eim_mock -tags netgo -ldflags \
"-s -w
-X eim/build.ServiceName=EIM-Mock
-X eim/build.Branch=${BRANCH:-master}
-X eim/build.Commit=${COMMIT:-dev}
-X eim/build.Date=${NOW}" \
./cmd/mock

docker build -t eim/mock -f Dockerfile_mock ./dist
docker build -t eim/gateway -f Dockerfile_gateway ./dist
docker build -t eim/seq -f Dockerfile_seq ./dist
docker build -t eim/dispatch -f Dockerfile_dispatch ./dist
docker build -t eim/storage -f Dockerfile_storage ./dist