# ---------- Stage 1: Build ----------
FROM golang:1.24.2 AS builder

# Install build dependencies untuk CGO (fiber prefork butuh ini)
RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy go mod and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY shared ./shared
COPY services/auth ./services/auth
COPY services/article ./services/article
COPY gateway/rest ./gateway/rest

# Build dengan CGO untuk fiber prefork
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gateway-service ./gateway/rest/cmd/main.go

# Make sure binary is executable
RUN chmod +x /app/gateway-service

# ---------- Stage 2: Run ----------
# Gunakan image yang punya shell tapi tetap kecil
FROM debian:bookworm-slim

# Install minimal runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/gateway-service /gateway-service

EXPOSE 8000

# Shell form untuk fiber prefork
CMD ./gateway-service