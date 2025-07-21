#!/bin/sh

echo "Menunggu Vault siap..."
sleep 5

export VAULT_ADDR=${VAULT_ADDR:-http://vault:8200}
export VAULT_TOKEN=${VAULT_TOKEN:-root}

vault secrets enable -path=posfin -version=2 kv || echo "KV sudah aktif"

echo "Menulis secret ke posfin/config..."
vault kv put posfin/config \
  ARTICLE_SERVICE_PORT="8002" \
  ARTICLE_SERVICE_URL="article-srv" \
  ARTICLE_SERVICE_PROXY="article-service-cluster.local" \
  AUTH_SERVICE_PORT="8001" \
  AUTH_SERVICE_URL="auth-srv" \
  AUTH_SERVICE_PROXY="auth-service-cluster.local" \
  GATEWAY_PORT="8000" \
  MODERATION_SERVICE_PORT="8003" \
  MODERATION_SERVICE_URL="moderation-srv" \
  MODERATION_SERVICE_PROXY="moderation-service-cluster.local" \
  RABBITMQ_URL="amqp://guest:guest@rabbitmq-server:5672/" \
  CONSUL_URL="consul-server:8500" \
  TRAEFIK_GRPC_PROXY_URL="traefik-server:9000" \
  TRAEFIK_WEB_PROXY_URL="traefik-server:8081"

echo "Menulis secret ke posfin/auth-service..."
vault kv put posfin/auth-service \
  DATABASE_URL="host=postgres-server user=root password=admin dbname=posfin port=5432 sslmode=disable TimeZone=Asia/Jakarta" \
  JWT_ACCESS_TOKEN_EXPIRATION_MINUTES="1440" \
  JWT_ISSUER="posfin" \
  JWT_SECRET_KEY="9232c8cd6cfc4c4ed3cb848682bc883dfb8964f3f04cc0811f56ff0c49ad20f68aec62c5eb40ce0235f0dc7f51bd8a3f"

echo "Menulis secret ke posfin/article-service..."
vault kv put posfin/article-service \
  DATABASE_URL="mongodb://root:admin@mongodb-server:27017/?authMechanism=SCRAM-SHA-1"

echo "Inisialisasi secret selesai."
