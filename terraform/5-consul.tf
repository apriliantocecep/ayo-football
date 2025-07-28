# resource "helm_release" "consul" {
#   chart            = "consul"
#   name             = "consul"
#   namespace        = "demo"
#   create_namespace = true
#   repository       = "https://helm.releases.hashicorp.com"
#   version          = "1.8.0"
#
#   values = [file("values/consul.yaml")]
# }