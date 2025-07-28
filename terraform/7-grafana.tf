# resource "helm_release" "grafana" {
#   chart            = "grafana"
#   name             = "grafana"
#   namespace        = "demo"
#   create_namespace = true
#   repository       = "https://grafana.github.io/helm-charts"
#   version          = "9.3.0"
#
#   values = [file("values/grafana.yaml")]
# }