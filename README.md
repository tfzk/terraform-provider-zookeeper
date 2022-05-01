# TFZK Terraform provider for ZooKeeper

A Terraform provider for [Apache ZooKeeper](https://zookeeper.apache.org/)
that can be leveraged as part of a bigger infrastructure that depends on having data stored in ZooKeeper
[ZNodes](https://zookeeper.apache.org/doc/r3.1.2/zookeeperProgrammers.html#sc_zkDataModel_znodes).

**NOTE:** This provider is built on top of
[Terraform Plugin SDKv2](https://www.terraform.io/plugin/sdkv2/sdkv2-intro),
and as such can be used with versions of Terraform `>= 0.12` (legacy).
It will work just fine with Terraform `>= 1.0`, but it's purpose is to
keep supporting users of legacy versions of Terraform.

## Provider features

* [x] support for ZK standard multi-server connection string
* [ ] support for ZK authentication
* [ ] support for ZK ACLs
* [x] "session timeout" configuration
* [x] create ZNode
* [x] create Sequential ZNode
* [x] read ZNode
* [x] update ZNode
* [x] delete ZNode
* [x] import ZNode
* [x] import Sequential ZNode

## Project 1.0.0 must haves

* [ ] Documentation
  * [x] setup `tfplugindocs` for autogeneration
  * [ ] provide content etc.
* [x] Build and Test automation
  * [x] `golangci-lint`
  * [x] verify that generated documentation is up-to-date
  * [x] builds cleanly
  * [ ] Spin up ZooKeeper service (dependency for acceptance testing)
  * [ ] Acceptance testing against all latest minor releases of Terraform >= 0.12
* [ ] Release automation
  * [ ] Triggered by semver tag detected
  * [ ] Generates `CHANGELOG` entry automatically
  * [ ] Publishes to registry.terraform.com

## Development

TBD

### Requirements

* [Go](https://go.dev/dl/) >= `1.17`
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)

### Ensemble to run tests against

To run tests, you will need a ZooKeeper Ensemble running.

In `tools/zk-local-ensemble` we provide a `docker-compose.yml` that can spin
up an ensemble made of 3 servers, running on `localhost` ports `2181, 2182 and 2183`.

Please take a look at the `Makefile` to understand how those are then passed to
go during (Acceptance) Tests.

## License

All the content of this repository is under [MIT License](./LICENSE)
