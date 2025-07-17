package shared

import (
	"context"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"log"
	"os"
)

type VaultClient struct {
	Client *vault.Client
}

func NewVaultClient() *VaultClient {
	// get vault env
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultAddr == "" || vaultToken == "" {
		log.Fatalf("VAULT_ADDR and VAULT_TOKEN is not set")
	}

	config := vault.DefaultConfig()
	config.Address = vaultAddr

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create Vault client: %v", err)
	}
	client.SetToken(vaultToken)

	return &VaultClient{Client: client}
}

func (v *VaultClient) GetSecret(path string) (map[string]interface{}, error) {
	secret, err := v.Client.KVv2("posfin").Get(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at %s", path)
	}

	return secret.Data, nil
}
