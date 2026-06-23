# Stage 1: Build biner Go menggunakan image golang resmi
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy modul dan download dependency dulu (memanfaatkan cache layer Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code proyek
COPY . .

# Build biner untuk arsitektur Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-binary cmd/api/main.go

# Stage 2: Gunakan alpine kecil untuk menjalankan aplikasinya saja
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add tzdata

# Ambil biner hasil compile dari Stage 1 tadi
COPY --from=builder /app/app-binary /app/api
COPY .env /app/.env

RUN chmod +x /app/api

EXPOSE 8080

CMD ["/app/api"]