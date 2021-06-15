# Copyright PeerFintech All Rights Reserved.
#
# -------------------------------------------------------------
# This makefile defines the following targets
#   - apiserver - builds a apiserver binary
#   - unit-test - runs the go-test based unit tests
#   - apiserver-docker - build the apiserver image

#apiserver每个项目默认从1.0.0开始
BASE_APISERVER_VERSION = 1.0.0
PREV_APISERVER_VERSION = 1.0.0


PROJECT_NAME = jzsg/fabric-explorer
DOCKERIMAGE = jzsg/fabric-explorer

export DOCKERIMAGE
EXTRA_VERSION ?= $(shell git rev-parse --short HEAD)
PROJECT_VERSION=$(BASE_APISERVER_VERSION)-snapshot-$(EXTRA_VERSION)

PKGNAME = github.com/$(PROJECT_NAME)
CGO_FLAGS = CGO_CFLAGS=" "
ARCH=$(shell go env GOARCH)
MARCH=$(shell go env GOOS)-$(shell go env GOARCH)
IMAGEIDS= $(shell docker images --quiet --filter=reference='$(DOCKERIMAGE)')
# defined in apiserver/common/metadata.go
METADATA_VAR = Version=$(BASE_APISERVER_VERSION)
METADATA_VAR += CommitSHA=$(EXTRA_VERSION)

GO_TAGS ?=


export GO_LDFLAGS GO_TAGS

.PHONY: apiserver
apiserver:
	@echo "Building apiserver...."
	go build -o ./bin/apiserver -mod=vendor -ldflags "-extldflags -static" $(PKGNAME)/apiserver

.PHONY: apiserver-docker
apiserver-docker: apiserver
	@echo "Building apiserver docker image...."
	@./buildImage.sh

.PHONY: unit-test
unit-test:
	@echo "unit test...."
	@MODULE="github.com/peerfintech/Hsbank" \
	PKG_ROOT="./" \
	./scripts/unit.sh

.PHONY: clean
clean:
	-@rm -rf ./bin
ifneq ($(IMAGEIDS),)
	    -docker rmi -f $(IMAGEIDS)
endif
