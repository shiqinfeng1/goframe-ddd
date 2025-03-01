#!/usr/bin/env bash

# shellcheck disable=SC2034 # Variables sourced in other scripts.

# Common utilities, variables and checks for all build scripts.
set -o errexit
set -o nounset
set -o pipefail

# Unset CDPATH, having it set messes up with script import paths
unset CDPATH

USER_ID=$(id -u)
GROUP_ID=$(id -g)

DOCKER_OPTS=${DOCKER_OPTS:-""}
IFS=" " read -r -a DOCKER <<< "docker ${DOCKER_OPTS}"
DOCKER_HOST=${DOCKER_HOST:-""}
DOCKER_MACHINE_NAME=${DOCKER_MACHINE_NAME:-"app-dev"}
readonly DOCKER_MACHINE_DRIVER=${DOCKER_MACHINE_DRIVER:-"virtualbox --virtualbox-cpu-count -1"}

# This will canonicalize the path
APP_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd -P)

source "${APP_ROOT}/scripts/lib/init.sh"

# Constants
readonly APP_BUILD_IMAGE_REPO=app-build
#readonly APP_BUILD_IMAGE_CROSS_TAG="$(cat "${APP_ROOT}/build/build-image/cross/VERSION")"

readonly APP_DOCKER_REGISTRY="${APP_DOCKER_REGISTRY:-k8s.gcr.io}"
readonly APP_BASE_IMAGE_REGISTRY="${APP_BASE_IMAGE_REGISTRY:-us.gcr.io/k8s-artifacts-prod/build-image}"

# This version number is used to cause everyone to rebuild their data containers
# and build image.  This is especially useful for automated build systems like
# Jenkins.
#
# Increment/change this number if you change the build image (anything under
# build/build-image) or change the set of volumes in the data container.
#readonly APP_BUILD_IMAGE_VERSION_BASE="$(cat "${APP_ROOT}/build/build-image/VERSION")"
#readonly APP_BUILD_IMAGE_VERSION="${APP_BUILD_IMAGE_VERSION_BASE}-${APP_BUILD_IMAGE_CROSS_TAG}"

# Here we map the output directories across both the local and remote _output
# directories:
#
# *_OUTPUT_ROOT    - the base of all output in that environment.
# *_OUTPUT_SUBPATH - location where golang stuff is built/cached.  Also
#                    persisted across docker runs with a volume mount.
# *_OUTPUT_BINPATH - location where final binaries are placed.  If the remote
#                    is really remote, this is the stuff that has to be copied
#                    back.
# OUT_DIR can come in from the Makefile, so honor it.
readonly LOCAL_OUTPUT_ROOT="${APP_ROOT}/${OUT_DIR:-_output}"
readonly LOCAL_OUTPUT_SUBPATH="${LOCAL_OUTPUT_ROOT}/platforms"
readonly LOCAL_OUTPUT_BINPATH="${LOCAL_OUTPUT_SUBPATH}"
readonly LOCAL_OUTPUT_GOPATH="${LOCAL_OUTPUT_SUBPATH}/go"
readonly LOCAL_OUTPUT_IMAGE_STAGING="${LOCAL_OUTPUT_ROOT}/images"

# This is the port on the workstation host to expose RSYNC on.  Set this if you
# are doing something fancy with ssh tunneling.
readonly APP_RSYNC_PORT="${APP_RSYNC_PORT:-}"

# This is the port that rsync is running on *inside* the container. This may be
# mapped to APP_RSYNC_PORT via docker networking.
readonly APP_CONTAINER_RSYNC_PORT=8730

# Get the set of master binaries that run in Docker (on Linux)
# Entry format is "<name-of-binary>,<base-image>".
# Binaries are placed in /usr/local/bin inside the image.
#
# $1 - server architecture
build::get_docker_wrapped_binaries() {
  local arch=$1
  local debian_base_version=v2.1.0
  local debian_iptables_version=v12.1.0
  ### If you change any of these lists, please also update DOCKERIZED_BINARIES
  ### in build/BUILD. And golang::server_image_targets
  local targets=(
    "app-apiserver,${APP_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "app-controller-manager,${APP_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "app-scheduler,${APP_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "app-proxy,${APP_BASE_IMAGE_REGISTRY}/debian-iptables-${arch}:${debian_iptables_version}"
  )

  echo "${targets[@]}"
}

# ---------------------------------------------------------------------------
# Basic setup functions

# Verify that the right utilities and such are installed for building app. Set
# up some dynamic constants.
# Args:
#   $1 - boolean of whether to require functioning docker (default true)
#
# Vars set:
#   APP_ROOT_HASH
#   APP_BUILD_IMAGE_TAG_BASE
#   APP_BUILD_IMAGE_TAG
#   APP_BUILD_IMAGE
#   APP_BUILD_CONTAINER_NAME_BASE
#   APP_BUILD_CONTAINER_NAME
#   APP_DATA_CONTAINER_NAME_BASE
#   APP_DATA_CONTAINER_NAME
#   APP_RSYNC_CONTAINER_NAME_BASE
#   APP_RSYNC_CONTAINER_NAME
#   DOCKER_MOUNT_ARGS
#   LOCAL_OUTPUT_BUILD_CONTEXT
function build::verify_prereqs() {
  local -r require_docker=${1:-true}
  log::status "Verifying Prerequisites...."
  build::ensure_tar || return 1
  build::ensure_rsync || return 1
  if ${require_docker}; then
    build::ensure_docker_in_path || return 1
    app::util::ensure_docker_daemon_connectivity || return 1

    if (( APP_VERBOSE > 6 )); then
      log::status "Docker Version:"
      "${DOCKER[@]}" version | log::info_from_stdin
    fi
  fi

  APP_GIT_BRANCH=$(git symbolic-ref --short -q HEAD 2>/dev/null || true)
  log::info "APP_GIT_BRANCH = ${APP_GIT_BRANCH}" 
  APP_ROOT_HASH=$(build::short_hash "${HOSTNAME:-}:${APP_ROOT}:${APP_GIT_BRANCH}")
  log::info "APP_ROOT_HASH = ${APP_ROOT_HASH}" 
  APP_BUILD_IMAGE_TAG_BASE="build-${APP_ROOT_HASH}"
  log::info "APP_BUILD_IMAGE_TAG_BASE = ${APP_BUILD_IMAGE_TAG_BASE}" 
  #APP_BUILD_IMAGE_TAG="${APP_BUILD_IMAGE_TAG_BASE}-${APP_BUILD_IMAGE_VERSION}"
  #APP_BUILD_IMAGE="${APP_BUILD_IMAGE_REPO}:${APP_BUILD_IMAGE_TAG}"
  APP_BUILD_CONTAINER_NAME_BASE="app-build-${APP_ROOT_HASH}"
  log::info "APP_BUILD_CONTAINER_NAME_BASE = ${APP_BUILD_CONTAINER_NAME_BASE}" 
  #APP_BUILD_CONTAINER_NAME="${APP_BUILD_CONTAINER_NAME_BASE}-${APP_BUILD_IMAGE_VERSION}"
  APP_RSYNC_CONTAINER_NAME_BASE="app-rsync-${APP_ROOT_HASH}"
  log::info "APP_RSYNC_CONTAINER_NAME_BASE = ${APP_RSYNC_CONTAINER_NAME_BASE}" 
  #APP_RSYNC_CONTAINER_NAME="${APP_RSYNC_CONTAINER_NAME_BASE}-${APP_BUILD_IMAGE_VERSION}"
  APP_DATA_CONTAINER_NAME_BASE="app-build-data-${APP_ROOT_HASH}"
  log::info "APP_DATA_CONTAINER_NAME_BASE = ${APP_DATA_CONTAINER_NAME_BASE}" 
  #APP_DATA_CONTAINER_NAME="${APP_DATA_CONTAINER_NAME_BASE}-${APP_BUILD_IMAGE_VERSION}"
  #DOCKER_MOUNT_ARGS=(--volumes-from "${APP_DATA_CONTAINER_NAME}")
  #LOCAL_OUTPUT_BUILD_CONTEXT="${LOCAL_OUTPUT_IMAGE_STAGING}/${APP_BUILD_IMAGE}"

  version::get_version_vars
  #version::save_version_vars "${APP_ROOT}/.dockerized-app-version-defs"
}

# ---------------------------------------------------------------------------
# Utility functions

function build::docker_available_on_osx() {
  if [[ -z "${DOCKER_HOST}" ]]; then
    if [[ -S "/var/run/docker.sock" ]]; then
      log::status "Using Docker for MacOS"
      return 0
    fi

    log::status "No docker host is set. Checking options for setting one..."
    if [[ -z "$(which docker-machine)" ]]; then
      log::status "It looks like you're running Mac OS X, yet neither Docker for Mac nor docker-machine can be found."
      log::status "See: https://docs.docker.com/engine/installation/mac/ for installation instructions."
      return 1
    elif [[ -n "$(which docker-machine)" ]]; then
      build::prepare_docker_machine
    fi
  fi
}

function build::prepare_docker_machine() {
  log::status "docker-machine was found."

  local available_memory_bytes
  available_memory_bytes=$(sysctl -n hw.memsize 2>/dev/null)

  local bytes_in_mb=1048576

  # Give virtualbox 1/2 the system memory. Its necessary to divide by 2, instead
  # of multiple by .5, because bash can only multiply by ints.
  local memory_divisor=2

  local virtualbox_memory_mb=$(( available_memory_bytes / (bytes_in_mb * memory_divisor) ))

  docker-machine inspect "${DOCKER_MACHINE_NAME}" &> /dev/null || {
    log::status "Creating a machine to build APP"
    docker-machine create --driver "${DOCKER_MACHINE_DRIVER}" \
      --virtualbox-memory "${virtualbox_memory_mb}" \
      --engine-env HTTP_PROXY="${APPRNETES_HTTP_PROXY:-}" \
      --engine-env HTTPS_PROXY="${APPRNETES_HTTPS_PROXY:-}" \
      --engine-env NO_PROXY="${APPRNETES_NO_PROXY:-127.0.0.1}" \
      "${DOCKER_MACHINE_NAME}" > /dev/null || {
      log::error "Something went wrong creating a machine."
      log::error "Try the following: "
      log::error "docker-machine create -d ${DOCKER_MACHINE_DRIVER} --virtualbox-memory ${virtualbox_memory_mb} ${DOCKER_MACHINE_NAME}"
      return 1
    }
  }
  docker-machine start "${DOCKER_MACHINE_NAME}" &> /dev/null
  # it takes `docker-machine env` a few seconds to work if the machine was just started
  local docker_machine_out
  while ! docker_machine_out=$(docker-machine env "${DOCKER_MACHINE_NAME}" 2>&1); do
    if [[ ${docker_machine_out} =~ "Error checking TLS connection" ]]; then
      echo "${docker_machine_out}"
      docker-machine regenerate-certs "${DOCKER_MACHINE_NAME}"
    else
      sleep 1
    fi
  done
  eval "$(docker-machine env "${DOCKER_MACHINE_NAME}")"
  log::status "A Docker host using docker-machine named '${DOCKER_MACHINE_NAME}' is ready to go!"
  return 0
}

function build::is_gnu_sed() {
  [[ $(sed --version 2>&1) == *GNU* ]]
}

function build::ensure_rsync() {
  if [[ -z "$(which rsync)" ]]; then
    log::error "Can't find 'rsync' in PATH, please fix and retry."
    return 1
  fi
}

function build::update_dockerfile() {
  if build::is_gnu_sed; then
    sed_opts=(-i)
  else
    sed_opts=(-i '')
  fi
  sed "${sed_opts[@]}" "s/APP_BUILD_IMAGE_CROSS_TAG/${APP_BUILD_IMAGE_CROSS_TAG}/" "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
}

function  build::set_proxy() {
  if [[ -n "${APPRNETES_HTTPS_PROXY:-}" ]]; then
    echo "ENV https_proxy $APPRNETES_HTTPS_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
  if [[ -n "${APPRNETES_HTTP_PROXY:-}" ]]; then
    echo "ENV http_proxy $APPRNETES_HTTP_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
  if [[ -n "${APPRNETES_NO_PROXY:-}" ]]; then
    echo "ENV no_proxy $APPRNETES_NO_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
}

function build::ensure_docker_in_path() {
  if [[ -z "$(which docker)" ]]; then
    log::error "Can't find 'docker' in PATH, please fix and retry."
    log::error "See https://docs.docker.com/installation/#installation for installation instructions."
    return 1
  fi
}

function build::ensure_tar() {
  if [[ -n "${TAR:-}" ]]; then
    return
  fi

  # Find gnu tar if it is available, bomb out if not.
  TAR=tar
  if which gtar &>/dev/null; then
      TAR=gtar
  else
      if which gnutar &>/dev/null; then
	  TAR=gnutar
      fi
  fi
  if ! "${TAR}" --version | grep -q GNU; then
    echo "  !!! Cannot find GNU tar. Build on Linux or install GNU tar"
    echo "      on Mac OS X (brew install gnu-tar)."
    return 1
  fi
}

function build::has_docker() {
  which docker &> /dev/null
}

function build::has_ip() {
  which ip &> /dev/null && ip -Version | grep 'iproute2' &> /dev/null
}

# Detect if a specific image exists
#
# $1 - image repo name
# $2 - image tag
function build::docker_image_exists() {
  [[ -n $1 && -n $2 ]] || {
    log::error "Internal error. Image not specified in docker_image_exists."
    exit 2
  }

  [[ $("${DOCKER[@]}" images -q "${1}:${2}") ]]
}

# Delete all images that match a tag prefix except for the "current" version
#
# $1: The image repo/name
# $2: The tag base. We consider any image that matches $2*
# $3: The current image not to delete if provided
function build::docker_delete_old_images() {
  # In Docker 1.12, we can replace this with
  #    docker images "$1" --format "{{.Tag}}"
  for tag in $("${DOCKER[@]}" images "${1}" | tail -n +2 | awk '{print $2}') ; do
    if [[ "${tag}" != "${2}"* ]] ; then
      V=3 log::status "Keeping image ${1}:${tag}"
      continue
    fi

    if [[ -z "${3:-}" || "${tag}" != "${3}" ]] ; then
      V=2 log::status "Deleting image ${1}:${tag}"
      "${DOCKER[@]}" rmi "${1}:${tag}" >/dev/null
    else
      V=3 log::status "Keeping image ${1}:${tag}"
    fi
  done
}

# Stop and delete all containers that match a pattern
#
# $1: The base container prefix
# $2: The current container to keep, if provided
function build::docker_delete_old_containers() {
  # In Docker 1.12 we can replace this line with
  #   docker ps -a --format="{{.Names}}"
  for container in $("${DOCKER[@]}" ps -a | tail -n +2 | awk '{print $NF}') ; do
    if [[ "${container}" != "${1}"* ]] ; then
      V=3 log::status "Keeping container ${container}"
      continue
    fi
    if [[ -z "${2:-}" || "${container}" != "${2}" ]] ; then
      V=2 log::status "Deleting container ${container}"
      build::destroy_container "${container}"
    else
      V=3 log::status "Keeping container ${container}"
    fi
  done
}

# Takes $1 and computes a short has for it. Useful for unique tag generation
function build::short_hash() {
  [[ $# -eq 1 ]] || {
    log::error "Internal error.  No data based to short_hash."
    exit 2
  }

  local short_hash
  if which md5 >/dev/null 2>&1; then
    short_hash=$(md5 -q -s "$1")
  else
    short_hash=$(echo -n "$1" | md5sum)
  fi
  echo "${short_hash:0:10}"
}

# Pedantically kill, wait-on and remove a container. The -f -v options
# to rm don't actually seem to get the job done, so force kill the
# container, wait to ensure it's stopped, then try the remove. This is
# a workaround for bug https://github.com/docker/docker/issues/3968.
function build::destroy_container() {
  "${DOCKER[@]}" kill "$1" >/dev/null 2>&1 || true
  if [[ $("${DOCKER[@]}" version --format '{{.Server.Version}}') = 17.06.0* ]]; then
    # Workaround https://github.com/moby/moby/issues/33948.
    # TODO: remove when 17.06.0 is not relevant anymore
    DOCKER_API_VERSION=v1.29 "${DOCKER[@]}" wait "$1" >/dev/null 2>&1 || true
  else
    "${DOCKER[@]}" wait "$1" >/dev/null 2>&1 || true
  fi
  "${DOCKER[@]}" rm -f -v "$1" >/dev/null 2>&1 || true
}

# ---------------------------------------------------------------------------
# Building


function build::clean() {
  if build::has_docker ; then
    build::docker_delete_old_containers "${APP_BUILD_CONTAINER_NAME_BASE}"
    build::docker_delete_old_containers "${APP_RSYNC_CONTAINER_NAME_BASE}"
    build::docker_delete_old_containers "${APP_DATA_CONTAINER_NAME_BASE}"
    build::docker_delete_old_images "${APP_BUILD_IMAGE_REPO}" "${APP_BUILD_IMAGE_TAG_BASE}"

    V=2 log::status "Cleaning all untagged docker images"
    "${DOCKER[@]}" rmi "$("${DOCKER[@]}" images -q --filter 'dangling=true')" 2> /dev/null || true
  fi

  if [[ -d "${LOCAL_OUTPUT_ROOT}" ]]; then
    log::status "Removing _output directory"
    rm -rf "${LOCAL_OUTPUT_ROOT}"
  fi
}

# Set up the context directory for the app-build image and build it.
function build::build_image() {
  mkdir -p "${LOCAL_OUTPUT_BUILD_CONTEXT}"
  # Make sure the context directory owned by the right user for syncing sources to container.
  chown -R "${USER_ID}":"${GROUP_ID}" "${LOCAL_OUTPUT_BUILD_CONTEXT}"

  cp /etc/localtime "${LOCAL_OUTPUT_BUILD_CONTEXT}/"

  cp "${APP_ROOT}/build/build-image/Dockerfile" "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  cp "${APP_ROOT}/build/build-image/rsyncd.sh" "${LOCAL_OUTPUT_BUILD_CONTEXT}/"
  dd if=/dev/urandom bs=512 count=1 2>/dev/null | LC_ALL=C tr -dc 'A-Za-z0-9' | dd bs=32 count=1 2>/dev/null > "${LOCAL_OUTPUT_BUILD_CONTEXT}/rsyncd.password"
  chmod go= "${LOCAL_OUTPUT_BUILD_CONTEXT}/rsyncd.password"

  build::update_dockerfile
  build::set_proxy
  build::docker_build "${APP_BUILD_IMAGE}" "${LOCAL_OUTPUT_BUILD_CONTEXT}" 'false'

  # Clean up old versions of everything
  build::docker_delete_old_containers "${APP_BUILD_CONTAINER_NAME_BASE}" "${APP_BUILD_CONTAINER_NAME}"
  build::docker_delete_old_containers "${APP_RSYNC_CONTAINER_NAME_BASE}" "${APP_RSYNC_CONTAINER_NAME}"
  build::docker_delete_old_containers "${APP_DATA_CONTAINER_NAME_BASE}" "${APP_DATA_CONTAINER_NAME}"
  build::docker_delete_old_images "${APP_BUILD_IMAGE_REPO}" "${APP_BUILD_IMAGE_TAG_BASE}" "${APP_BUILD_IMAGE_TAG}"

  build::ensure_data_container
  build::sync_to_container
}

# Build a docker image from a Dockerfile.
# $1 is the name of the image to build
# $2 is the location of the "context" directory, with the Dockerfile at the root.
# $3 is the value to set the --pull flag for docker build; true by default
function build::docker_build() {
  local -r image=$1
  local -r context_dir=$2
  local -r pull="${3:-true}"
  local -ra build_cmd=("${DOCKER[@]}" build -t "${image}" "--pull=${pull}" "${context_dir}")

  log::status "Building Docker image ${image}"
  local docker_output
  docker_output=$("${build_cmd[@]}" 2>&1) || {
    cat <<EOF >&2
+++ Docker build command failed for ${image}

${docker_output}

To retry manually, run:

${build_cmd[*]}

EOF
    return 1
  }
}

function build::ensure_data_container() {
  # If the data container exists AND exited successfully, we can use it.
  # Otherwise nuke it and start over.
  local ret=0
  local code=0

  code=$(docker inspect \
      -f '{{.State.ExitCode}}' \
      "${APP_DATA_CONTAINER_NAME}" 2>/dev/null) || ret=$?
  if [[ "${ret}" == 0 && "${code}" != 0 ]]; then
    build::destroy_container "${APP_DATA_CONTAINER_NAME}"
    ret=1
  fi
  if [[ "${ret}" != 0 ]]; then
    log::status "Creating data container ${APP_DATA_CONTAINER_NAME}"
    # We have to ensure the directory exists, or else the docker run will
    # create it as root.
    mkdir -p "${LOCAL_OUTPUT_GOPATH}"
    # We want this to run as root to be able to chown, so non-root users can
    # later use the result as a data container.  This run both creates the data
    # container and chowns the GOPATH.
    #
    # The data container creates volumes for all of the directories that store
    # intermediates for the Go build. This enables incremental builds across
    # Docker sessions. The *_cgo paths are re-compiled versions of the go std
    # libraries for true static building.
    local -ra docker_cmd=(
      "${DOCKER[@]}" run
      --volume "${REMOTE_ROOT}"   # white-out the whole output dir
      --volume /usr/local/go/pkg/linux_386_cgo
      --volume /usr/local/go/pkg/linux_amd64_cgo
      --volume /usr/local/go/pkg/linux_arm_cgo
      --volume /usr/local/go/pkg/linux_arm64_cgo
      --volume /usr/local/go/pkg/linux_ppc64le_cgo
      --volume /usr/local/go/pkg/darwin_amd64_cgo
      --volume /usr/local/go/pkg/darwin_386_cgo
      --volume /usr/local/go/pkg/windows_amd64_cgo
      --volume /usr/local/go/pkg/windows_386_cgo
      --name "${APP_DATA_CONTAINER_NAME}"
      --hostname "${HOSTNAME}"
      "${APP_BUILD_IMAGE}"
      chown -R "${USER_ID}":"${GROUP_ID}"
        "${REMOTE_ROOT}"
        /usr/local/go/pkg/
    )
    "${docker_cmd[@]}"
  fi
}

# Build all app commands.
function build::build_command() {
  log::status "Running build command..."
  make -C "${APP_ROOT}" build.multiarch BINS="mgrid"
}
