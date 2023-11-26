#!/usr/bin/env bash
#
# generate completions for the various shells
#

set -o errexit
set -o pipefail
set -o nounset

SCRIPT_HOME="$( cd "$( dirname "$0" )" && pwd )"
REPO_HOME=$(realpath "${SCRIPT_HOME}/..")
COMPLETIONS_DIR="${REPO_HOME}/completions"

echo "INFO: completions starting at $(date -u +%Y-%m-%dT%H:%M:%SZ)"

rm -rf "${COMPLETIONS_DIR}"
mkdir "${COMPLETIONS_DIR}"
for sh in bash zsh fish; do
	echo "INFO: generating completions for $sh..."
	LOG_LEVEL=12 go run "${REPO_HOME}/main.go" completion "$sh" >""${COMPLETIONS_DIR}"/social-post.$sh"
done

echo "INFO: completions complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
