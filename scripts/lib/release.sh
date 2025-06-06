#!/usr/bin/env bash

# This file creates release artifacts (tar files, container images) that are
# ready to distribute to install or distribute to end users.

###############################################################################
# Most of the ::release:: namespace functions have been moved to
# github.com/app/release.  Have a look in that repo and specifically in
# lib/releaselib.sh for ::release::-related functionality.
###############################################################################

# Tencent cos configuration

# This is where the final release artifacts are created locally
readonly RELEASE_STAGE="${LOCAL_OUTPUT_ROOT}/release-stage"
readonly RELEASE_TARS="${LOCAL_OUTPUT_ROOT}/release-tars"
readonly RELEASE_IMAGES="${LOCAL_OUTPUT_ROOT}/release-images"

# APP github account info
readonly APP_GITHUB_ORG=shiqinfeng1
readonly APP_GITHUB_REPO=goframe-ddd

readonly ARTIFACT=app.tar.gz
readonly CHECKSUM=${ARTIFACT}.sha1sum

APP_BUILD_CONFORMANCE=${APP_BUILD_CONFORMANCE:-y}
APP_BUILD_PULL_LATEST_IMAGES=${APP_BUILD_PULL_LATEST_IMAGES:-y}

# Validate a ci version
#
# Globals:
#   None
# Arguments:
#   version
# Returns:
#   If version is a valid ci version
# Sets:                    (e.g. for '1.2.3-alpha.4.56+abcdef12345678')
#   VERSION_MAJOR          (e.g. '1')
#   VERSION_MINOR          (e.g. '2')
#   VERSION_PATCH          (e.g. '3')
#   VERSION_PRERELEASE     (e.g. 'alpha')
#   VERSION_PRERELEASE_REV (e.g. '4')
#   VERSION_BUILD_INFO     (e.g. '.56+abcdef12345678')
#   VERSION_COMMITS        (e.g. '56')
function release::parse_and_validate_ci_version() {
  # Accept things like "v1.2.3-alpha.4.56+abcdef12345678" or "v1.2.3-beta.4"
  local -r version_regex="^v(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)-([a-zA-Z0-9]+)\\.(0|[1-9][0-9]*)(\\.(0|[1-9][0-9]*)\\+[0-9a-f]{7,40})?$"
  local -r version="${1-}"
  [[ "${version}" =~ ${version_regex} ]] || {
    log::error "Invalid ci version: '${version}', must match regex ${version_regex}"
    return 1
  }

  # The VERSION variables are used when this file is sourced, hence
  # the shellcheck SC2034 'appears unused' warning is to be ignored.

  # shellcheck disable=SC2034
  VERSION_MAJOR="${BASH_REMATCH[1]}"
  # shellcheck disable=SC2034
  VERSION_MINOR="${BASH_REMATCH[2]}"
  # shellcheck disable=SC2034
  VERSION_PATCH="${BASH_REMATCH[3]}"
  # shellcheck disable=SC2034
  VERSION_PRERELEASE="${BASH_REMATCH[4]}"
  # shellcheck disable=SC2034
  VERSION_PRERELEASE_REV="${BASH_REMATCH[5]}"
  # shellcheck disable=SC2034
  VERSION_BUILD_INFO="${BASH_REMATCH[6]}"
  # shellcheck disable=SC2034
  VERSION_COMMITS="${BASH_REMATCH[7]}"
}

# ---------------------------------------------------------------------------
# Build final release artifacts
function release::clean_cruft() {
  # Clean out cruft
  find "${RELEASE_STAGE}" -name '*~' -exec rm {} \;
  find "${RELEASE_STAGE}" -name '#*#' -exec rm {} \;
  find "${RELEASE_STAGE}" -name '.DS*' -exec rm {} \;
}

function release::package_tarballs() {
  # Clean out any old releases
  rm -rf "${RELEASE_STAGE}" "${RELEASE_TARS}" "${RELEASE_IMAGES}"
  mkdir -p "${RELEASE_TARS}"
  release::package_src_tarball &
  release::package_app_manifests_tarball &
  release::package_server_tarballs &
  util::wait-for-jobs || { log::error "previous tarball phase failed"; return 1; }

  release::package_final_tarball & # _final depends on some of the previous phases
  util::wait-for-jobs || { log::error "previous tarball phase failed"; return 1; }
}


# Package the source code we built, for compliance/licensing/audit/yadda.
function release::package_src_tarball() {
  local -r src_tarball="${RELEASE_TARS}/app-src.tar.gz"
  log::status "Building tarball: src"
  if [[ "${APP_GIT_TREE_STATE-}" = 'clean' ]]; then
    git archive -o "${src_tarball}" HEAD
  else
  # 用于排除一些不需要的文件和目录
    find "${APP_ROOT}" -mindepth 1 -maxdepth 1 \
      ! \( \
      \( -path "${APP_ROOT}"/_\* -o \
      -path "${APP_ROOT}"/.git\* -o \
      -path "${APP_ROOT}"/.gitignore\* -o \
      -path "${APP_ROOT}"/.gsemver.yaml\* -o \
      -path "${APP_ROOT}"/.config\* -o \
      -path "${APP_ROOT}"/.chglog\* -o \
      -path "${APP_ROOT}"/.gitlint -o \
      -path "${APP_ROOT}"/.golangci.yaml -o \
      -path "${APP_ROOT}"/.goreleaser.yml -o \
      -path "${APP_ROOT}"/.note.md -o \
      -path "${APP_ROOT}"/.todo.md -o \
      -path "${APP_ROOT}"/deployments\* \
      \) -prune \
      \) -print0 \
      | "${TAR}" czf "${src_tarball}" --transform "s|${APP_ROOT#/*}|app|" --null -T -
  fi
}

# Package up all of the server binaries
function release::package_server_tarballs() {
  # Find all of the built client binaries
  # 将 LOCAL_OUTPUT_BINPATH 目录下的两级子目录列表赋值给数组变量 long_platforms
  # 如果 LOCAL_OUTPUT_BINPATH 是 /path/to/bin，且该目录下有 linux/amd64 和 darwin/arm64 两个二级子目录，
  # 那么 long_platforms 数组将包含这两个路径。
  local long_platforms=("${LOCAL_OUTPUT_BINPATH}"/*/*)
  # 判断 APP_BUILD_PLATFORMS 是否有值时，避免出现未定义变量的错误
  if [[ -n ${APP_BUILD_PLATFORMS-} ]]; then
    # 如果 APP_BUILD_PLATFORMS 有值，就将其值按空格分割后存储到 long_platforms 数组中。
    read -ra long_platforms <<< "${APP_BUILD_PLATFORMS}"
  fi

  # 遍历 long_platforms 数组中的每个平台路径，为每个平台创建一个独立的发布目录结构，
  # 将服务器二进制文件复制到相应的目录中，清理不必要的文件，然后将该平台的发布内容打包成一个压缩文件（.tar.gz）。
  # 同时，为了提高效率，每个平台的处理过程会在后台以子进程的方式并行执行
  for platform_long in "${long_platforms[@]}"; do
    local platform
    local platform_tag
    platform=${platform_long##${LOCAL_OUTPUT_BINPATH}/} # Strip LOCAL_OUTPUT_BINPATH
    platform_tag=${platform/\//-} # Replace a "/" for a "-"
    log::status "Starting tarball: server $platform_tag"

    (
    local release_stage="${RELEASE_STAGE}/server/${platform_tag}/app"
    rm -rf "${release_stage}"
    mkdir -p "${release_stage}/server/bin"

    local server_bins=("${APP_SERVER_BINARIES[@]}")

    # This fancy expression will expand to prepend a path
    # (${LOCAL_OUTPUT_BINPATH}/${platform}/) to every item in the
    # server_bins array.
    # 为 server_bins 数组中的每个元素添加前缀 ${LOCAL_OUTPUT_BINPATH}/${platform}/
    # 然后使用 cp 命令将这些文件复制到 release_stage/server/bin/ 目录中
    cp "${server_bins[@]/#/${LOCAL_OUTPUT_BINPATH}/${platform}/}" \
      "${release_stage}/server/bin/"

    release::clean_cruft

    local package_name="${RELEASE_TARS}/app-server-${platform_tag}.tar.gz"
    release::create_tarball "${package_name}" "${release_stage}/.."
    ) &
  done

  log::status "Waiting on tarballs"
  util::wait-for-jobs || { log::error "server tarball creation failed"; exit 1; }
}

# Package up all of the server binaries in docker images
function release::build_server_images() {
  # Clean out any old images
  rm -rf "${RELEASE_IMAGES}"
  local platform
  for platform in "${APP_SERVER_PLATFORMS[@]}"; do
    local platform_tag
    local arch
    platform_tag=${platform/\//-} # Replace a "/" for a "-"
    arch=$(basename "${platform}")
    log::status "Building images: $platform_tag"

    local release_stage
    release_stage="${RELEASE_STAGE}/server/${platform_tag}/app"
    rm -rf "${release_stage}"
    mkdir -p "${release_stage}/server/bin"

    release::create_docker_images_for_server "${release_stage}/server/bin" "${arch}"
  done
}

function release::md5() {
  if which md5 >/dev/null 2>&1; then
    md5 -q "$1"
  else
    md5sum "$1" | awk '{ print $1 }'
  fi
}

function release::sha1() {
  if which sha1sum >/dev/null 2>&1; then
    sha1sum "$1" | awk '{ print $1 }'
  else
    shasum -a1 "$1" | awk '{ print $1 }'
  fi
}

function release::build_conformance_image() {
  local -r arch="$1"
  local -r registry="$2"
  local -r version="$3"
  local -r save_dir="${4-}"
  log::status "Building conformance image for arch: ${arch}"
  ARCH="${arch}" REGISTRY="${registry}" VERSION="${version}" \
    make -C cluster/images/conformance/ build >/dev/null

  local conformance_tag
  conformance_tag="${registry}/conformance-${arch}:${version}"
  if [[ -n "${save_dir}" ]]; then
    "${DOCKER[@]}" save "${conformance_tag}" > "${save_dir}/conformance-${arch}.tar"
  fi
  log::status "Deleting conformance image ${conformance_tag}"
  "${DOCKER[@]}" rmi "${conformance_tag}" &>/dev/null || true
}

# This builds all the release docker images (One docker image per binary)
# Args:
#  $1 - binary_dir, the directory to save the tared images to.
#  $2 - arch, architecture for which we are building docker images.
function release::create_docker_images_for_server() {
  # Create a sub-shell so that we don't pollute the outer environment
  (
    local binary_dir
    local arch
    local binaries
    local images_dir
    binary_dir="$1"
    arch="$2"
    binaries=$(build::get_docker_wrapped_binaries "${arch}")
    images_dir="${RELEASE_IMAGES}/${arch}"
    mkdir -p "${images_dir}"

    # k8s.gcr.io is the constant tag in the docker archives, this is also the default for config scripts in GKE.
    # We can use APP_DOCKER_REGISTRY to include and extra registry in the docker archive.
    # If we use APP_DOCKER_REGISTRY="k8s.gcr.io", then the extra tag (same) is ignored, see release_docker_image_tag below.
    local -r docker_registry="k8s.gcr.io"
    # Docker tags cannot contain '+'
    local docker_tag="${APP_GIT_VERSION/+/_}"
    if [[ -z "${docker_tag}" ]]; then
      log::error "git version information missing; cannot create Docker tag"
      return 1
    fi

    # provide `--pull` argument to `docker build` if `APP_BUILD_PULL_LATEST_IMAGES`
    # is set to y or Y; otherwise try to build the image without forcefully
    # pulling the latest base image.
    local docker_build_opts
    docker_build_opts=
    if [[ "${APP_BUILD_PULL_LATEST_IMAGES}" =~ [yY] ]]; then
        docker_build_opts='--pull'
    fi

    for wrappable in $binaries; do

      local binary_name=${wrappable%%,*}
      local base_image=${wrappable##*,}
      local binary_file_path="${binary_dir}/${binary_name}"
      local docker_build_path="${binary_file_path}.dockerbuild"
      local docker_file_path="${docker_build_path}/Dockerfile"
      local docker_image_tag="${docker_registry}/${binary_name}-${arch}:${docker_tag}"

      log::status "Starting docker build for image: ${binary_name}-${arch}"
      (
        rm -rf "${docker_build_path}"
        mkdir -p "${docker_build_path}"
        ln "${binary_file_path}" "${docker_build_path}/${binary_name}"
        ln "${APP_ROOT}/build/nsswitch.conf" "${docker_build_path}/nsswitch.conf"
        chmod 0644 "${docker_build_path}/nsswitch.conf"
        cat <<EOF > "${docker_file_path}"
FROM ${base_image}
COPY ${binary_name} /usr/local/bin/${binary_name}
EOF
        # ensure /etc/nsswitch.conf exists so go's resolver respects /etc/hosts
        if [[ "${base_image}" =~ busybox ]]; then
          echo "COPY nsswitch.conf /etc/" >> "${docker_file_path}"
        fi

        "${DOCKER[@]}" build ${docker_build_opts:+"${docker_build_opts}"} -q -t "${docker_image_tag}" "${docker_build_path}" >/dev/null
        # If we are building an official/alpha/beta release we want to keep
        # docker images and tag them appropriately.
        local -r release_docker_image_tag="${APP_DOCKER_REGISTRY-$docker_registry}/${binary_name}-${arch}:${APP_DOCKER_IMAGE_TAG-$docker_tag}"
        if [[ "${release_docker_image_tag}" != "${docker_image_tag}" ]]; then
          log::status "Tagging docker image ${docker_image_tag} as ${release_docker_image_tag}"
          "${DOCKER[@]}" rmi "${release_docker_image_tag}" 2>/dev/null || true
          "${DOCKER[@]}" tag "${docker_image_tag}" "${release_docker_image_tag}" 2>/dev/null
        fi
        "${DOCKER[@]}" save -o "${binary_file_path}.tar" "${docker_image_tag}" "${release_docker_image_tag}"
        echo "${docker_tag}" > "${binary_file_path}.docker_tag"
        rm -rf "${docker_build_path}"
        ln "${binary_file_path}.tar" "${images_dir}/"

        log::status "Deleting docker image ${docker_image_tag}"
        "${DOCKER[@]}" rmi "${docker_image_tag}" &>/dev/null || true
      ) &
    done

    if [[ "${APP_BUILD_CONFORMANCE}" =~ [yY] ]]; then
      release::build_conformance_image "${arch}" "${docker_registry}" \
        "${docker_tag}" "${images_dir}" &
    fi

    util::wait-for-jobs || { log::error "previous Docker build failed"; return 1; }
    log::status "Docker builds done"
  )

}

# This will pack app-system manifests files for distros such as COS.
function release::package_app_manifests_tarball() {
  log::status "Building tarball: manifests"

  local src_dir="${APP_ROOT}/deployments"

  local release_stage="${RELEASE_STAGE}/manifests/app"
  rm -rf "${release_stage}"

  local dst_dir="${release_stage}"
  mkdir -p "${dst_dir}"
  cp -r ${src_dir}/* "${dst_dir}"
  #cp "${APP_ROOT}/cluster/gce/gci/health-monitor.sh" "${dst_dir}/health-monitor.sh"

  release::clean_cruft

  local package_name="${RELEASE_TARS}/app-manifests.tar.gz"
  release::create_tarball "${package_name}" "${release_stage}/.."
}

function release::updload_tarballs() {
  log::info "upload ${RELEASE_TARS}/* ..."
  for file in $(ls ${RELEASE_TARS}/*)
  do
    echo "[TODO] upload ${file} to anywhere you want !!!"
  done
}
# This is all the platform-independent stuff you need to run/install app.
# Arch-specific binaries will need to be downloaded separately (possibly by
# using the bundled cluster/get-app-binaries.sh script).
# Included in this tarball:
#   - Cluster spin up/down scripts and config for various cloud providers
#   - Tarballs for manifest config that are ready to be uploaded
#   - Examples (which may or may not still work)
#   - The remnants of the docs/ directory
function release::package_final_tarball() {
  log::status "Building tarball: final"

  # This isn't a "full" tarball anymore, but the release lib still expects
  # artifacts under "full/app/"
  local release_stage="${RELEASE_STAGE}/full/app"
  rm -rf "${release_stage}"
  mkdir -p "${release_stage}"

  mkdir -p "${release_stage}/client"
  cat <<EOF > "${release_stage}/client/README"
Client binaries are no longer included in the APP final tarball.

Run release/get-app-binaries.sh to download client and server binaries.
EOF

  # We want everything in /scripts.
  mkdir -p "${release_stage}/release"
  cp -R "${APP_ROOT}/scripts/release" "${release_stage}/"
  cat <<EOF > "${release_stage}/release/get-app-binaries.sh"
#!/usr/bin/env bash

# This file download app client and server binaries from tencent cos bucket.

os=linux arch=amd64 version=${APP_GIT_VERSION} 
EOF
  chmod +x ${release_stage}/release/get-app-binaries.sh

  mkdir -p "${release_stage}/server"
  cp "${RELEASE_TARS}/app-manifests.tar.gz" "${release_stage}/server/"
  cat <<EOF > "${release_stage}/server/README"
Server binary tarballs are no longer included in the APP final tarball.

Run release/get-app-binaries.sh to download client and server binaries.
EOF

  # Include hack/lib as a dependency for the cluster/ scripts
  #mkdir -p "${release_stage}/hack"
  #cp -R "${APP_ROOT}/hack/lib" "${release_stage}/hack/"

  cp -R ${APP_ROOT}/{docs,config,scripts,deployments,init,README.md,LICENSE} "${release_stage}/"

  echo "${APP_GIT_VERSION}" > "${release_stage}/version"

  release::clean_cruft

  local package_name="${RELEASE_TARS}/${ARTIFACT}"
  release::create_tarball "${package_name}" "${release_stage}/.."
}

# Build a release tarball.  $1 is the output tar name.  $2 is the base directory
# of the files to be packaged.  This assumes that ${2}/appis what is
# being packaged.
function release::create_tarball() {
  build::ensure_tar

  local tarfile=$1
  local stagingdir=$2

  "${TAR}" czf "${tarfile}" -C "${stagingdir}" app --owner=0 --group=0
}

function release::install_github_release(){
  GO111MODULE=on go install github.com/github-release/github-release@latest
}

# Require the following tools:
# - github-release
# - gsemver
# - git-chglog
# - coscmd or coscli
function release::verify_prereqs(){
  if [ -z "$(which github-release 2>/dev/null)" ]; then
    log::info "'github-release' tool not installed, try to install it."

    if ! release::install_github_release; then
      log::error "failed to install 'github-release'"
      return 1
    fi
  fi

  if [ -z "$(which git-chglog 2>/dev/null)" ]; then
    log::info "'git-chglog' tool not installed, try to install it."

    if ! go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest &>/dev/null; then
      log::error "failed to install 'git-chglog'"
      return 1
    fi
  fi

  if [ -z "$(which gsemver 2>/dev/null)" ]; then
    log::info "'gsemver' tool not installed, try to install it."

    if ! go install github.com/arnaud-deprez/gsemver@latest &>/dev/null; then
      log::error "failed to install 'gsemver'"
      return 1
    fi
  fi
}
# details:
# https://github.com/github-release/github-release
function release::check_github_token() {
  # 检查 GITHUB_TOKEN 是否已经导出
  if [ -z "${GITHUB_TOKEN:-}" ]; then
      echo "GITHUB_TOKEN 未设置，需要你手动输入。"
      # 提示用户输入 GitHub 访问令牌
      read -s -p "请输入你的 GitHub 访问令牌：" GITHUB_TOKEN

      # 验证输入是否为空
      if [ -z "$GITHUB_TOKEN" ]; then
          echo "输入为空，请重新运行脚本并输入有效的令牌。"
          exit 1
      fi

      # 导出为全局环境变量
      export GITHUB_TOKEN
  else
      echo "GITHUB_TOKEN 已经设置，值为: $GITHUB_TOKEN"
  fi
}
# Create a github release with specified tarballs.
# NOTICE: Must export 'GITHUB_TOKEN' env in the shell, details:
# https://github.com/github-release/github-release
function release::github_release() {
  # create a github release
  log::info "\n create a new github release ..."
  log::info "github-release release --user ${APP_GITHUB_ORG} --repo ${APP_GITHUB_REPO} --tag ${APP_GIT_VERSION} --description '' --pre-release"
  github-release release \
    --user ${APP_GITHUB_ORG} \
    --repo ${APP_GITHUB_REPO} \
    --tag ${APP_GIT_VERSION} \
    --description "" \
    --pre-release

  # update app tarballs
  log::info "upload ${ARTIFACT} to release ${APP_GIT_VERSION}"
  github-release upload \
    --user ${APP_GITHUB_ORG} \
    --repo ${APP_GITHUB_REPO} \
    --tag ${APP_GIT_VERSION} \
    --name ${ARTIFACT} \
    --file ${RELEASE_TARS}/${ARTIFACT}

  log::info "upload app-src.tar.gz to release ${APP_GIT_VERSION}"
  github-release upload \
    --user ${APP_GITHUB_ORG} \
    --repo ${APP_GITHUB_REPO} \
    --tag ${APP_GIT_VERSION} \
    --name "app-src.tar.gz" \
    --file ${RELEASE_TARS}/app-src.tar.gz
}

function release::generate_changelog() {
  log::info "generate CHANGELOG-${APP_GIT_VERSION#v}.md and commit it"

  git-chglog ${APP_GIT_VERSION} > ${APP_ROOT}/CHANGELOG/CHANGELOG-${APP_GIT_VERSION#v}.md

  set +o errexit
  git add ${APP_ROOT}/CHANGELOG/CHANGELOG-${APP_GIT_VERSION#v}.md
  git commit -a -m "docs(changelog): add CHANGELOG-${APP_GIT_VERSION#v}.md"
  git push -f origin main # 最后将 CHANGELOG 也 push 上去
}

