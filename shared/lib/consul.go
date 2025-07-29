package lib

import (
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	capi "github.com/hashicorp/consul/api"
	"log"
	"time"
)

type ConsulClient struct {
	Client *capi.Client
}

func (c *ConsulClient) StartTTLHeartbeat(checkID string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		err := c.Client.Agent().UpdateTTL(checkID, "online", capi.HealthPassing)
		if err != nil {
			log.Printf("failed to update TTL for worker %s: %v", checkID, err)
		}
	}
}

func (c *ConsulClient) ServiceRegister(serviceID, serviceName, address string, port int, checkConfig *capi.AgentServiceCheck, tags []string) {
	registration := &capi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Tags:    tags,
		Port:    port,
		Address: address,
		Check:   checkConfig,
	}

	err := c.Client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("failed to register gRPC service to consul agent: %v", err)
	}
}

func NewConsulClient(vaultClient *shared.VaultClient) *ConsulClient {
	secret := utils.GetVaultSecretConfig(vaultClient)
	url := secret["CONSUL_URL"]
	if url == nil || url == "" {
		log.Fatalln("CONSUL_URL is not set")
	}

	config := capi.DefaultConfig()
	config.Address = url.(string)

	client, err := capi.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create consul client: %v", err)
	}

	return &ConsulClient{
		Client: client,
	}
}
