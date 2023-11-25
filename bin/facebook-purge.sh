#!/usr/bin/env bash
#
# testing the fb api locally
#
# https://developers.facebook.com/tools/explorer/
# https://developers.facebook.com/tools/debug/accesstoken/

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

TMP_DIR="${BASE_DIR}/tmp"
if [ ! -d "${TMP_DIR}" ]; then
    echo "INFO: creating '${TMP_DIR}'"
    mkdir -p "${TMP_DIR}"
fi

curl -X GET "https://graph.facebook.com/v18.0/${FACEBOOK_PAGE_ID}/feed?access_token=${FACEBOOK_PAGE_ACCESS_TOKEN}" | cb

cb  | jq --raw-output '.data[] | .id' >${TMP_DIR}/fbdelete.tmp

for line in $(cat "${TMP_DIR}/fbdelete.tmp"); do
    # delete the post
    echo "INFO: deleting '${line}'"
    curl -i -X DELETE "https://graph.facebook.com/v18.0/${line}?access_token=${FACEBOOK_PAGE_ACCESS_TOKEN}"
done


#curl -i -X GET "https://graph.facebook.com/v18.0/${FACEBOOK_PAGE_ID}/feed?access_token=${FACEBOOK_PAGE_ACCESS_TOKEN}"
exit 1

curl \
    --request POST \
    --header "Content-Type: application/json" \
    -d "{\"message\":\"your_message_text\",\"access_token\":\"${FACEBOOK_PAGE_ACCESS_TOKEN}\"}" \
    https://graph.facebook.com/v18.0/${FACEBOOK_PAGE_ID}/feed \
    | jq .
