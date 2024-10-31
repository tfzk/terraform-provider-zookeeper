# TFZK: Terraform provider for ZooKeeper

[![Build and Test](https://github.com/tfzk/terraform-provider-zookeeper/actions/workflows/build-test.yml/badge.svg)](https://github.com/tfzk/terraform-provider-zookeeper/actions/workflows/build-test.yml)
[![Terraform Provider ZooKeeper Documentation](https://img.shields.io/badge/terraform-%235835CC.svg?style=flat&logo=terraform&logoColor=white&label=docs&labelColor=305)](https://registry.terraform.io/providers/tfzk/zookeeper/latest/docs)

A Terraform provider for [Apache ZooKeeper](https://zookeeper.apache.org/).

To get started, go
on [`terraform-provider-zookeeper` registry page](https://registry.terraform.io/providers/tfzk/zookeeper).

This provider focuses on _Data Management_ for ZooKeeper:
it offers a _CRUD_ for **Persistent ZNodes** and **Persistent Sequential ZNodes**.

For example, it can be leveraged as part of a bigger infrastructure, where sharing data across multiple
live services via ZooKeeper is desirable. Good examples can be _runtime configuration data_ or
_large architectures topology data_ and so forth.

Data can be stored both as UTF-8 and binary (via Base64 encoding) inside ZooKeeper
[ZNodes](https://zookeeper.apache.org/doc/r3.1.2/zookeeperProgrammers.html#sc_zkDataModel_znodes).

## Compatibility

Compatibility table between this provider,
the [Registry Protocol](https://www.terraform.io/internals/provider-registry-protocol)
version it implements, and Terraform:

| Provider | Registry Protocol | Terraform |
|:--------:|:-----------------:|:---------:|
| `>= 1.x` |        `5`        | `>= 0.12` |

### CI Testing

The provider test suite is run against all Terraform versions from `0.12` to `1.9`,
as well as all ZooKeeper versions from `3.5` to `3.9`. 

See the [Build and Test](https://github.com/tfzk/terraform-provider-zookeeper/blob/main/.github/workflows/build-test.yml)
workflow for details.

## Provider features

* [x] support for ZK standard multi-server connection string
* [x] support for ZK authentication
* [x] support for ZK ACLs
* [x] "session timeout" configuration
* [x] create ZNode
* [x] create Sequential ZNode
* [x] read ZNode
* [x] update ZNode
* [x] delete ZNode
* [x] import ZNode
* [x] import Sequential ZNode
* [x] support for binary data in Base64 format

## Development

### Requirements

* [Go](https://go.dev/dl/) >= `1.23`
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [golangci-lint](https://golangci-lint.run/)

### Run acceptance tests locally

To run acceptance tests, you will need a ZooKeeper Ensemble running:

```shell
$ make local.zk.up

$ make local.test

# ... do your development / fixing ...

$ make local.zk.down
```
In `scripts/zk-local-ensemble` we provide a `docker-compose.yml` that can spin
up an ensemble made of 3 servers, running on `localhost` ports `2181, 2182 and 2183`.
Everything can be controlled via the `make local.*` commands provided.

If you are curious, please take a look at the `Makefile` to understand how those are then passed to
go during (Acceptance) Tests.

## License

All the content of this repository is under [MIT License](./LICENSE)
