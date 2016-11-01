#! /usr/bin/env bash

rm -rf ./vendor
rm glide.lock
glide up

grep -R --exclude-dir vendor --exclude-dir .git --exclude build.sh TODO ./

GOOS=linux GOARCH=amd64 go build .
