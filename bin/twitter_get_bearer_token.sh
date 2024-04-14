#!/usr/bin/env bash
#
# get a bearer token from twitter (not working)
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
    -d grant_type=client_credentials \
    -u $TWITTER_KEY:$TWITTER_SECRET \
    https://api.twitter.com/oauth2/token
