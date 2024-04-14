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

FILENAME="$(realpath "${BASE_DIR}/docs/favicon512.png")"
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
    "${UPLOAD_URL}"

curl \
    --data "{\"files\": [{\"id\":\"${FILE_ID}\"}]}" \
    --header "Authorization: Bearer $SLACK_BOT_TOKEN" \
    --header "Content-Type: application/json; charset=UTF-8" \
    --request POST \
    "https://slack.com/api/files.completeUploadExternal"
