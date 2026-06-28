# BE Pasar Malam

Backend API untuk aplikasi Pasar Malam — marketplace jajanan pasar malam. User bisa browsing produk, masukin ke keranjang, checkout, dan bayar. Dibangun pakai Go, Gin framework, MySQL, dan Firebase Auth.

## Fitur

### Autentikasi
- Login lewat Firebase Auth (Google Sign-In atau email/password)
- Verifikasi token Firebase ke backend
- JWT buat autentikasi endpoint yang butuh login
- Verifikasi email (dev mode ada endpoint khusus buat skip verifikasi)

### Produk
- List semua produk (bisa filter kategori)
- Detail produk berdasarkan ID
- Data produk di-seed otomatis saat pertama kali jalan (12 item jajanan pasar malam)

### Keranjang Belanja (Cart)
- Lihat isi keranjang
- Tambah produk ke keranjang
- Update jumlah item
- Hapus item dari keranjang
- Kosongkan seluruh keranjang

### Pesanan (Order)
- Checkout dari keranjang jadi pesanan
- Lihat daftar pesanan milik user
- Detail pesanan berdasarkan ID
- Support beberapa metode pembayaran (VA number, GoPay deeplink)

## Struktur Project

```
be_pasar_malam/
├── config/           # Konfigurasi dari .env
│   └── config.go
├── database/         # Koneksi MySQL & Firebase
│   ├── mysql.go
│   └── firebase.go   (via config)
├── handlers/         # Handler HTTP (controller)
│   ├── auth.go       # Login, verifikasi token
│   ├── products.go   # List & detail produk
│   ├── cart.go       # CRUD keranjang
│   ├── orders.go     # Checkout & riwayat pesanan
│   └── health.go     # Health check
├── middleware/        # Middleware JWT & logger
│   ├── auth.go
│   └── logger.go
├── models/           # Model database (GORM)
│   ├── user.go
│   ├── product.go
│   ├── cart.go
│   └── order.go
├── routes/           # Definisi route API
│   ├── routes.go
│   └── models_alias.go
├── seed/             # Data awal produk
│   └── seed.go
├── services/         # Business logic
│   ├── jwt.go
│   └── firebase_auth.go
├── main.go           # Entry point
├── go.mod
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Bahasa | Go 1.22 |
| HTTP Framework | Gin |
| ORM | GORM |
| Database | MySQL |
| Autentikasi | Firebase Auth + JWT |
| Container | Docker + Docker Compose |

## API Endpoints

### Public
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | `/v1/health` | Cek server hidup |
| POST | `/v1/auth/verify-token` | Verifikasi token Firebase |
| POST | `/v1/auth/dev-verify-email` | Verifikasi email (dev only) |
| GET | `/v1/products` | List semua produk |
| GET | `/v1/products/:id` | Detail produk |

### Butuh Login (JWT)
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | `/v1/auth/me` | Info user yang sedang login |
| PUT | `/v1/auth/fcm-token` | Update FCM token |
| GET | `/v1/cart` | Lihat isi keranjang |
| POST | `/v1/cart` | Tambah item ke keranjang |
| PUT | `/v1/cart/:id` | Update jumlah item |
| DELETE | `/v1/cart/:id` | Hapus item dari keranjang |
| DELETE | `/v1/cart` | Kosongkan keranjang |
| GET | `/v1/orders` | Daftar pesanan saya |
| POST | `/v1/orders/checkout` | Checkout keranjang |
| GET | `/v1/orders/:id` | Detail pesanan |

## Data Seed (Produk Awal)

Saat pertama kali server jalan, otomatis di-seed 12 produk:

**Makanan:** Sate Ayam Madura, Bakso Urat, Bakso Bakar, Soto Ayam Kampung, Nasi Goreng Spesial, Mie Ayam Bakso

**Minuman:** Es Teh Manis, Es Jeruk Peras, Bajigur, Bandrek

**Snack:** Pisang Goreng, Tahu Bulat

## Cara Menjalankan

### 1. Persiapan

Buat file `.env` dari contoh:
```bash
cp .env.example .env
```

Isi konfigurasi di `.env` sesuai kebutuhan.

Pastikan `firebase_service_account.json` ada. Backend ini share Firebase project yang sama dengan BE Dompet Digital.

### 2. Jalankan dengan Docker

```bash
docker compose up --build
```

API jalan di port `8082`. Backend ini pakai network yang sama dengan `be_dompet_digital` supaya bisa share MySQL.

### 3. Jalankan manual

```bash
go run main.go
```

### 4. Cek server

```bash
curl http://localhost:8082/v1/health
```

### 5. Matikan Docker

```bash
docker compose down
```

## Koneksi dari Flutter

| Device | URL |
|--------|-----|
| Android emulator | `http://10.0.2.2:8082` |
| HP fisik | `http://<IP_LAPTOP>:8082` |
| iOS simulator | `http://localhost:8082` |

## Proyek Terkait

| Folder | Hubungan |
|--------|----------|
| [📂 `pasar_malam`](../pasar_malam) | Flutter app (frontend) yang pakai backend ini |
| [📂 `be_dompet_digital`](../be_dompet_digital) | Backend dompet digital — user bisa bayar pakai saldo dari sana |
| [📂 `dompet_digital`](../dompet_digital) | Flutter app e-money — sumber saldo buat transaksi di marketplace ini |

## Environment Variables

| Variable | Default | Keterangan |
|----------|---------|------------|
| `PORT` | `8082` | Port server |
| `DB_HOST` | `localhost` | Host MySQL |
| `DB_PORT` | `3306` | Port MySQL |
| `DB_USER` | `useremoney` | User MySQL |
| `DB_PASSWORD` | `Password#123` | Password MySQL |
| `DB_NAME` | `pasarmalam` | Nama database |
| `JWT_SECRET` | - | Secret key buat JWT |
| `JWT_EXPIRY_HOURS` | `168` | Masa berlaku JWT (7 hari) |
| `FIREBASE_CREDENTIALS_PATH` | `firebase_service_account.json` | Path service account |
| `FIREBASE_API_KEY` | - | Firebase Web API Key |
# be_pasar_malam
