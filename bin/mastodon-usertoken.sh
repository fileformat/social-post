#!/usr/bin/env bash
#
# get a user token needed for posting
#


set -o errexit
set -o pipefail
set -o nounset

SCRIPT_HOME="$( cd "$( dirname "$0" )" && pwd )"
BASE_DIR=$(realpath "${SCRIPT_HOME}/..")

#
# load an .env file if it exists
#
ENV_FILE="${BASE_DIR}/.env"
if [ -f "${ENV_FILE}" ]; then
    echo "INFO: loading '${ENV_FILE}'"
    export $(cat "${ENV_FILE}")
fi

echo "INFO: opening auth URL"
open \
    "https://${MASTODON_SERVER}/oauth/authorize?client_id=${MASTODON_CLIENT_KEY}&scope=read+write&redirect_uri=urn:ietf:wg:oauth:2.0:oob&response_type=code"
    

read -p "Enter the code: " CODE

echo "INFO: requesting access token"
curl \
    --request POST \
    --data "client_id=${MASTODON_CLIENT_KEY}" \
    --data "client_secret=${MASTODON_CLIENT_SECRET}" \
    --data "grant_type=authorization_code" \
    --data "code=${CODE}" \
    --data "redirect_uri=urn:ietf:wg:oauth:2.0:oob" \
    https://${MASTODON_SERVER}/oauth/token \
    | cb

cb | cat

cb | jq .