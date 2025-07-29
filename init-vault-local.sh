#!/bin/sh

echo "Menunggu Vault siap..."
sleep 5

export VAULT_ADDR=${VAULT_ADDR:-http://localhost:8200}
export VAULT_TOKEN=${VAULT_TOKEN:-root}

vault secrets enable -path=ayofootball -version=2 kv || echo "KV sudah aktif"

echo "Menulis secret ke ayofootball/config..."
vault kv put ayofootball/config \
  DATABASE_URL="host=localhost user=root password=admin dbname=ayofootball port=5432 sslmode=disable TimeZone=Asia/Jakarta" \
  TEAM_SERVICE_PORT="8002" \
  TEAM_SERVICE_URL="localhost" \
  AUTH_SERVICE_PORT="8001" \
  AUTH_SERVICE_URL="localhost" \
  GATEWAY_PORT="8000" \
  PLAYER_SERVICE_PORT="8003" \
  PLAYER_SERVICE_URL="localhost" \
  RABBITMQ_URL="amqp://guest:guest@rabbitmq-server:5672/" \
  CONSUL_URL="consul-server:8500" \
  TRAEFIK_GRPC_PROXY_URL="traefik-server:9000" \
  TRAEFIK_WEB_PROXY_URL="traefik-server:8081" \
  OTEL_COLLECTOR_ENDPOINT_HTTP="otel-collector:4318" \
  OTEL_COLLECTOR_ENDPOINT_GRPC="otel-collector:4317"

echo "Menulis secret ke ayofootball/auth-service..."
vault kv put ayofootball/auth-service \
  JWT_ACCESS_TOKEN_EXPIRATION_MINUTES="1440" \
  JWT_ISSUER="ayofootball" \
  JWT_SECRET_KEY="9232c8cd6cfc4c4ed3cb848682bc883dfb8964f3f04cc0811f56ff0c49ad20f68aec62c5eb40ce0235f0dc7f51bd8a3f"

echo "Menulis secret ke ayofootball/article-service..."
vault kv put ayofootball/article-service \
  DATABASE_URL="mongodb://root:admin@mongodb-server:27017/?authMechanism=SCRAM-SHA-1"

echo "Inisialisasi secret selesai."
