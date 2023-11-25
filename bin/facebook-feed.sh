#!/usr/bin/env bash
#
#
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

curl --request GET \
    "https://graph.facebook.com/v18.0/me/photos?fields=alt_text_custom%2Calt_text%2Cid%2Cfrom%2Cpage_story_id&access_token=${FACEBOOK_PAGE_ACCESS_TOKEN}" \
    | cb

cb | jq .
