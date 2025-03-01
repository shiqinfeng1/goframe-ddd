#!/usr/bin/env bash

set -o errexit
set +o nounset
set -o pipefail

# Unset CDPATH so that path interpolation can work correctly
# https://github.com/apprnetes/apprnetes/issues/52255
unset CDPATH

# Default use go modules
export GO111MODULE=on

# The root of the build/dist directory
APP_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

source "${APP_ROOT}/scripts/lib/util.sh"
source "${APP_ROOT}/scripts/lib/logging.sh"
source "${APP_ROOT}/scripts/lib/color.sh"

log::install_errexit

source "${APP_ROOT}/scripts/lib/version.sh"
source "${APP_ROOT}/scripts/lib/golang.sh"
