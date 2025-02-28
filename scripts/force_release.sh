#!/bin/bash

APP_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${APP_ROOT}/scripts/lib/init.sh"

if [ $# -ne 1 ];then
  log::error "Usage: force_release.sh v1.0.0"
  exit 1
fi

version="$1"

set +o errexit
# 1. delete old version
# 删除本地和远程仓库中的特定 Git 标签
git tag -d ${version}
git push origin --delete ${version}

# 2. create a new tag
# 创建新的tag
git tag -a ${version} -m "release ${version}"
git push origin master
git push origin ${version}

# 3. release the new release
# 将当前工作目录切换到 ${APP_ROOT} 所代表的目录
pushd ${APP_ROOT}
# try to delete target github release if exist to avoid create error
log::info "delete github release with tag ${version} if exist"
# 使用 github-release 工具从 GitHub 上删除指定仓库的特定版本标签对应的发布版本
github-release delete  \
  --user shiqinfeng1\
  --repo goframe-ddd  \
  --tag ${version} &> /dev/null

make release
