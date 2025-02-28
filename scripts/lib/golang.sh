#!/usr/bin/env bash

# shellcheck disable=SC2034 # Variables sourced in other scripts.

# The server platform we are building on.
readonly APP_SUPPORTED_SERVER_PLATFORMS=(
  linux/amd64
)

# The set of server targets that we are only building for Linux
# If you update this list, please also update build/BUILD.
# 返回一个包含特定元素的数组中的所有元素
golang::server_targets() {
  local targets=(
    mgrid
  )
  echo "${targets[@]}"
}

# 获取 golang::server_targets 函数输出的结果，将其按空格分割成多个元素，存储到数组 APP_SERVER_TARGETS 中
IFS=" " read -ra APP_SERVER_TARGETS <<< "$(golang::server_targets)"
readonly APP_SERVER_TARGETS
# 从 APP_SERVER_TARGETS 数组的每个元素中提取文件名部分，存储到另一个只读数组 APP_SERVER_BINARIES 中
readonly APP_SERVER_BINARIES=("${APP_SERVER_TARGETS[@]##*/}")


# ------------
# NOTE: All functions that return lists should use newlines.
# bash functions can't return arrays, and spaces are tricky, so newline
# separators are the preferred pattern.
# To transform a string of newline-separated items to an array, use app::util::read-array:
# app::util::read-array FOO < <(golang::dups a b c a)
#
# ALWAYS remember to quote your subshells. Not doing so will break in
# bash 4.3, and potentially cause other issues.
# ------------

# Returns a sorted newline-separated list containing only duplicated items.
# 找出传递给该函数的所有参数中的重复项，并将这些重复项输出。
golang::dups() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort | uniq -d
}

# Returns a sorted newline-separated list with duplicated items removed.
# 去除传递给该函数的所有参数中的重复项，并将去重后的参数按字典序排序后输出。
golang::dedup() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort -u
}

# Depends on values of user-facing APP_BUILD_PLATFORMS, APP_FASTBUILD,
# and APP_BUILDER_OS.
# Configures APP_SERVER_PLATFORMS, then sets them
# to readonly.
# The configured vars will only contain platforms allowed by the
# APP_SUPPORTED* vars at the top of this file.
# 声明一个名为 APP_SERVER_PLATFORMS 的数组变量
declare -a APP_SERVER_PLATFORMS
# 根据环境变量APP_BUILD_PLATFORMS或者APP_FASTBUILD，设置APP_SERVER_PLATFORMS
golang::setup_platforms() {
  # 如果APP_BUILD_PLATFORMS非空
  if [[ -n "${APP_BUILD_PLATFORMS:-}" ]]; then
    # APP_BUILD_PLATFORMS needs to be read into an array before the next
    # step, or quoting treats it all as one element.
    # 声明一个名为 platforms 的局部数组变量
    local -a platforms
    IFS=" " read -ra platforms <<< "${APP_BUILD_PLATFORMS}"

    # Deduplicate to ensure the intersection trick with golang::dups
    # is not defeated by duplicates in user input.
    app::util::read-array platforms < <(golang::dedup "${platforms[@]}")

    # Use golang::dups to restrict the builds to the platforms in
    # APP_SUPPORTED_*_PLATFORMS. Items should only appear at most once in each
    # set, so if they appear twice after the merge they are in the intersection.
    app::util::read-array APP_SERVER_PLATFORMS < <(golang::dups \
        "${platforms[@]}" \
        "${APP_SUPPORTED_SERVER_PLATFORMS[@]}" \
      )
    readonly APP_SERVER_PLATFORMS
  # 如果 APP_FASTBUILD 已定义
  elif [[ "${APP_FASTBUILD:-}" == "true" ]]; then
    APP_SERVER_PLATFORMS=(linux/amd64)
    readonly APP_SERVER_PLATFORMS
  else
    APP_SERVER_PLATFORMS=("${APP_SUPPORTED_SERVER_PLATFORMS[@]}")
    readonly APP_SERVER_PLATFORMS
  fi
}

golang::setup_platforms

# 定义了一个只读数组 APP_ALL_TARGETS，它包含了另一个数组 APP_SERVER_TARGETS 的所有元素
readonly APP_ALL_TARGETS=(
  "${APP_SERVER_TARGETS[@]}"
)
# 定义了另一个只读数组 APP_ALL_BINARIES，它存储了 APP_ALL_TARGETS 数组中每个元素去除路径前缀后剩余的文件名部分
readonly APP_ALL_BINARIES=("${APP_ALL_TARGETS[@]##*/}")

# Asks golang what it thinks the host platform is. The go tool chain does some
# slightly different things when the target platform matches the host platform.
golang::host_platform() {
  echo "$(go env GOHOSTOS)/$(go env GOHOSTARCH)"
}

# Ensure the go tool exists and is a viable version.
golang::verify_go_version() {
  if [[ -z "$(command -v go)" ]]; then
    log::usage_from_stdin <<EOF
Can't find 'go' in PATH, please fix and retry.
See http://golang.org/doc/install for installation instructions.
EOF
    return 2
  fi

  local go_version
  IFS=" " read -ra go_version <<< "$(go version)"
  local minimum_go_version
  minimum_go_version=go1.24.0
  if [[ "${minimum_go_version}" != $(echo -e "${minimum_go_version}\n${go_version[2]}" | sort -s -t. -k 1,1 -k 2,2n -k 3,3n | head -n1) && "${go_version[2]}" != "devel" ]]; then
    log::usage_from_stdin <<EOF
Detected go version: ${go_version[*]}.
APP requires ${minimum_go_version} or greater.
Please install ${minimum_go_version} or later.
EOF
    return 2
  fi
}

# golang::setup_env will check that the `go` commands is available in
# ${PATH}. It will also check that the Go version is good enough for the
# APP build.
#
# Outputs:
#   env-var GOBIN is unset (we want binaries in a predictable place)
#   env-var GO15VENDOREXPERIMENT=1
#   env-var GO111MODULE=on
golang::setup_env() {
  golang::verify_go_version

  # Unset GOBIN in case it already exists in the current session.
  unset GOBIN

  # This is for sanity.  Without it, user umasks leak through into release
  # artifacts.
  # 将当前 shell 会话的权限掩码设置为 0022， 表示在创建文件或目录时，要从所属组和其他用户的权限中去除写权限
  umask 0022
}
