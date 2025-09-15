# Stage 1: build
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

# Build services
WORKDIR /app/users
RUN go mod tidy && go build -o /app/bin/users

WORKDIR /app/bookings_view
RUN go mod tidy && go build -o /app/bin/bookings_view

WORKDIR /app/events
RUN go mod tidy && go build -o /app/bin/events

WORKDIR /app/gateway
RUN go mod tidy && go build -o /app/bin/gateway

# Build consumers
WORKDIR /app/bookings/bookings_consumer
RUN go mod tidy && go build -o /app/bin/bookings_consumer

WORKDIR /app/bookings/cancel_consumer
RUN go mod tidy && go build -o /app/bin/cancel_consumer

WORKDIR /app/bookings/update_seats_consumer
RUN go mod tidy && go build -o /app/bin/update_seats_consumer

# Stage 2: runtime
FROM ubuntu:24.04

WORKDIR /app

# Copy binaries
COPY --from=builder /app/bin/* /app/bin/

# Install supervisor
RUN apt-get update && apt-get install -y supervisor && rm -rf /var/lib/apt/lists/*
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

EXPOSE 8080

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
