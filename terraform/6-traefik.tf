# resource "helm_release" "traefik" {
#   chart            = "traefik"
#   name             = "traefik"
#   namespace        = "demo"
#   create_namespace = true
#   repository       = "https://traefik.github.io/charts"
#   version          = "36.3.0"
#
#   values = [file("values/traefik.yaml")]
# }