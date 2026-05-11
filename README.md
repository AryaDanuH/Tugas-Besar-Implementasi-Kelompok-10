# LibMatch - Sistem Manajemen Perpustakaan

LibMatch adalah aplikasi web full-stack untuk peminjaman buku, memberikan review, dan admin untuk mengelola katalog buku.

## Fitur Utama

- **Autentikasi Pengguna**: Registrasi dan login dengan keamanan password hashing
- **Katalog Buku**: Browse, search, dan filter buku berdasarkan kategori
- **Sistem Peminjaman**: Keranjang peminjaman dan checkout
- **Review Buku**: Pengguna dapat memberikan rating dan review untuk buku yang dipinjam
- **Admin Panel**: Kelola buku, verifikasi peminjaman, dan lihat riwayat
- **Profile**: Kelola data pengguna dan upload foto profil
- **Populer Books**: Buku yang paling sering dipinjam
- **Lokasi Pickup**: Kelola lokasi pengambilan buku

## Tech Stack

**Backend:**
- Go (Golang)
- MySQL
- REST API

**Frontend:**
- HTML5
- CSS3
- JavaScript (Vanilla)

## Instalasi

### Prerequisites

- Go 1.16+
- MySQL 5.7+
- Git
- Browser modern

### Step 1: Clone Repository

```bash
git clone <repository-url>
cd LibMatch
```

### Step 2: Setup Database

1. Buka MySQL command line:
```bash
mysql -u root -p
```

2. Buat database:
```sql
CREATE DATABASE libmatch;
USE libmatch;
```

3. Import schema (jika ada file .sql):
```sql
SOURCE BackEnd/database/schema.sql;
```

4. Atau buat tabel secara manual sesuai struktur di `BackEnd/Model/`

### Step 3: Setup Backend

1. Navigate ke folder backend:
```bash
cd BackEnd
```

2. Update konfigurasi database di `main.go`:
```go
// Ubah sesuai credentials MySQL Anda
const (
    dsn = "root:password@tcp(localhost:3306)/libmatch"
)
```

3. Jalankan server:
```bash
go run main.go
```

Output akan menampilkan:
```
Server running on :8080
```

### Step 4: Buka Frontend

1. Buka file dengan browser:
```
FrontEnd/index.html
```

Atau jika perlu local server:
```bash
cd FrontEnd
python3 -m http.server 3000
# Buka http://localhost:3000
```

## Struktur Folder

```
LibMatch/
├── BackEnd/
│   ├── main.go              # Entry point server
│   ├── Handle/              # Handler HTTP endpoints
│   ├── Model/               # Database models
│   ├── middleware/          # Middleware (auth, logging)
│   └── config.go            # Konfigurasi
├── FrontEnd/
│   ├── index.html           # Halaman utama
│   ├── dashboard-logged-in.html
│   ├── admin-all-books.html
│   ├── js/                  # JavaScript files
│   ├── css/                 # Stylesheet
│   └── images/              # Assets
├── sequence-diagrams/       # Dokumentasi UML
├── DPPL_Libmatch.pdf       # Dokumentasi lengkap
└── README.md               # File ini
```

## Penggunaan

### Pengguna Regular

1. **Register**: Klik "Daftar" di halaman utama
2. **Login**: Masukkan email dan password
3. **Browse Buku**: Lihat katalog, cari, atau filter berdasarkan kategori
4. **Tambah ke Keranjang**: Pilih buku dan tambahkan ke cart
5. **Checkout**: Isi detail peminjaman dan pilih lokasi pickup
6. **Review**: Setelah meminjam, berikan rating dan review

### Admin

1. **Login sebagai Admin**: Gunakan akun admin
2. **Kelola Buku**: Tambah, edit, atau hapus buku
3. **Verifikasi Peminjaman**: Approve atau reject peminjaman pengguna
4. **Lihat Laporan**: Analytics populer buku dan pengguna aktif

## API Endpoints

### User Management
- `POST /api/auth/register` - Register pengguna baru
- `POST /api/auth/login` - Login pengguna
- `GET /api/users/{id}` - Ambil data pengguna
- `PUT /api/users/{id}` - Update profil pengguna

### Books
- `GET /api/books` - Ambil semua buku
- `GET /api/books/{id}` - Ambil detail buku
- `GET /api/books/search?q=query` - Search buku
- `GET /api/books/popular` - Buku populer
- `POST /api/books` - Tambah buku (admin)
- `PUT /api/books/{id}` - Edit buku (admin)
- `DELETE /api/books/{id}` - Hapus buku (admin)

### Borrowing
- `POST /api/borrows` - Buat peminjaman
- `GET /api/users/{id}/borrows` - Riwayat peminjaman user
- `PUT /api/borrows/{id}/status` - Update status peminjaman

### Reviews
- `POST /api/reviews` - Tambah review
- `GET /api/books/{id}/reviews` - Ambil review buku

## Troubleshooting

### Error: Port 8080 already in use
```bash
# Cari proses yang menggunakan port 8080
lsof -i :8080

# Kill proses
kill -9 <PID>

# Atau ubah port di main.go
```

### Error: MySQL connection refused
- Pastikan MySQL sudah running: `mysql -u root -p`
- Check credentials di `main.go`
- Pastikan database `libmatch` sudah dibuat

### Frontend tidak connect ke backend
- Pastikan backend running di port 8080
- Check API_URL di `FrontEnd/js/config.js` atau di file JavaScript
- Buka Developer Console (F12) untuk lihat error

### Database not found
```bash
# Login MySQL
mysql -u root -p

# Buat database
CREATE DATABASE libmatch;

# Import schema jika ada
USE libmatch;
SOURCE BackEnd/database/schema.sql;
```
## License

Project ini dibuat untuk keperluan akademik.

---

**Dibuat dengan ❤️ oleh Tim LibMatch**
