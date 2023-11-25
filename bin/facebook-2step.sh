#!/usr/bin/env bash
#
# create a photo and then a post with the photo
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
    --form "published=false" \
    --form "caption=test from facebook-2step.sh on $(date -u)" \
    --form "access_token=${FACEBOOK_PAGE_ACCESS_TOKEN}" \
    "https://graph.facebook.com/me/photos" \
    | cb

cb | cat

cb | jq .

PHOTO_ID=$(cb | jq --raw-output .id)
echo "INFO: photo id is '${PHOTO_ID}'"

curl \
    --request POST \
    --data "message=Testing multi-photo post from facebook-2step!" \
    --data "attached_media[0]={"media_fbid":"${PHOTO_ID}"}" \
    --data "access_token=${FACEBOOK_PAGE_ACCESS_TOKEN}" \
    "https://graph.facebook.com/me/feed"