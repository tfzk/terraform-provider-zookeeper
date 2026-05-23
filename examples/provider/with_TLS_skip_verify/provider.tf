provider "zookeeper" {
  servers         = "zk-server-01:2182,zk-server-02:2182"
  session_timeout = 30
  tls_enabled     = true
  tls_skip_verify = true
}
