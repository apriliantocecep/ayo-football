# resource "helm_release" "vault" {
#   chart            = "vault"
#   name             = "vault"
#   namespace        = "demo"
#   create_namespace = true
#   repository       = "https://helm.releases.hashicorp.com"
#   version          = "0.30.0"
#
#   values = [file("values/vault.yaml")]
# }
#
# resource "terraform_data" "vault_secrets" {
#   triggers_replace = [helm_release.vault.id]
#
#   provisioner "local-exec" {
#     command = "scripts/init-vault.sh"
#   }
#
#   depends_on = [helm_release.vault]
# }