# ================================================================
# project custom variables
IMAGE=ghcr.io/kofj/ipi
ROOT_PACKAGE=github.com/kofj/ipi
VERSION_PACKAGE := $(ROOT_PACKAGE)/pkg/version
# git version
GIT_VERSION := $(shell git describe --tags --always --dirty)

# ================================================================
# Includes makefiles
include build/lib/common.mk
include build/lib/golang.mk

# ================================================================
# go build info
.PHONY: build.info
build.info:
	@$(MAKE) go.build.info

## build: Build source code for host platform.
.PHONY: build
build:
	@echo -e "${GREEN} build binary for ${PLATFORM} ${NC}"
	@$(MAKE) go.build.multiarch

image:
	@echo -e "${GREEN} build image $(IMAGE):${GIT_VERSION} ${NC}"
	@docker build -t $(IMAGE):${GIT_VERSION} -f Dockerfile .
