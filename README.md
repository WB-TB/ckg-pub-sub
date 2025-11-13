# PubSub Service untuk Integrasi TB di ASIK CKG

Sistem Pub/Sub untuk integrasi data skrining CKG (Cegah Komplikasi Gizi) dan TB (Tuberkulosis) antara CKG dan SITB.

## Deskripsi Project

Project ini adalah sistem messaging yang menggunakan Google Cloud Pub/Sub untuk menghubungkan dua sistem:
- **ASIK CKG (Cek Kesehatan Gratis)** - Sistem nasional untuk skrining kesehatan gratis
- **SITB (Sistem Informasi Tuberkulosis)** - Sistem informasi TB

Sistem ini memungkinkan pertukaran data pasien secara semi-realtime maupun terjadwal antara kedua sistem dengan memastikan data yang dikirim adalah data yang valid dan belum diproses sebelumnya.

## Fitur Utama

- **Producer**: Mengirim data skrining TB yang terduga dari CKG ke SITB
- **Consumer**: Menerima dan memproses data status pasien TB dari SITB
- **Change Stream**: Mendeteksi perubahan data secara real-time di MongoDB
- **Duplicate Prevention**: Mencegah pengiriman atau pemrosesan data yang sama berulang kali
- **Batch Processing**: Mengirim data dalam batch untuk performa yang lebih baik
- **Multi-Database Support**: Mendukung MongoDB dan SQL database

## Struktur Project

```
pubsub-ckg-tb/
├── cmd/
│   ├── consumer/          # Consumer application
│   └── producer/          # Producer application
├── internal/
│   ├── app/               # Application layer
│   │   └── ckg/          # CKG specific logic
│   ├── config/            # Configuration management
│   ├── db/                # Database layer
│   │   ├── connection/   # Database connections
│   │   ├── mongo/        # MongoDB implementation
│   │   ├── sql/          # SQL implementation
│   │   └── utils/        # Database utilities
│   ├── models/           # Data models
│   └── pubsub/           # Pub/Sub implementation
├── schema-sitb-ckg.sql   # Database schema
└── go.mod               # Go module file
```

## Komponen Utama

### 1. Producer (`cmd/producer/main.go`)
- Mengambil data skrining TB yang terduga dari database CKG
- Mengirim data ke Google Cloud Pub/Sub
- Mendukung mode watch untuk monitoring real-time

### 2. Consumer (`cmd/consumer/main.go`)
- Menerima data dari Google Cloud Pub/Sub
- Memproses data status pasien TB
- Menyimpan hasil ke database CKG

### 3. CKG Transmitter (`internal/app/ckg/trasmitter.go`)
- Menyiapkan data untuk dikirim
- Mendeteksi perubahan data melalui change stream
- Mengelola batch pengiriman

### 4. CKG Receiver (`internal/app/ckg/receiver.go`)
- Memvalidasi pesan masuk
- Memproses data status pasien
- Mencegah duplikasi data

## Konfigurasi

Project menggunakan file konfigurasi yang dapat diatur melalui:
- File `.env`
- Environment variables
- Default values

### Konfigurasi Utama

```yaml
app:
  env: development
  loglevel: info

google:
  project: your-gcp-project-id
  credentials: path/to/service-account.json
  debug: false

pubsub:
  topic: ckg-tb-topic
  subscription: ckg-tb-subscription
  messageordering: true

db:
  driver: mongodb  # atau mysql, postgresql
  host: localhost
  port: 27017
  database: ckg_db
  username: ""
  password: ""

ckg:
  tablemasterwilayah: master_wilayah
  tablemasterfaskes: master_faskes
  tableskrining: skrining_ckg
  tablestatus: status_pasien
  tableincoming: ckg_pubsub_incoming
  tableoutgoing: ckg_pubsu_outgoing
  markerfield: marker
  markerconsume: consumed
  markerproduce: produced
```

## Instalasi

### Prasyarat

- Go 1.25 atau lebih tinggi
- Google Cloud SDK
- Service account dengan akses ke Google Cloud Pub/Sub
- Docker dan Docker Compose (untuk penggunaan container)

### Instalasi Lokal

1. Clone repository:
```bash
git clone <repository-url>
cd pubsub-ckg-tb
```

2. Install dependencies:
```bash
go mod download
```

3. Konfigurasi environment:
```bash
cp .env.example .env
# Edit .env dengan konfigurasi Anda
```

4. Buat topik dan subscription di Google Cloud Pub/Sub:
```bash
gcloud pubsub topics create ckg-tb-topic
gcloud pubsub subscriptions create ckg-tb-subscription --topic=ckg-tb-topic
```

### Instalasi Menggunakan Docker

1. Clone repository:
```bash
git clone <repository-url>
cd pubsub-ckg-tb
```

2. Buat file credentials.json untuk akses Google Cloud:
```bash
# Buat service account di Google Cloud Console
# Download JSON key dan simpan sebagai credentials.json
```

3. Jalankan dengan Docker Compose:
```bash
# Jalankan semua service
docker-compose up -d

# Untuk development dengan live reload
docker-compose up -d producer consumer

# Lihat logs
docker-compose logs -f producer
docker-compose logs -f consumer
```

4. Build untuk production:
```bash
# Build binary production
docker-compose run builder

# Binary akan tersedia di folder ./bin
```

### Live Reload Development Mode

Project mendukung live reload menggunakan library [Air](https://github.com/air-verse/air) untuk development yang lebih efisien.

#### Menggunakan Air Lokal

1. Install Air:
```bash
go install github.com/air-verse/air@latest
```

2. Jalankan dengan live reload:
```bash
# Untuk producer
air --build.cmd "go build -o /app/producer /app/cmd/producer/main.go" --build.bin "/app/producer"

# Untuk consumer
air --build.cmd "go build -o /app/consumer /app/cmd/consumer/main.go" --build.bin "/app/consumer"
```

#### Menggunakan Air di Docker Container

1. Jalankan container development:
```bash
docker-compose up -d producer consumer
```

2. Air akan secara otomatis:
   - Monitor perubahan file
   - Rebuild aplikasi saat ada perubahan
   - Restart aplikasi dengan binary terbaru
   - Menampilkan log real-time

3. Konfigurasi Air dapat disesuaikan di file [`.air.toml`](.air.toml):
```toml
[build]
cmd = "go build -o /app/consumer /app/cmd/consumer/main.go"
bin = "/app/consumer"

[log]
time = false
main_only = false
silent = false

[color]
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"
```

### Langkah-langkah

1. Clone repository:
```bash
git clone <repository-url>
cd pubsub-ckg-tb
```

2. Install dependencies:
```bash
go mod download
```

3. Konfigurasi environment:
```bash
cp .env.example .env
# Edit .env dengan konfigurasi Anda
```

4. Buat topik dan subscription di Google Cloud Pub/Sub:
```bash
gcloud pubsub topics create ckg-tb-topic
gcloud pubsub subscriptions create ckg-tb-subscription --topic=ckg-tb-topic
```

## Penggunaan

### Menjalankan Producer Lokal

```bash
# Mode one-time production
go run cmd/producer/main.go

# Mode watch (untuk MongoDB)
# Akan secara otomatis mendeteksi perubahan data
go run cmd/producer/main.go
```

### Menjalankan Producer di Docker Container

```bash
# Jalankan producer di container
docker-compose up -d producer

# Lihat logs producer
docker-compose logs -f producer

# Stop producer
docker-compose stop producer
```

### Menjalankan Consumer Lokal

```bash
# Consumer akan berjalan terus-menerus
go run cmd/consumer/main.go
```

### Menjalankan Consumer di Docker Container

```bash
# Jalankan consumer di container
docker-compose up -d consumer

# Lihat logs consumer
docker-compose logs -f consumer

# Stop consumer
docker-compose stop consumer
```

### Menjalankan Semua Service

```bash
# Jalankan semua service (MongoDB, Pub/Sub Emulator, Producer, Consumer)
docker-compose up -d

# Lihat semua logs
docker-compose logs -f

# Stop semua service
docker-compose down
```

### Development dengan Live Reload

```bash
# Development mode dengan live reload untuk producer
docker-compose up -d producer

# Development mode dengan live reload untuk consumer
docker-compose up -d consumer

# Development mode untuk kedua service
docker-compose up -d producer consumer
```

Perubahan kode akan secara otomatis terdeteksi dan aplikasi akan di-rebuild serta di-restart tanpa perlu stop dan start manual.

## Model Data

### SkriningCKGRaw
Data mentah dari database CKG berisi informasi lengkap pasien dan hasil pemeriksaan.

### SkriningCKGResult
Data hasil yang telah diproses dan siap dikirim ke SITB.

### StatusPasien
Data status pasien TB yang diterima dari SITB.

## Database Schema

Project mendukung dua jenis database:

### MongoDB
- Menggunakan change stream untuk monitoring real-time
- Data disimpan dalam format BSON
- Mendukukung skema fleksibel

### SQL (MySQL/PostgreSQL)
- Menggunakan schema yang terstruktur
- Data disimpan dalam format JSON untuk field dinamis
- Lihat `schema-sitb-ckg.sql` untuk detail struktur

## Monitoring dan Logging

Sistem menggunakan structured logging dengan level:
- DEBUG: Informasi detail untuk debugging
- INFO: Informasi umum tentang operasi sistem
- WARN: Peringatan untuk kondisi yang tidak normal
- ERROR: Kesalahan yang memerlukan perhatian

Log dapat dikonfigurasi melalui environment variable `APP_LOGLEVEL`.

## Error Handling

Sistem memiliki mekanisme error handling yang komprehensif:
- Retry mechanism untuk transient errors
- Dead letter queue untuk pesan yang gagal diproses
- Graceful shutdown untuk aplikasi
- Connection pooling untuk database

## Testing

Untuk menjalankan test:
```bash
go test ./...
```

Untuk coverage report:
```bash
go test -cover ./...
```

## Contributing

1. Fork repository
2. Buat feature branch
3. Commit perubahan
4. Push ke branch
5. Buat pull request

### Development Workflow dengan Docker

1. Buat branch baru:
```bash
git checkout -b feature/nama-fitur
```

2. Lakukan perubahan kode

3. Test dengan Docker Compose:
```bash
# Rebuild dan jalankan service
docker-compose up -d --build producer consumer

# Lihat logs untuk debugging
docker-compose logs -f producer
```

4. Commit perubahan:
```bash
git add .
git commit -m "feat: tambah fitur baru"
git push origin feature/nama-fitur
```

## Lisensi

Lihat file [LICENSE](LICENSE) untuk informasi lisensi.

## Support

Untuk pertanyaan atau bantuan, silakan hubungi tim development atau buat issue di repository.