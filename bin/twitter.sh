#!/usr/bin/env bash
#
# run locally
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
    --request GET \
    --header "Authorization: Bearer $TWITTER_BEARER_TOKEN" \
    'https://api.twitter.com/2/tweets/search/recent?query=from:twitterdev' 
    