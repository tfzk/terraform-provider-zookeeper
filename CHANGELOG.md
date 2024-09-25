## NEXT (MONTH DAY, 2024)

IMPROVEMENTS:

* Enabling CI testing for versions `1.9` of Terraform

NOTES:

* Moved build to [Golang `v1.23`](https://go.dev/blog/go1.23)
* Updated all dependencies to latest
* Updated [golangci-lint](https://golangci-lint.run/) linters
  * Removed `exportloopref`
  * Added `copyloopvar`
* Updated all GitHub Actions used in this repository:
  * [golangci/golangci-lint-action](https://github.com/golangci/golangci-lint-action) to `v6`
  * [goreleaser/goreleaser-action](https://github.com/goreleaser/goreleaser-action) to `v6`

## 1.1.0 (April 21, 2024)

IMPROVEMENTS:

* Enabling CI testing for versions `1.5`, `1.6`, `1.7` and `1.8` of Terraform

NOTES:

* Re-enabled [golangci-lint](https://golangci-lint.run/) linters, now that they are fully supported again:
  * `contextcheck` ([repo](https://github.com/kkHAIKE/contextcheck))
  * `wastedassign` ([repo](https://github.com/sanposhiho/wastedassign))
* Updated all dependencies to latest
* Moved build to [Golang `v1.22`](https://go.dev/blog/go1.22)
* Updated all GitHub Actions used in this repository:
  * [actions/checkout](https://github.com/actions/checkout) to `v4`
  * [actions/setup-go](https://github.com/actions/setup-go) to `v5`
  * [golangci/golangci-lint-action](https://github.com/golangci/golangci-lint-action) to `v4`
  * [hashicorp/setup-terraform](https://github.com/hashicorp/setup-terraform) to `v3`
  * [crazy-max/ghaction-import-gpg](https://github.com/crazy-max/ghaction-import-gpg) to `v6`
  * [goreleaser/goreleaser-action](https://github.com/goreleaser/goreleaser-action) to `v5`

## 1.0.4 (March 8, 2023)

NOTES:

* No new feature: updated dependencies (since `v1.0.3`).
* CI: Expanded testing matrix to include Terraform versions `1.2.*` to `1.4.*`.

## 1.0.3 (November 27, 2022)

NOTES:

* No new feature: updated dependencies (since `v1.0.2`).

## 1.0.2 (August 22, 2022)

NOTES:

* No new feature: updated [core dependencies](https://github.com/tfzk/terraform-provider-zookeeper/commit/f350b6cd70455c105636bd08f6169fd3743f0e36) (since `v1.0.1`), and moved build to [Golang `v1.18`](https://github.com/tfzk/terraform-provider-zookeeper/commit/f7451189924cc642adac9939f7d11f5610cc69db).

## 1.0.1 (July 27, 2022)

NOTES:

* No new feature: cutting a release because all major dependencies have received an update since release `v1.0.0`.
* Updated GH Action that imports the GPG key during release: [hashicorp/ghaction-import-gpg](https://github.com/hashicorp/ghaction-import-gpg#warning-this-action-as-been-deprecated) was deprecated.

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
