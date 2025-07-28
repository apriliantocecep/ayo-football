# resource "helm_release" "pgsql" {
#   chart            = "postgresql"
#   name             = "postgresql"
#   namespace        = "demo"
#   create_namespace = true
#   repository       = "https://charts.bitnami.com/bitnami"
#   version          = "16.7.21"
#
#   values = [file("values/postgresql.yaml")]
# }