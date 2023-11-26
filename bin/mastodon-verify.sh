#!/usr/bin/env bash
#
# verify mastodon creds
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
	--header "Authorization: Bearer ${MASTODON_ACCESS_TOKEN}" \
	https://${MASTODON_SERVER}/api/v1/apps/verify_credentials \
    | cb

cb | cat

cb | jq .
