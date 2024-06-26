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

#social-post email --subject="testing on $(date -u)" --to=fileformat@gmail.com "custom ua!"
#social-post facebook --image=${TEST_IMAGE} --image-caption="$(date -u)" "this is the shields!"
#social-post mastodon "from run.sh without image at $(date -u)!"
#social-post mastodon --image=${TEST_IMAGE} --image-caption="$(date -u)" "from run.sh at $(date -u)!"
#go test -timeout 30s -run "^TestBadger$" github.com/FileFormatInfo/fflint/cmd/fflint
./dist/social-post slack --channel ${SLACK_CHANNEL} --image=/Users/andrew/Downloads/vlz300.png "from *social-post-cli* at _$(date -u)_"
