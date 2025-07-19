ui = true
disable_mlock = true

listener "tcp" {
  tls_disable = true
  address = "[::]:8200"
}

storage "file" {
  path = "/vault/data"
}