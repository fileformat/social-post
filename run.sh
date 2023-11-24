#!/usr/bin/env bash
#
# run locally
#

set -o errexit
set -o pipefail
set -o nounset



rm -rf ./dist/social-post
mkdir -p ./dist
echo "INFO: building social-post"
go build -o ./dist/social-post ./main.go

if [ -f ".env" ]; then
    echo "INFO: loading .env"
    export $(cat .env)
fi
export PATH=$PATH:$(pwd)/dist
social-post email --subject="testing on $(date -u)" --to=fileformat@gmail.com "this is the body!"
#go test -timeout 30s -run "^TestBadger$" github.com/FileFormatInfo/fflint/cmd/fflint