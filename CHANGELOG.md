
IMPROVEMENTS:

* Enabling CI testing for versions `1.14` of Terraform

## 1.2.10 (Oct 22, 2025)

NOTES:

* Updated repository Go version to `1.25.3`
  * See the [Go `1.25.3` milestone](https://github.com/golang/go/issues?q=milestone%3AGo1.25.3)
  * See the [Go `1.25.3` release notes](https://go.dev/doc/devel/release#go1.25.3)
* Updated all dependencies to latest

IMPROVEMENTS:

* Minor tweak to `Makefile` to ensure building and dependency updates commands
  always include test code.

## 1.2.9 (Sep 23, 2025)

NOTES:

* Updated repository Go version to `1.25.1`
  * See the [Go `1.25.1` milestone](https://github.com/golang/go/issues?q=milestone%3AGo1.25.1)
  * See the [Go `1.25` milestone](https://github.com/golang/go/issues?q=milestone%3AGo1.25)
* Updated all dependencies to latest
* Updated repository `golangci-lint` to `2.5.0`

IMPROVEMENTS:

* Enabling CI testing for versions `1.12` and `1.13` of Terraform
* Enabling CI testing for versions `3.3` and `3.4` of ZooKeeper

## 1.2.8 (Jul 12, 2025)

NOTES:

* Updated repository Go version to `1.24.5`
  * See the [Go `1.24.5` milestone](https://github.com/golang/go/issues?q=milestone%3AGo1.24.5) on our issue tracker for details.
* Updated all dependencies to latest
  * Includes updates to `github.com/hashicorp/terraform-plugin-docs` to `v0.22.0`, with minor tweaks to docs in the `Import` section

IMPROVEMENTS:

* Execute `tfplugindocs validate` as part of `//go:generate` step for the documentation 

## 1.2.7 (Jun 19, 2025)

NOTES:

* Updated repository Go version to `1.24.4`
* Updated repository `golangci-lint` version to `v2` (`2.1.6`)
* Updated all dependencies to latest
* Added note in README about [asdf](https://asdf-vm.com/)
* Address 2 (duplicate) code scanning issues, related to `build-test.yml` workflow missing permission restriction
  * https://github.com/tfzk/terraform-provider-zookeeper/security/code-scanning/2
  * https://github.com/tfzk/terraform-provider-zookeeper/security/code-scanning/1

## 1.2.6 (May 02, 2025)

NOTES:

* Updated repository Go version to `1.24.2`
  * Adding `toolchain` hint in `go.mod`
* Updated all dependencies to latest
* Updated repository `golangci-lint` version to `v2` (`2.1.5`)
  * Now used for both linting _and_ formatting
* Updated GitHub Action `golangci/golangci-lint-action` from `v6` to `v7`

IMPROVEMENTS:

* Switching to `golangci-lint@v2` highlighted places where internal doc and formatting could improved, so did that
* `Makefile` `fmt`, `lint` and `lint-fix` now solely rely on `golangci-lint`

## 1.2.5 (March 14, 2025)

NOTES:

* Updated repository Go version to `1.24.1`
  * Adding `toolchain` hint in `go.mod`
* Updated all dependencies to latest

## 1.2.4 (March 04, 2025)

IMPROVEMENTS:

* Enabling CI testing for versions `1.10` and `1.11` of Terraform

NOTES:

* Updated repository Go version to `1.24.0`
  * Adding `toolchain` hint in `go.mod`
* Updated all dependencies to latest

## 1.2.3 (December 21, 2024)

NOTES:

* Updated all dependencies to latest

## 1.2.2 (December 12, 2024)

NOTES:

* Updated repository Go version to `1.23.4` 
* Updated all dependencies to latest
  * Addresses advisory [CVE-2024-45337](https://github.com/advisories/GHSA-v778-237x-gjrc)

## 1.2.1 (November 24, 2024)

NOTES:

* Updated repository Go version to `1.23.3`
* Updated all dependencies to latest
* Updated [golangci-lint](https://golangci-lint.run/) linters
  * Added `asasalint`
  * Added `canonicalheader`
  * Added `containedctx`
  * Added `dupl`
  * Added `dupword`
  * Added `err113`
  * Added `errchkjson`
  * Added `fatcontext`
  * Added `gocognit`
  * Added `gomodguard`
  * Added `noctx`
  * Added `nolintlint`
  * Added `perfsprint`
  * Added `reassign`
  * Added `sqlclosecheck`
  * Added `tagalign`
  * Added `tagliatelle`
  * Added `testifylint`
  * Added `thelper`
  * Added `tparallel`
  * Added `usestdlibvars`
* Because of `err113` linter, created a couple of error types

## 1.2.0 (October 31, 2024)

NEW FEATURES:

* Added support for digest authentication in provider configuration (thanks to [@abarabash-sift](https://github.com/abarabash-sift))
* Added support for ZNode ACL management in `zookeeper_znode` and `zookeeper_sequential_znode` resources (thanks to [@abarabash-sift](https://github.com/abarabash-sift))

IMPROVEMENTS:

* Enabling CI testing for versions `1.9` of Terraform
* Enabling CI testing for version `3.6`, `3.7`, `3.8` and `3.9` of ZooKeeper
* Introduced internal pooling of Clients, reducing amount of ZooKeeper connections necessary to operate
* Minor re-touches to documentation
* Tweaked tests to better handle the closure of Clients after use, reducing amount of concurrent ZooKeeper connections

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
