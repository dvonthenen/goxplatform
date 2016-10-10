#! /usr/bin/env bash

grep -R --exclude-dir vendor --exclude-dir .git --exclude build.sh TODO ./

GOOS=linux GOARCH=amd64 go install
