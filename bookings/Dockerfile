# ---------- Build Stage ----------
FROM golang:1.24-bullseye AS builder

# Set working directory inside container
WORKDIR /app

# Cache dependencies
COPY go.mod ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the bookings service binary
RUN go build -o bookings-service 

# ---------- Runtime Stage ----------
FROM ubuntu:22.04

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/bookings-service .

# Expose API port (adjust if needed)
EXPOSE 8080

# Run the binary
CMD ["./bookings-service"]
