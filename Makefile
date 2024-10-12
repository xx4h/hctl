BINDIR				:= $(CURDIR)/bin
DISTDIR				:= $(CURDIR)/dist
INSTALL_PATH	?= /usr/local/bin
DIST_DIRS			:= find * -type d -exec
BINNAME				?= hctl
LOCAL_INSTALL_PATH := $(HOME)/.local/bin

GOBIN					= $(shell go env GOBIN)
ifeq ($(GOBIN),)
	GOBIN 			= $(shell go env GOPATH)/bin
endif
GOIMPORTS			= $(GOBIN)/goimports
ARCH					= $(shell go env GOARCH)

GIT_COMMIT		= $(shell git rev-parse HEAD)
GIT_SHA				= $(shell git rev-parse --short HEAD)
GIT_TAG				= $(shell git describe --tags --abrev=0 --exact-match 2>/dev/null)
GIT_DIRTY			= $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

# go option
PKG						:= ./...
TAGS					:=
TESTS					:= .
TESTFLAGS			:=
LDFLAGS				:= -w -s\
								 -X github.com/xx4h/hctl/cmd.version=$(shell git rev-parse --abbrev-ref HEAD)-$(GIT_DIRTY)\
								 -X github.com/xx4h/hctl/cmd.commit=$(GIT_COMMIT)\
								 -X github.com/xx4h/hctl/cmd.date=$(shell date -Iseconds)
GOFLAGS				:=
CGO_ENABLED		?= 0

# rebuild binary if any of these files change
SRC						:= $(shell find . -type f -name '*.go' -print) go.mod go.sum

# required for flobs to work correctly
SHELL					= /usr/bin/env bash

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif

.PHONY: all
all: build

# ---
# build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) ./main.go

# ------------------------------------------------------------------------------
#  install

.PHONY: install
install: build
	@install "$(BINDIR)/$(BINNAME)" "$(INSTALL_PATH)/$(BINNAME)"

# ------------------------------------------------------------------------------
#  local-install

.PHONY: local-install
local-install: build
	@install "$(BINDIR)/$(BINNAME)" "$(LOCAL_INSTALL_PATH)/$(BINNAME)"

# ------------------------------------------------------------------------------
#  test

.PHONY: test
test: build
ifeq ($(ARCH),s390x)
test: TESTFLAGS += -v
else
test: TESTFLAGS += -race -v
endif
test: test-style

.PHONY: test-unit
test-unit:
	@echo
	@echo "==> Running unit tests <=="
	go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)
	@echo
	@echo "==> Running unit test(s) with ldflags <=="
	go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS) -ldflags '$(LDFLAGS)'

.PHONY: test-style
test-style:
	golangci-lint run ./...

.PHONY: test-goreleaser
test-goreleaser:
	goreleaser release --snapshot

.PHONY: format
format: $(GOIMPORTS)
	go list -f '{{.Dir}}' ./... | xargs $(GOIMPORTS) -w -local github.com/xx4h/hctl

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)' '$(DISTDIR)'

.PHONY: info
info:
	 @echo "Version:           ${VERSION}"
	 @echo "Git Tag:           ${GIT_TAG}"
	 @echo "Git Commit:        ${GIT_COMMIT}"
	 @echo "Git Tree State:    ${GIT_DIRTY}"
