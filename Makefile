TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=github.com
NAMESPACE=tfzk
NAME=zookeeper-legacy
BINARY=bin/terraform-provider-${NAME}
VERSION=0.0.0-dev
OS_ARCH=darwin_arm64

# NOTE: This is valid only while the `tools/zk-local-ensemble` is up
ZOOKEEPER_SERVERS?=localhost:2181,localhost:2182,localhost:2183

default: install

build:
	go build -o ${BINARY}

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	ZOOKEEPER_SERVERS=$(ZOOKEEPER_SERVERS) \
	TF_ACC=true \
		go test $(TEST) -v $(TESTARGS) -timeout 120m