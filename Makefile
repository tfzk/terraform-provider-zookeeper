HOSTNAME=github.com
NAMESPACE=tfzk
NAME=zookeeper
BINARY=bin/terraform-provider-${NAME}

# Terraform requires providers to be sem-versioned.
# The format allows to append a string like `-dev`.
# So, we use version `0.0.0-dev` during development and testing,
# and then leave it to `goreleaser` to use the git tag to establish
# the actual release version.
VERSION=0.0.0-dev

# Detect Operating System and Architecture
ifeq ($(OS),Windows_NT)
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		OS_ARCH="windows_amd64"
	else ifeq ($(PROCESSOR_ARCHITECTURE),ARM64)
		OS_ARCH="windows_arm64"
	else ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		OS_ARCH="windows_386"
	endif
else
	OS_ARCH=$(shell uname -s | tr A-Z a-z )_$(shell uname -m)
endif

# NOTE: This is valid only while the `tools/zk-local-ensemble` is up
ZOOKEEPER_SERVERS?="localhost:2181,localhost:2182,localhost:2183"

default: build

build:
	go build -o ${BINARY} -v .

fmt:
	gofmt -s -w -e .

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -cover -timeout=60s -parallel=4 ./...

testacc: 
	ZOOKEEPER_SERVERS=$(ZOOKEEPER_SERVERS) \
	TF_ACC=true \
		go test -v -cover -timeout=60s ./...

.PHONY: build fmt release install test testacc
