# ---------- Stage 1: Build ----------
FROM golang:1.24.2 AS builder

WORKDIR /app

# Copy go mod and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY shared ./shared
COPY services/article ./services/article

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o article-service ./services/article/cmd/main.go

# ---------- Stage 2: Run ----------
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/article-service /article-service

EXPOSE 8002

# Run binary
ENTRYPOINT ["/article-service"]
