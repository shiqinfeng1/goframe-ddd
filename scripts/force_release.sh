#!/bin/bash

APP_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${APP_ROOT}/scripts/lib/init.sh"

if [ $# -ne 1 ];then
  app::log::error "Usage: force_release.sh v1.0.0"
  exit 1
fi

version="$1"

set +o errexit
# 1. delete old version
git tag -d ${version}
git push origin --delete ${version}

# 2. create a new tag
git tag -a ${version} -m "release ${version}"
git push origin master
git push origin ${version}

# 3. release the new release
pushd ${APP_ROOT}
# try to delete target github release if exist to avoid create error
app::log::info "delete github release with tag ${version} if exist"
github-release delete  \
  --user marmotedu\
  --repo app  \
  --tag ${version} &> /dev/null

make release
