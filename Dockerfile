
# NOTE: Must be run in the context of the repo's root directory

####################################
## (1) Setup the build environment
FROM golang:1.24-bullseye AS build-setup

RUN sed -i 's/deb.debian.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list && \
    sed -i 's/security.debian.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list
RUN apt-get update
RUN apt-get -y install  gcc-aarch64-linux-gnu

## (2) Setup crypto dependencies
FROM build-setup AS build-env

# Build the app binary in /app
RUN mkdir /app
WORKDIR /app

ARG TARGET
ARG COMMIT
ARG VERSION

ENV GOPRIVATE=

COPY . .

####################################
## (3) Build the production app binary
FROM build-env AS build-production
WORKDIR /app

ARG GOARCH=amd64
# TAGS can be overriden to modify the go build tags (e.g. build without netgo)
ARG TAGS="netgo,osusergo"
# CC flag can be overwritten to specify a C compiler
ARG CC=""
# CGO_FLAG uses ADX instructions by default, flag can be overwritten to build without ADX
ARG CGO_FLAG=""

# 后续可添加其他构建步骤，例如设置 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

# Keep Go's build cache between builds.
# https://github.com/golang/go/issues/27719#issuecomment-514747274
RUN --mount=type=cache,sharing=locked,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux GOARCH=${GOARCH} CC="${CC}" CGO_CFLAGS="${CGO_FLAG}" go build --tags "${TAGS}" -ldflags "-s -w -extldflags -static \
    -X 'github.com/shiqinfeng1/goframe-ddd/pkg/version.GitCommit=${COMMIT}' -X  'github.com/shiqinfeng1/goframe-ddd/cmd/pkg/version.GitVersion=${VERSION}'" \
    -o ./app ${TARGET}

RUN chmod a+x /app/app

## (4) Add the statically linked production binary to a distroless image
FROM gcr.io/distroless/base-debian11 AS production

COPY --from=build-production /app/app /bin/app

ENTRYPOINT ["/bin/app"]