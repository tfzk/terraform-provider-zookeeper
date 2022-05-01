# NOTE: This is valid only while the `make local.zk.up`
ZOOKEEPER_SERVERS=localhost:2181,localhost:2182,localhost:2183

default: build

build:
	go build -v ./...

install: build
	go install -v ./...

# See https://golangci-lint.run/
lint:
	golangci-lint run

generate:
	go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=2m -parallel=4 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout=2m -parallel=4 ./...

local.zk.up:
	./scripts/zk-local-ensemble/up

local.zk.down:
	./scripts/zk-local-ensemble/down

local.zk.restart:
	./scripts/zk-local-ensemble/restart

local.testacc:
	ZOOKEEPER_SERVERS=$(ZOOKEEPER_SERVERS) make testacc

.PHONY: build install lint generate fmt test testacc local.zk.up local.zk.down local.zk.restart local.testacc
