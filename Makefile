# NOTE: This is valid only while the `make local.zk.up`
LOCAL_ZOOKEEPER_SERVERS=localhost:2181,localhost:2182,localhost:2183

default: build

build:
	go build -v ./...

install: build
	go install -v ./...

# Executes golangci-lint.
# See: https://golangci-lint.run/.
lint:
	golangci-lint run

# Generates the documentation that eventually gets published here:
# https://registry.terraform.io/providers/tfzk/zookeeper/latest/docs.
generate:
	go generate ./...

# Formats the codebase.
fmt:
	gofmt -s -w -e .

# Updates all dependencies, recursively.
deps.update:
	go get -u ./...
	go mod tidy

test:
	go test -v -cover -timeout=2m -parallel=4 ./...

# Stands up a ZooKeeper Ensemble, for testing.
local.zk.up:
	./scripts/zk-local-ensemble/up

# Shuts down the ZooKeeper Ensemble.
local.zk.down:
	./scripts/zk-local-ensemble/down

# Restarts the ZooKeeper Ensemble.
local.zk.restart:
	./scripts/zk-local-ensemble/restart

# Runs Acceptance Tests against the ZooKeeper Ensemble running locally.
local.test:
	ZOOKEEPER_SERVERS=$(LOCAL_ZOOKEEPER_SERVERS) TF_ACC=1 make test

.PHONY: build install lint generate fmt deps.update test local.zk.up local.zk.down local.zk.restart local.test
