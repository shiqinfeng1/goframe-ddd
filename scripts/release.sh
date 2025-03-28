#!/usr/bin/env bash

# Build a APP release.  This will build the binaries, create the Docker
# images and other build artifacts.

set -o errexit
set -o nounset
set -o pipefail

APP_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${APP_ROOT}/scripts/common.sh"
source "${APP_ROOT}/scripts/lib/release.sh"

APP_RELEASE_RUN_TESTS=${APP_RELEASE_RUN_TESTS-y}

golang::setup_env
build::verify_prereqs false
release::verify_prereqs
#build::build_image
build::build_command
release::package_tarballs
release::updload_tarballs
git push origin ${VERSION}
release::check_github_token
release::github_release
release::generate_changelog

