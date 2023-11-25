#!/usr/bin/env bash
#
# creates a single "photo-post" in one step
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

curl \
    --request POST \
    --form "source=@${FACEBOOK_TEST_IMAGE}" \
    --form "published=true" \
    --form "caption=test from facebook-postphotos.sh on $(date -u)" \
    --form "access_token=${FACEBOOK_PAGE_ACCESS_TOKEN}" \
    "https://graph.facebook.com/me/photos" \
    | cb

cb | cat

cb | jq .
