FROM alpine:latest

WORKDIR /app

# Install tzdata agar tidak ada kendala format waktu atau timezone di aplikasi Go
RUN apk --no-cache add tzdata

# Salin file binary matang (app-binary) yang sudah kamu build di lokal tadi
COPY app-binary /app/api

RUN chmod +x /app/api

# Buka port sesuai dengan yang digunakan aplikasi Go kamu
EXPOSE 8080

# Jalankan file binary-nya saat container dinyalakan
CMD ["/app/api"]