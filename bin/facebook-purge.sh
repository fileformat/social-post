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

echo "==="
echo ${FACEBOOK_PAGE_ACCESS_TOKEN}
echo "==="

curl -X GET "https://graph.facebook.com/v18.0/${FACEBOOK_PAGE_ID}/feed?access_token=${FACEBOOK_PAGE_ACCESS_TOKEN}" | cb

cb  | jq --raw-output '.data[] | .id' >${SCRIPT_HOME}/fbdelete.txt

# for each line in fbdelete.txt
for line in $(cat "${SCRIPT_HOME}/fbdelete.txt"); do
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
