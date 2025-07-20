node_name  = "consul-server"
server     = true
datacenter = "dc1"
data_dir = "consul/data"
ui_config {
  enabled = true
}
bind_addr = "0.0.0.0"
client_addr = "0.0.0.0"
