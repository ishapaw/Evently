# Stage 1: build
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

# Build users service
WORKDIR /app/users
RUN go mod tidy && go build -o /bin/users

# Build bookings service
WORKDIR /app/bookings_view
RUN go mod tidy && go build -o /bin/bookings_view

# Build events service
WORKDIR /app/events
RUN go mod tidy && go build -o /bin/events

# Build gateway service
WORKDIR /app/gateway
RUN go mod tidy && go build -o /bin/gateway

# Build consumers (example: bookings_consumer, cancel_consumer, update_seats_consumer)
WORKDIR /app/bookings/bookings_consumer
RUN go mod tidy && go build -o /bin/bookings_consumer

WORKDIR /app/bookings/cancel_consumer
RUN go mod tidy && go build -o /bin/cancel_consumer

WORKDIR /app/bookings/update_seats_consumer
RUN go mod tidy && go build -o /bin/update_seats_consumer

# Stage 2: runtime
FROM debian:bookworm-slim

WORKDIR /app

# Copy binaries
COPY --from=builder /bin/* /bin/

# Install supervisor
RUN apt-get update && apt-get install -y supervisor && rm -rf /var/lib/apt/lists/*
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

EXPOSE 8080

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
