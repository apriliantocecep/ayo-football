# resource "helm_release" "rabbitmq" {
#   chart            = "rabbitmq"
#   name             = "rabbitmq"
#   namespace        = "demo"
#   create_namespace = true
#   repository       = "https://charts.bitnami.com/bitnami"
#   version          = "16.0.11"
#
#   values = [
#     templatefile("${path.module}/templates/rabbitmq-values.yaml.tmpl", {
#       load_definition = indent(6, file("${path.module}/scripts/rabbitmq-definitions.json"))
#     })
#   ]
# }