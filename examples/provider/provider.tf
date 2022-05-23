provider "zookeeper" {
  servers         = "zk-server-01:2181,zk-server-02:2181"
  session_timeout = 30
}
