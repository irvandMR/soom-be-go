FROM alpine:latest
WORKDIR /app
# Kita HANYA meng-copy biner yang sudah jadi dari laptopmu
COPY app-binary .
# Jika ada file .env atau config, copy juga
COPY .env .
# Jalankan binernya
CMD ["./app-binary"]