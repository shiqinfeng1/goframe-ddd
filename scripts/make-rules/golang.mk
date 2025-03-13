
# ==============================================================================
# Makefile helper functions for golang
#

GO := go
VERSION_PACKAGE=$(ROOT_PACKAGE)/pkg/version
GIT_COMMIT:=$(shell git rev-parse HEAD)
# Check if the tree is dirty.  default to dirty
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif

GO_LDFLAGS += -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
ifneq ($(DLV),)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
	LDFLAGS = ""
endif
GO_BUILD_FLAGS += -ldflags "-s -w $(GO_LDFLAGS)"

ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif
ifeq ($(GOOS),linux)
	SED := sed -i
endif
ifeq ($(GOOS),darwin)
	SED := sed -i ''
	UPX_FLAG := --force-macos
endif
ifeq ($(GOARCH),arm64)
	CC := aarch64-linux-gnu-gcc
endif

ifeq ($(ROOT_PACKAGE),)
	$(error the variable ROOT_PACKAGE must be set prior to including golang.mk)
endif

GOPATH := $(shell go env GOPATH)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

# 先取出所有命令的main.go所在的文件夹路径，遍历COMMANDS，以cmd下的文件夹名称作为bin的名称
COMMANDS ?= $(filter-out %.md, $(wildcard ${ROOT_DIR}/cmd/*))
BINS ?= $(foreach cmd,${COMMANDS},$(notdir ${cmd}))

ifeq (${COMMANDS},)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif
ifeq (${BINS},)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

# 测试报告过滤的文件夹
EXCLUDE_TESTS=github.com/shiqinfeng1/goframe-ddd/api'|'github.com/shiqinfeng1/goframe-ddd/cmd'|'github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent

# Minimum test coverage
ifeq ($(origin COVERAGE_TARGET——TARFGET),undefined)
COVERAGE_TARGET := 60
endif


.PHONY: go.build.verify
go.build.verify:

# 执行的具体的编译操作
.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@mkdir -p $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)
	CGO_ENABLED=1 GOOS=$(OS) GOARCH=$(ARCH) CC=$(CC) $(GO) build $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_PACKAGE)/cmd/$(COMMAND)
	upx $(UPX_FLAG) $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT)

# 编译当前平台的bin
.PHONY: go.build
go.build: go.build.verify $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))

# 编译跨平台的bin
.PHONY: go.build.multiarch
go.build.multiarch: go.build.verify $(foreach p,$(PLATFORMS),$(addprefix go.build., $(addprefix $(p)., $(BINS))))

# 清空所有编译输出
.PHONY: go.clean
go.clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

# 静态代码分析工具
.PHONY: go.lint
go.lint: tools.verify.golangci-lint
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

# 执行单元测试用例，输出每个函数
.PHONY: go.test
go.test: tools.verify.go-junit-report
	@echo "===========> Run unit test"
# 开启并发竞态检测，找出代码中可能存在的并发问题
# 打乱测试用例的执行顺序，增强测试的随机性
# 设置测试超时时间，避免测试无限期运行
# 运行简短测试，跳过一些耗时较长的测试用例
# 把EXCLUDE_TESTS列表中的空格替换为|（在正则表达式中|表示或），再过滤掉当前项目下所有Go包中在EXCLUDE_TESTS中的包
# 执行测试用例，并生成的代码覆盖率数据文件coverage.out
# 将测试结果同时输出到终端和 JUnit 格式的 XML 报告文件
	@set -o pipefail;$(GO) test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out \
		-timeout=10m -shuffle=on -short -v `go list ./...|egrep -v $(EXCLUDE_TESTS)` 2>&1 | \
		tee >(go-junit-report --set-exit-code >$(OUTPUT_DIR)/report.xml)
# remove mock_.*.go files from test coverage
	@$(SED) '/api/d' $(OUTPUT_DIR)/coverage.out 
	@$(SED) '/server/d' $(OUTPUT_DIR)/coverage.out 
	@$(SED) '/adapters/d' $(OUTPUT_DIR)/coverage.out 
# 根据指定的代码覆盖率数据文件coverage.out生成 HTML 格式的报告
	@$(GO) tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html

# 查看代码覆盖率统计信息, 执行go tool cover命令后，会显示每个函数的名称、执行的语句数量以及该函数的代码覆盖率百分比
# 执行coverage.awk 对覆盖率镜像统计和判断
.PHONY: go.test.cover
go.test.cover: go.test
	@$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out | \
		awk -v target=$(COVERAGE_TARGET) -f $(ROOT_DIR)/scripts/coverage.awk

# 检查依赖包是否有更新，比较耗时
.PHONY: go.updates
go.updates: tools.verify.go-mod-outdated
	@$(GO) list -u -m -json all | go-mod-outdated -update -direct
