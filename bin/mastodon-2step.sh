#!/usr/bin/env bash
#
# post a photo, then use it in a post
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
    --form "file=@${TEST_IMAGE}" \
    --form "description=Testing at $(date -u)" \
	--header "Authorization: Bearer ${MASTODON_USER_TOKEN}" \
	https://${MASTODON_SERVER}/api/v2/media \
    | cb

cb | cat

cb | jq .

PHOTO_ID=$(cb | jq --raw-output .id)

echo "INFO: photo id is '${PHOTO_ID}'"

#echo "INFO: sleeping 30 seconds to allow photo to be processed"
#sleep 30

echo "INFO: posting status with photo id '${PHOTO_ID}'"

curl \
    --verbose \
    --request POST \
    --form "status=Testing with photo media_ids[]=${PHOTO_ID} at $(date -u)" \
    --form "visibility=public" \
    --form "media_ids[]=${PHOTO_ID}" \
	--header "Authorization: Bearer ${MASTODON_USER_TOKEN}" \
	https://${MASTODON_SERVER}/api/v1/statuses \
    | cb

cb | cat

cb | jq .

