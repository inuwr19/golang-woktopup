# ----------- Build Stage -----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod & go.sum terlebih dahulu untuk cache efisien
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Pindah ke lokasi main.go
WORKDIR /app/cmd/server

# Build binary
RUN go build -o /woktopup

# ----------- Runtime Stage -----------
FROM alpine:latest

RUN apk --no-cache add ca-certificates mysql-client

WORKDIR /root/

# Copy hasil build
COPY --from=builder /woktopup .
COPY .env .

# Tambahkan dockerize untuk menunggu database
RUN wget https://github.com/jwilder/dockerize/releases/download/v0.7.0/dockerize-alpine-linux-amd64-v0.7.0.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-v0.7.0.tar.gz \
    && rm dockerize-alpine-linux-amd64-v0.7.0.tar.gz

# Jalankan aplikasi setelah database siap
CMD ["dockerize", "-wait", "tcp://db:3306", "-timeout", "60s", "./woktopup"]

EXPOSE 8080
