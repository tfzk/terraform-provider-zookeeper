resource "zookeeper_znode" "best_team" {
  path = "/serie/a/best_team"
  data = "It's SSC Napoli, of course!"
}

data "zookeeper_znode" "best_team" {
  path = zookeeper_znode.best_team.path
}

output "best_team_znode_path" {
  value = data.zookeeper_znode.best_team.path
}

output "best_team_znode_data" {
  value = data.zookeeper_znode.best_team.data
}

output "best_team_znode_data_base64" {
  value = data.zookeeper_znode.best_team.data_base64
}
