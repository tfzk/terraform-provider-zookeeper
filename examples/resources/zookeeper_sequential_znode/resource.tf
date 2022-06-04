# Creating 3 sequential nodes in the same "directory".
# Nodes will have final `.path` of:
#
#  /sequential/nodes/go/in/here/0000000001
#  /sequential/nodes/go/in/here/0000000002
#  /sequential/nodes/go/in/here/0000000003
#
# NOTE: the order they appear in the config (i.e. `seqA`, `seqB`, `seqC`)
# does NOT guarantee they will be respectively named
# `0000000001`, `0000000002`, `0000000003`.
#
# Their creation is parallelized by Terraform, so ZooKeeper will give each
# a unique name but not in the order we see here.
#
# Consider using [`depends_on`](https://www.terraform.io/language/meta-arguments/depends_on)
# if you want to gain control of this, OR consider using `zookeeper_znode`s instead.

resource "zookeeper_znode" "dir" {
  path = "/sequential/nodes/go/in/here"
}

resource "zookeeper_sequential_znode" "seqA" {
  path_prefix = zookeeper_znode.dir.path + "/"
  data        = "some data"
}

resource "zookeeper_sequential_znode" "seqB" {
  path_prefix = zookeeper_znode.dir.path + "/"
  data        = "some more data"
}

resource "zookeeper_sequential_znode" "seqC" {
  path_prefix = zookeeper_znode.dir.path + "/"
  data        = "even more data"
}
