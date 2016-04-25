#!/bin/bash -x
rm -rf target
go get -t -v ./...
go test -v ./...
go build

mkdir target -p
cp -r docker target/
cp -r html target/docker/
cp ccTray2Slack target/docker/

cd target/docker/
docker-compose build
