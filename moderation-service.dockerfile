# ---------- Stage 1: Build ----------
FROM golang:1.24.2 AS builder

WORKDIR /app

# Copy go mod and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY shared ./shared
COPY services/moderation ./services/moderation

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o moderation-service ./services/moderation/cmd/main.go

# ---------- Stage 2: Run ----------
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/moderation-service /moderation-service

EXPOSE 8003

# Run binary
ENTRYPOINT ["/moderation-service"]
