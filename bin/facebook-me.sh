#!/usr/bin/env bash
#
#
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

curl --request GET \
    "https://graph.facebook.com/v18.0/me?fields=app_id,about,can_post,description,description_html,cover,has_added_app,link,name,name_with_location_descriptor,page_token,username,verification_status,website&access_token=${FACEBOOK_PAGE_ACCESS_TOKEN}" \
    | cb

cb | jq .
