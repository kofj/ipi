COMMA:=,
SHELL:=/bin/bash
# colors
RED:=\033[0;31m
GREEN:=\033[0;32m
BLUE:=\033[0;34m
YELLOW:=\033[0;33m
NC:=\033[0m

# include the common make file
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif

ifeq ($(origin VERSION), undefined)
VERSION ?= $(shell git describe --match 'v[0-9]*' --dirty='-dirty' --tags --always)
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)
BUILD_DATE:=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

PLATFORMS ?= linux_amd64 linux_arm64
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
else
    GOOS := $(word 1,$(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2,$(subst _, ,$(PLATFORM)))
endif

