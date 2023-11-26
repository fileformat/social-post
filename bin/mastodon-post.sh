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
    --request POST \
    --form "status=Testing at $(date -u)" \
    --form "visibility=public" \
	--header "Authorization: Bearer ${MASTODON_USER_TOKEN}" \
	https://${MASTODON_SERVER}/api/v1/statuses \
    | cb

cb | cat

cb | jq .
