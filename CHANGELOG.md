## 1.0.0 (July 11, 2022)

NEW FEATURES:

* provider: attribute `servers` supports configuration for a full ensemble of ZooKeeper servers, comma-separated.
* provider: attribute `session_timeout` supports configuration for the internal ZooKeeper client session timeout

* data-source/zookeeper_znode: Reads `data` and `stat` of a ZNode located at `path`; binary content is accessible via `data_base64`.

* resource/zookeeper_znode: Full CRUD for a ZNode located at `path`. Fields `data` and `data_base64` provide access to the content, `stat` to the ZNode [stat structure](https://registry.terraform.io/providers/tfzk/zookeeper/latest/docs#the-stat-structure).
* resource/zookeeper_znode: Support for `import`ing.

* resource/zookeeper_sequential_znode: Same as `zookeeper_znode`, though the full path is computed starting from an initial `path_prefix`.
* resource/zookeeper_sequential_znode: Support for `import`ing.

NOTES:

* Provider is based on [terraform-plugin-sdk](https://github.com/hashicorp/terraform-plugin-sdk/) `v2`: this makes it compatible with Terraform `>= 0.12`.
* Finalised tooling and configuration for static code analysis, using [golangci-lint](https://golangci-lint.run/).
* Updated all dependencies since previous release `v1.0.0-alpha`.

## 1.0.0-alpha (June 05, 2022)

NOTES:

* feature complete for `v1.0.0`, but not releasing it yet
* testing release process before `v1.0.0`
* `CHANGELOG.md` not ready yet: creating first entry for `v1.0.0-alpha`
