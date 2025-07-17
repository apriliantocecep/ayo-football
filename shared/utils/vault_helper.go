package utils

import (
	"github.com/apriliantocecep/posfin-blog/shared"
	"log"
	"strconv"
)

func GetVaultSecret(client *shared.VaultClient, path string) map[string]interface{} {
	secret, err := client.GetSecret(path)
	if err != nil {
		log.Fatalf("vault error: %s", err.Error())
	}
	return secret
}

func GetVaultSecretConfig(client *shared.VaultClient) map[string]interface{} {
	secret := GetVaultSecret(client, "config")
	return secret
}

func GetVaultSecretAuthSvc(client *shared.VaultClient) map[string]interface{} {
	secret := GetVaultSecret(client, "auth-service")
	return secret
}

func GetVaultSecretArticleSvc(client *shared.VaultClient) map[string]interface{} {
	secret := GetVaultSecret(client, "article-service")
	return secret
}

func ParsePort(portStr string) int {
	value, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("invalid port number: %s", err.Error())
	}
	return value
}
