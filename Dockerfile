# syntax=docker/dockerfile:1
FROM golang:1.18.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server

EXPOSE 8080
CMD ["/app/server"]