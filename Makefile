# NOTE: This is valid only while the `tools/zk-local-ensemble` is up
ZOOKEEPER_SERVERS="localhost:2181,localhost:2182,localhost:2183"

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
	go test -v -cover -timeout=120s -parallel=4 ./...

testacc:
	ZOOKEEPER_SERVERS=$(ZOOKEEPER_SERVERS) \
	TF_ACC=1 \
		go test -v -cover -parallel=4 -timeout 2m ./...

.PHONY: build install lint generate fmt test testacc
