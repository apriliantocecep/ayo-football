# If you want using helm chart from bitnami/mongodb
# resource "helm_release" "mongodb" {
#   chart            = "mongodb"
#   name             = "mongodb"
#   namespace        = "demo"
#   create_namespace = true
#   repository       = "https://charts.bitnami.com/bitnami"
#   version          = "13.18.5"
#
#   values = [file("values/mongodb.yaml")]
# }

# If you want using kubectl manifest
# resource "kubectl_manifest" "demo_namespace" {
#   yaml_body = file("${path.module}/templates/mongodb/namespace.yaml")
# }
#
# resource "kubectl_manifest" "mongodb_deployment" {
#   yaml_body = file("${path.module}/templates/mongodb/deployment.yaml")
# }
#
# resource "kubectl_manifest" "mongodb_service" {
#   depends_on = [kubectl_manifest.mongodb_deployment]
#
#   yaml_body = file("${path.module}/templates/mongodb/service.yaml")
# }