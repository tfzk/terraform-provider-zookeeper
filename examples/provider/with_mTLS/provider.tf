provider "zookeeper" {
  servers         = "zk-server-01:2182,zk-server-02:2182"
  session_timeout = 30
  tls_enabled     = true
  tls_ca_file     = "/path/to/ca_cert.pem"
  tls_cert_file   = "/path/to/client_cert.pem"
  tls_key_file    = "/path/to/client_key.key"
}
