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

The provider is tested against all combinations of:

* **Terraform** - versions `0.12` to `1.15`
* **ZooKeeper** - versions `3.3` to `3.9` 

See the [Build and Test](https://github.com/tfzk/terraform-provider-zookeeper/blob/main/.github/workflows/build-test.yml)
workflow for details.

## Provider features

* [x] support for ZK standard multi-server connection string
* [x] support for ZK authentication
* [x] support for TLS and mTLS
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

* [Go](https://go.dev/dl/) >= `1.27`
* [Task](https://taskfile.dev/)
* [golangci-lint](https://golangci-lint.run/)
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)

#### [asdf](https://asdf-vm.com/)

The project comes pre-configured with `.tool-versions`: if you already use [asdf](https://asdf-vm.com/),
just run `asdf install` to get all required development tools.

### Testing locally

#### One-shot acceptance tests
To run the full suite of acceptance tests automatically (spinning up the local 3-server ZooKeeper ensemble, running all tests, and shutting down the ensemble afterwards), execute:

```shell
$ task local.test
```

#### Running `task test` directly
By default, running:

```shell
$ task test
```

will run unit tests and **skip any acceptance tests** (so no running ZooKeeper cluster is required).

#### Iterative acceptance testing against a running ensemble
If you want to run acceptance tests iteratively during development using `task test`:

1. Start the local ZooKeeper ensemble:
   ```shell
   $ task local.zk.up
   ```

2. Run tests with `TF_ACC=1`:
   ```shell
   $ TF_ACC=1 task test
   # or export TF_ACC=1 and then run: task test
   ```

3. When finished developing, shut down the local ensemble:
   ```shell
   $ task local.zk.down
   ```

> [!NOTE]
> You can provide your own ZooKeeper ensemble by setting the `ZOOKEEPER_SERVERS` environment variable to a comma-separated list
> of `host:port` pairs (e.g. `export ZOOKEEPER_SERVERS="localhost:2181,localhost:2182,localhost:2183"`).
> In that case, just do step 2 above.

In `scripts/zk-local-ensemble` we provide a `docker-compose.yml` that spins up an ensemble made of 3 servers running on `localhost` ports `2181, 2182, and 2183`. Everything can be controlled via the `task local.zk.*` commands provided.

If you are curious, please take a look at `Taskfile.yml` to understand how variables and environment configurations are passed to Go during tests.

## License

All the content of this repository is under [MIT License](./LICENSE)
