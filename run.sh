#!/usr/bin/env bash
#
# run locally
#

set -o errexit
set -o pipefail
set -o nounset

rm -rf ./dist/social-post
mkdir -p ./dist
go build -o ./dist/social-post ./main.go
export PATH=$PATH:$(pwd)/dist
social-post version
#go test -timeout 30s -run "^TestBadger$" github.com/FileFormatInfo/fflint/cmd/fflint