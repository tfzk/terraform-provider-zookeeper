# UTF-8 encoded content
resource "zookeeper_znode" "napoli" {
  path = "/forza/napoli"
  data = "Sempre!"
}

# Base64 encoded content
resource "zookeeper_znode" "napoli_logo" {
  path        = "/forza/napoli/logo"
  data_base64 = filebase64("logo.png")
}
