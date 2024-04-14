#!/usr/bin/env bash
#
# post a slack message with curl
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

#MSG="Hello, World at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
MSG=$(cat "${SCRIPT_HOME}/block_sectionimage.json")

echo '{}' | jq --arg CHANNEL "${SLACK_CHANNEL}" --arg MSG "${MSG}" --compact-output '.channel|=$CHANNEL|.blocks|=$MSG' | cb

echo "INFO: json is $(cb)"



curl \
    --data "$(cb)" \
    --header "Authorization: Bearer $SLACK_BOT_TOKEN" \
    --header "Content-Type: application/json; charset=UTF-8" \
    --request POST \
    "https://slack.com/api/chat.postMessage"
 