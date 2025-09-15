# 1️⃣ Build stage
FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Copy go.mod and go.sum from both modules
COPY go.mod go.sum ./
COPY users/go.mod users/
COPY users/go.sum users/
RUN go mod download

# Copy all the source code for both modules
COPY . .

# Build the gateway binary from its subdirectory
WORKDIR /app/gateway
RUN go build -o /app/gateway .

# 2️⃣ Runtime stage
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/gateway .

EXPOSE 8080
CMD ["./gateway"]