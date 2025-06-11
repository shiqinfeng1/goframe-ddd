# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
all: tidy gen cover build #lint  add-copyright 

# ==============================================================================
# Build options
ROOT_PACKAGE := $(shell go list -m)

# ==============================================================================
# Includes

include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/golang.mk
# include scripts/make-rules/image.mk
include scripts/make-rules/deploy.mk
# include scripts/make-rules/copyright.mk
include scripts/make-rules/gen.mk
include scripts/make-rules/release.mk
include scripts/make-rules/dependencies.mk
include scripts/make-rules/tools.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  BINS             要构建的二进制文件。默认是 cmd 目录下的所有文件。
				   此选项在使用 make build 或 make build.multiarch 时可用。
				   示例: make build BINS="mgrid tool"

  IMAGES           要制作的后端镜像。默认是 cmd 目录下所有文件。
				   此选项在使用 make image、make image.multiarch、make push 或 make push.multiarch 时可用。
				   示例: make image.multiarch IMAGES="app-apiserver app-authz-server"
						 
  REGISTRY_PREFIX  Docker 镜像仓库前缀。默认是 marmotedu。
				   示例: make push REGISTRY_PREFIX=ccr.ccs.tencentyun.com/marmotedu VERSION=v1.6.2

  PLATFORMS        要构建的多个平台。默认是 linux_amd64 和 linux_arm64。
				   此选项在使用 make build.multiarch、make image.multiarch 或 make push.multiarch 时可用。
				   示例: make image.multiarch IMAGES="app-apiserver app-pump" PLATFORMS="linux_amd64 linux_arm64"

  VERSION          编译到二进制文件中的版本信息。
				   默认从 gsemver 或 git 获取。

  V                Set to 1 enable verbose build. Default is 0.
				   设置为 1 可启用详细构建过程。默认值为 0。
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets

## build: Build source code for host platform.
.PHONY: build
build:
	@$(MAKE) go.build

## build.multiarch: Build source code for multiple platforms. See option PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch

## image: Build docker images for host arch.
.PHONY: image-nats
image-nats:
	docker build  -t nats-for-mgrid -f Dockerfile-nats .

.PHONY: image-arm
image-arm:
	docker build -f Dockerfile --build-arg TARGET=./cmd/mgrid --build-arg COMMIT=$(GIT_COMMIT)  --build-arg VERSION=$(VERSION) --build-arg CC=aarch64-linux-gnu-gcc --build-arg GOARCH=arm64 --target production \
		 --build-arg GOPRIVATE=$(GOPRIVATE) \
		-t "mgrid:$(VERSION)-arm"  .
# -t "$(CONTAINER_REGISTRY)/mgrid:$(IMAGE_TAG_ARM)"  .

## image.multiarch: Build docker images for multiple platforms. See option PLATFORMS.
.PHONY: image.multiarch
image.multiarch:
	@$(MAKE) image.build.multiarch

## push: Build docker images for host arch and push images to registry.
.PHONY: push
push:
	@$(MAKE) image.push

## push.multiarch: Build docker images for multiple platforms and push images to registry.
.PHONY: push.multiarch
push.multiarch:
	@$(MAKE) image.push.multiarch

## deploy: Deploy updated components to development env.
.PHONY: deploy
deploy:
	@$(MAKE) deploy.run

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## test: Run unit test.
.PHONY: test
test:
	@$(MAKE) go.test

## cover: Run unit test and get test coverage.
.PHONY: cover 
cover:
	@$(MAKE) go.test.cover

.PHONY: release.build
release.build:
	@$(MAKE) push.multiarch

## release: Release app
.PHONY: release
release:
	@$(MAKE) release.run

## verify-copyright: Verify the boilerplate headers for all files.
.PHONY: verify-copyright
verify-copyright:
	@$(MAKE) copyright.verify

## add-copyright: Ensures source code files have copyright license headers.
.PHONY: add-copyright
add-copyright:
	@$(MAKE) copyright.add

## gen: Generate all necessary files, such as error code files.
.PHONY: gen
gen:
	@$(MAKE) gen.run

## swagger: Generate swagger document.
.PHONY: swagger
swagger:
	@$(MAKE) swagger.run

## serve-swagger: Serve swagger spec and docs.
.PHONY: swagger.serve
serve-swagger:
	@$(MAKE) swagger.serve

## dependencies: Install necessary dependencies.
.PHONY: dependencies
dependencies:
	@$(MAKE) dependencies.run

## tools: install dependent tools.
.PHONY: tools
tools:
	@$(MAKE) tools.install

## check-updates: Check outdated dependencies of the go projects.
.PHONY: check-updates
check-updates:
	@$(MAKE) go.updates

.PHONY: tidy
tidy:
	@$(MAKE) dependencies.packages



## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
