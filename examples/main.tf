terraform {
  required_providers {
    zookeeper = {
      source  = "github.com/tfzk/zookeeper-legacy"
      # NOTE: Make sure the pick the latest version when using the provider.
      # Here we refer to the one installed via `Makefile` as part of the
      # local development lifecycle.
      version = "0.0.0-dev"
    }
  }
}

provider "zookeeper" {
  # NOTE: This servers string connects to a local ensemble that
  # can be spawn using the `docker-compose.yml` inside
  # `<prj root>/tools/zk-local-ensemble`.
  servers = "localhost:2181,localhost:2182,localhost:2183"
}

resource "zookeeper_znode" "znode_with_json" {
  path = "/examples/zookeeper-legacy/json"

  data = jsonencode({
    ivan        = "uno1"
    fabio       = "due2"
    lorena      = "tre3"
    adriano     = "quattro4"
    marcello    = "cinque5"
  })
}

resource "zookeeper_znode" "empty_parent_znode" {
  path = "/examples/zookeeper-legacy/parent"
}

resource "zookeeper_sequential_znode" "znode_seq_type1" {
  path_prefix = format("%s/znode_seq_type1/", zookeeper_znode.empty_parent_znode.path)
  data        = "this is not json"
}

resource "zookeeper_sequential_znode" "znode_seq_type2" {
  path_prefix = format("%s/znode_seq_type2-", zookeeper_znode.empty_parent_znode.path)
  data        = "this is still not json"
}

data "zookeeper_znode" "data_znode_seq_type2" {
  path = zookeeper_sequential_znode.znode_seq_type2.path
}

output "znodes" {
  value = {
    resources = {
      znode_with_json     = zookeeper_znode.znode_with_json
      znode_seq_type1 = zookeeper_sequential_znode.znode_seq_type1
      znode_seq_type2 = zookeeper_sequential_znode.znode_seq_type2
    }
    data = {
      data_znode_seq_type2 = data.zookeeper_znode.data_znode_seq_type2
    }
  }
}