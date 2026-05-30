# soom-be-go

Backend service untuk aplikasi **Soom**, dibangun menggunakan Go, Gin, GORM, dan PostgreSQL.

---

## Prasyarat

Pastikan sudah terinstall:

- [Go](https://go.dev/dl/) >= 1.21
- [PostgreSQL](https://www.postgresql.org/download/) >= 14
- [Air](https://github.com/air-verse/air) *(opsional, untuk hot reload)*

Install Air:
```bash
go install github.com/air-verse/air@latest
```

---

## Cara Menjalankan

### 1. Clone Repository

```bash
git clone <repo-url>
cd soom-be-go
```

### 2. Setup Environment

```bash
cp .env.example .env
```

Sesuaikan konfigurasi di file `.env` dengan environment lokal kamu.

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Jalankan Aplikasi

**Development (hot reload):**
```bash
air
```

**Normal:**
```bash
go run ./cmd/api
```

**Generate file migration**
```bash
goose -dir db/migrations create <script_name> sql
```

Aplikasi berjalan di `http://localhost:8080`

> Migrasi database berjalan **otomatis** saat aplikasi pertama kali distart.

---

## Troubleshooting

**Gagal connect database** — pastikan PostgreSQL sudah berjalan dan konfigurasi di `.env` sudah benar.

**Port sudah dipakai** — ganti `APP_PORT` di `.env`.

**Air tidak ditemukan** — pastikan `$GOPATH/bin` sudah ada di `PATH`.
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```