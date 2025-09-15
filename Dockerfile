# Step 1: Build stage
FROM golang:1.24 AS builder
WORKDIR /src

# copy everything
COPY . .

# build all services
RUN go build -o /bin/users ./users
RUN go build -o /bin/bookings ./bookings
RUN go build -o /bin/events ./events
RUN go build -o /bin/bookings_view ./bookings_view
RUN go build -o /bin/gateway ./gateway

# Step 2: Runtime
FROM debian:bookworm-slim
WORKDIR /app

# install supervisor to manage multiple processes
RUN apt-get update && apt-get install -y supervisor && rm -rf /var/lib/apt/lists/*

# copy binaries
COPY --from=builder /bin/* /app/

# copy supervisor config
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

EXPOSE 8080
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
