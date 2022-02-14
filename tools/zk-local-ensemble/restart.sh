#!/usr/bin/env bash
# Bash Strict Mode (http://redsymbol.net/articles/unofficial-bash-strict-mode/)
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

${SCRIPT_DIR}/down.sh
${SCRIPT_DIR}/up.sh
