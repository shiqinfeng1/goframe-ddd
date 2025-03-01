#!/usr/bin/env bash

version="${VERSION}"
if [ "${version}" == "" ];then
  version=v`gsemver bump`
fi

# 定义语义化版本号的正则表达式
semver_pattern='^v[0-9]+\.[0-9]+\.[0-9]+$'
if [[ $version =~ $semver_pattern ]]; then
    echo "version 有效，值为: $version"
else
    echo "version 无效，不符合语义化版本号格式: ${version}"
    exit 1
fi

if [ -z "`git tag -l ${version}`" ];then
  git tag -a -m "release version ${version}" ${version}
fi
