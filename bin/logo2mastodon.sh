#!/usr/bin/env bash
#
# post a logo to mastodon after converting it from svg
#

set -o errexit
set -o pipefail
set -o nounset

SCRIPT_HOME="$( cd "$( dirname "$0" )" && pwd )"
REPO_HOME=$(realpath "${SCRIPT_HOME}/..")

LOGOHANDLE=${1:-BAD}
if [ "${LOGOHANDLE}" == "BAD" ]; then
    echo "usage: $0 <LOGOHANDLE>"
    echo "       LOGOHANDLE is the id of a logo on [VectorLogoZone](https://www.vectorlogo.zone/)"
    exit 2
fi

RSVG=$(which rsvg-convert)
if [ -z "${RSVG}" ]; then
    echo "ERROR: rsvg-convert not found"
    echo "       try 'brew install librsvg' on macOS"
    echo "       try 'apt-get install librsvg2-bin' on Ubuntu"
    exit 1
fi

ENV=${REPO_HOME}/.env
if [ ! -f "${ENV}" ]; then
    echo "ERROR: ${ENV} not found"
    exit 1
fi
export $(cat ${ENV} | xargs)

SOCIAL_POST="${REPO_HOME}/dist/social-post"
if [ ! -f "${SOCIAL_POST}" ]; then
    echo "ERROR: ${SOCIAL_POST} not found"
    echo "       try 'ruh.sh' in the repo root"
    exit 1
fi

echo "INFO: downloading ${LOGOHANDLE} from VectorLogoZone"
SVGDATA=$(curl -s "https://www.vectorlogo.zone/logos/${LOGOHANDLE}/${LOGOHANDLE}-ar21.svg")

PNGFILE=$(mktemp -t ${LOGOHANDLE}.XXXXXX.png)
echo "INFO: converting to PNG (${PNGFILE})"
echo "${SVGDATA}" | ${RSVG} --background-color=white --width=1280 --dpi-x=600 --keep-aspect-ratio --format png --output "${PNGFILE}"

echo "INFO: downloading metadata"
SVGMETA=$(curl -s "https://upload.vectorlogo.zone/logos/${LOGOHANDLE}/logoinfo.json")
NAME=$(echo "${SVGMETA}" | jq -r '.name')
WEBSITE=$(echo "${SVGMETA}" | jq -r '.website')



${SOCIAL_POST} mastodon --image="${PNGFILE}" --image-caption="${NAME} Logo" "Logo for ${NAME} https://vlz.one/${LOGOHANDLE}"

echo "INFO: complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"