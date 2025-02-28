#!/usr/bin/env bash

# Copyright @2025-2028 <SieYuan> . All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


version="${VERSION}"
if [ "${version}" == "" ];then
  version=v`gsemver bump patch`
fi
if [ -z "`git tag -l ${version}`" ];then
  git tag -a -m "release version ${version}" ${version}
fi
