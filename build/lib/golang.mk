GO := go
GO_SUPPORTED_VERSIONS ?= 1.18|1.20|1.22|1.24
# set root package path
ifeq ($(origin ROOT_PACKAGE), undefined)
ROOT_PACKAGE=git.woa.com/tcr-cloud/base
endif
# set version package path
ifeq ($(origin VERSION_PACKAGE), undefined)
VERSION_PACKAGE := $(ROOT_PACKAGE)/pkg/version
endif
GO_LDFLAGS += -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).BuildDate=$(BUILD_DATE) \

COMMANDS ?= $(filter-out %.md, $(wildcard ${ROOT_DIR}/cmd/*))
BINS ?= $(foreach cmd,${COMMANDS},$(notdir ${cmd}))
ifeq (${COMMANDS},)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif
ifeq (${BINS},)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

.PHONY: go.lint.verify
go.lint.verify:
ifeq (,$(shell which golangci-lint))
	@echo -e "${RED}install golangci-lint${NC}"
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4
	@golangci-lint version
endif

.PHONY: go.air.verify
go.air.verify:
ifeq (,$(shell which air))
	@echo "${RED}install air to ${NC}"
	@GO111MODULE=off $(GO) install github.com/air-verse/air@latest
	@air --version
endif

.PHONY: go.lint
go.lint: go.lint.verify
	@echo -e "${BLUE}linting code...${NC}"
	@golangci-lint run --verbose --timeout 5m ./...

.PHONY: gomv
go.mv:
	@echo checking go mod tidy...
	@go mod tidy
	@go mod vendor

.PHONY: go.build.verify
go.build.verify:
ifneq ($(shell $(GO) version | grep -q -E '\bgo($(GO_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	@echo -e "${RED}Unsupported go version. Please make install one of the following supported version: '$(GO_SUPPORTED_VERSIONS)'${NC}"
else
	@echo -e "${GREEN}Go version check passedm${NC}"
endif

.PHONY: go.build.info
go.build.info:
	@echo -e "${YELLOW}version: ${VERSION}${NC}"
	@echo -e "${BLUE}ROOT_DIR: $(ROOT_DIR)${NC}"
	@echo -e "${BLUE}ROOT_PACKAGE: $(ROOT_PACKAGE)${NC}"
	@echo -e "${BLUE}VERSION_PACKAGE: $(VERSION_PACKAGE)${NC}"
	@echo -e "${BLUE}GO_LDFLAGS: $(GO_LDFLAGS)${NC}"
	@echo -e "${BLUE}COMMANDS: $(COMMANDS)${NC}"
	@echo -e "${GREEN}BINS: $(BINS)${NC}"

.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo -e "${YELLOW}==>Building binary ${GREEN}$(COMMAND) version:$(VERSION) for $(OS)/$(ARCH)${NC}"
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build -o $(OUTPUT_DIR)/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) -ldflags "$(GO_LDFLAGS)" $(ROOT_PACKAGE)/cmd/$(COMMAND)

.PHONY: go.build
go.build: go.build.verify go.build.info go.lint $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))

.PHONY: go.build.multiarch
go.build.multiarch: go.build.verify $(foreach p,$(PLATFORMS),$(addprefix go.build., $(addprefix $(p)., $(BINS))))
