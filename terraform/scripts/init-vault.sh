#!/bin/sh

echo "Menunggu Vault siap..."
sleep 5

export VAULT_ADDR="http://vault.minikube"
export VAULT_TOKEN="root"

echo "Retry jika belum ready"
until curl -s $VAULT_ADDR/v1/sys/health; do echo "Waiting for Vault... [Only MacOS] Make sure you run 'minikube tunnel'"; sleep 5; done

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
  RABBITMQ_URL="amqp://guest:guest@rabbitmq.minikube:5672/" \
  CONSUL_URL="consul.minikube" \
  TRAEFIK_GRPC_PROXY_URL="traefik-server:9000" \
  TRAEFIK_WEB_PROXY_URL="traefik-server:8081" \
  OTEL_COLLECTOR_ENDPOINT_HTTP="otel-collector:4318" \
  OTEL_COLLECTOR_ENDPOINT_GRPC="otel-collector:4317"

echo "Menulis secret ke posfin/auth-service..."
vault kv put posfin/auth-service \
  DATABASE_URL="host=postgresql.demo.svc.cluster.local user=root password=admin dbname=posfin port=5432 sslmode=disable TimeZone=Asia/Jakarta" \
  JWT_ACCESS_TOKEN_EXPIRATION_MINUTES="1440" \
  JWT_ISSUER="posfin" \
  JWT_SECRET_KEY="9232c8cd6cfc4c4ed3cb848682bc883dfb8964f3f04cc0811f56ff0c49ad20f68aec62c5eb40ce0235f0dc7f51bd8a3f"

echo "Menulis secret ke posfin/article-service..."
vault kv put posfin/article-service \
  DATABASE_URL="mongodb://root:admin@mongodb.demo.svc.cluster.local:27017/?authMechanism=SCRAM-SHA-1"

echo "Inisialisasi secret selesai."
