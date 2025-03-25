
SHELL := /bin/bash

# include the common make file
# MAKEFILE_LIST 是 Makefile 中的一个自动变量。
# 它是一个包含所有已经被 Make 读取的 Makefile 文件名称的列表，
# 这些文件包括通过 include 指令引入的其他 Makefile 文件。列表中的文件名按读取的顺序排列，最新读取的文件位于列表末尾。
# lastword函数会返回 MAKEFILE_LIST 列表中的最后一个文件名，也就是当前正在处理的 Makefile 文件的名称
# dir 函数会返回该文件所在的目录路径。
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

# origin 函数可能返回的结果包括 undefined（未定义）、default（默认定义）、environment（从环境变量中获取）、file（在 Makefile 文件中定义）等
# 进入 COMMON_SELF_DIR 目录向上两层的目录，并打印当前工作目录的绝对路径，abspath将参数转换为绝对路径
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif
# 如果不存在_output
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif
ifeq ($(origin TOOLS_DIR),undefined)
TOOLS_DIR := $(OUTPUT_DIR)/tools
$(shell mkdir -p $(TOOLS_DIR))
endif
ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))
endif

# set the version number. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif


# 跨平台编译时，指定需要编译的平台，默认是只编译当前环境所在的平台
# The OS must be linux when building docker images
PLATFORMS ?= linux_amd64 linux_arm64
# The OS can be linux/windows/darwin when building binaries
# PLATFORMS ?= darwin_amd64 windows_amd64 linux_amd64 linux_arm64

# Set a specific PLATFORM
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	# Use linux as the default OS when building images
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

# Linux command settings
# FIND := find . ! -path './third_party/*' ! -path './vendor/*'
ifeq ($(GOOS),linux)
	XARGS := xargs --no-run-if-empty
endif
ifeq ($(GOOS),darwin)
	XARGS := xargs
endif


# Makefile settings
ifndef V
MAKEFLAGS += --no-print-directory
endif

# Copy githook scripts when execute makefile
# COPY_GITHOOK:=$(shell cp -f scripts/githooks/* .git/hooks/)

# Specify tools severity, include: BLOCKER_TOOLS, CRITICAL_TOOLS, TRIVIAL_TOOLS.
# Missing BLOCKER_TOOLS can cause the CI flow execution failed, i.e. `make all` failed.
# Missing CRITICAL_TOOLS can lead to some necessary operations failed. i.e. `make release` failed.
# TRIVIAL_TOOLS are Optional tools, missing these tool have no affect.
BLOCKER_TOOLS ?= go-junit-report golangci-lint gsemver gf # golines goimports codegen addlicense
CRITICAL_TOOLS ?= go-gitlint go-mod-outdated mockgen git-chglog github-release protoc protoc-gen-go protoc-gen-go-grpc cc-arm64 #swagger gotests   protoc-gen-go-http protoc-gen-go-errors protoc-gen-validate protoc-gen-openapi  wire coscmd 
TRIVIAL_TOOLS ?= #depth go-callvis gothanks richgo rts kube-score

COMMA := ,
SPACE :=
SPACE +=
