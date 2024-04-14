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

FILENAME="$(realpath "${BASE_DIR}/docs/favicon512dark.png")"
BASENAME=$(basename "${FILENAME}")
echo "INFO: base name: '${BASENAME}'"

FILELEN=$(cat "${FILENAME}" | wc -c | tr -d ' ')

echo "INFO: file length: '${FILELEN}'"


#
# get the upload URL
#

curl \
    --header "Authorization: Bearer $SLACK_BOT_TOKEN" \
    --header "Content-Type: application/json; charset=UTF-8" \
    --request GET \
    "https://slack.com/api/files.getUploadURLExternal?filename=${BASENAME}&length=${FILELEN}" \
    | cb

RESULT=$(cb | jq -r '.ok')
if [ "${RESULT}" != "true" ]; then
    echo "ERROR: failed to get upload URL $(cb)"
    exit 1
fi
FILE_ID=$(cb | jq -r '.file_id')
echo "INFO: file ID: '${FILE_ID}'"

UPLOAD_URL=$(cb | jq -r '.upload_url')
echo "INFO: upload URL: '${UPLOAD_URL}'"

MIMETYPE=$(file --mime-type -b "${FILENAME}")
echo "INFO: mime type: '${MIMETYPE}'"

curl \
    --header "Content-Type: ${MIMETYPE}" \
    --data-binary "@${FILENAME}" \
    --request POST \
    "${UPLOAD_URL}" \
    | cb

echo "INFO: result from file upload '$(cb)'"

curl \
    --data "{\"files\": [{\"id\":\"${FILE_ID}\"}]}" \
    --header "Authorization: Bearer $SLACK_BOT_TOKEN" \
    --header "Content-Type: application/json; charset=UTF-8" \
    --request POST \
    "https://slack.com/api/files.completeUploadExternal" \
    | cb

echo "INFO: result from upload complete '$(cb)'"

#FILE_URL=$(cb | jq --raw-output '.files.[0].permalink_public')
#echo "INFO: file URL: '${FILE_URL}'"
FILE_ID=$(cb | jq --raw-output '.files.[0].id')
echo "INFO: file ID: '${FILE_ID}'"

echo "INFO: sleeping for a bit..."
sleep 1

MSG="Hello, **World** at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
BLOCK=$(cat "${SCRIPT_HOME}/block_sectionimageid.json" | jq --compact-output --arg FILE_ID "${FILE_ID}" --arg MSG "${MSG}" '.[0].text.text=$MSG|.[1].slack_file.id=$FILE_ID')
echo "INFO: block is ${BLOCK}"

echo '{}' | jq --arg CHANNEL "${SLACK_CHANNEL}" --arg BLOCK "${BLOCK}" --arg MSG "${MSG}" --compact-output '.channel|=$CHANNEL|.blocks|=$BLOCK|.text=$MSG' | cb

echo "INFO: json is $(cb)"

curl \
    --data "$(cb)" \
    --header "Authorization: Bearer $SLACK_BOT_TOKEN" \
    --header "Content-Type: application/json; charset=UTF-8" \
    --request POST \
    "https://slack.com/api/chat.postMessage"
 