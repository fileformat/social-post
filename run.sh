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
#social-post email --subject="testing on $(date -u)" --to=fileformat@gmail.com "custom ua!"
social-post facebook --image=bin/shields.png --image-caption="$(date -u)" "this is the shields!"
#go test -timeout 30s -run "^TestBadger$" github.com/FileFormatInfo/fflint/cmd/fflint