provider "zookeeper" {
  servers               = "zk-server-01:2182,zk-server-02:2182"
  session_timeout       = 30
  tls_enable            = true
  tls_root_ca_cert_path = "/path/to/root_CA_certificate.pem"
  tls_cert_path         = "/path/to/client_certificate.pem"
  tls_key_path          = "/path/to/client_key.key"
}
