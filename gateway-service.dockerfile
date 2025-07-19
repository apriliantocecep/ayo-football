# ---------- Stage 1: Build ----------
FROM golang:1.24.2 AS builder

WORKDIR /app

# Copy go mod and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY shared ./shared
COPY services ./services
COPY gateway/rest ./gateway/rest

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway-service ./gateway/rest/cmd/main.go

# ---------- Stage 2: Run ----------
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/gateway-service /gateway-service

EXPOSE 8000

# Run binary
ENTRYPOINT ["/gateway-service"]
