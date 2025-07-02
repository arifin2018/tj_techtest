# Fleet Management System

Sistem manajemen armada kendaraan yang dibangun dengan Go Fiber, PostgreSQL, MQTT, dan RabbitMQ untuk tes teknis Backend Engineer Transjakarta.

## Fitur

- ✅ Menerima data lokasi kendaraan melalui MQTT
- ✅ Menyimpan data lokasi ke PostgreSQL
- ✅ API untuk mendapatkan lokasi terkini dan riwayat perjalanan kendaraan
- ✅ Sistem geofence dengan notifikasi event melalui RabbitMQ
- ✅ Containerized dengan Docker untuk deployment yang mudah

## Teknologi yang Digunakan

- **Golang** - Backend development
- **Fiber** - Web framework
- **GORM** - ORM untuk database
- **PostgreSQL** - Database
- **MQTT (Eclipse Mosquitto)** - Menerima data lokasi kendaraan
- **RabbitMQ** - Message broker untuk event geofence
- **Docker & Docker Compose** - Containerization

## Struktur Proyek

```
tj_techtest/
├── app/
│   ├── http/
│   │   ├── controllers/     # HTTP controllers
│   │   └── middleware/      # HTTP middleware
│   ├── models/              # Database models
│   ├── providers/           # Service providers
│   └── services/            # Business logic services
├── config/                  # Configuration files
├── database/
│   ├── migrations/          # Database migrations
│   └── seeders/             # Database seeders
├── pkg/
│   └── mqtt/               # MQTT client
├── routes/                  # API routes
├── scripts/
│   ├── mqtt_publisher/     # Script untuk publish data lokasi
│   ├── mqtt_subscriber/    # Script untuk subscribe MQTT
│   └── testing/            # Script testing
├── postman/                # Postman collection untuk testing API
├── storage/
│   └── logs/               # Log files
├── docker-compose.yml
├── Dockerfile
└── main.go
```

## Cara Menjalankan Aplikasi

### Prasyarat

- Docker dan Docker Compose terinstall
- Git untuk clone repository

### Langkah 1: Clone Repository

```bash
git clone <repository-url>
cd tj_techtest
```

### Langkah 2: Jalankan dengan Docker Compose

```bash
# Jalankan semua services (PostgreSQL, RabbitMQ, MQTT, App)
docker-compose up -d

# Cek status containers
docker-compose ps

# Lihat logs jika diperlukan
docker-compose logs -f app
```

### Langkah 3: Verifikasi Services

Setelah menjalankan `docker-compose up -d`, pastikan semua services berjalan:

- **API Application**: http://localhost:3000
- **PostgreSQL**: localhost:5432
- **RabbitMQ Management**: http://localhost:15672 (username: guest, password: guest)
- **MQTT Broker**: localhost:1883

### Langkah 4: Testing Data Lokasi

#### 4.1 Buat Kendaraan Baru (Opsional)

```bash
curl -X POST http://localhost:3000/vehicles \
  -H "Content-Type: application/json" \
  -d '{"name": "Bus Transjakarta 001"}'
```

#### 4.2 Jalankan MQTT Publisher untuk Simulasi Data

```bash
# Masuk ke container aplikasi
docker-compose exec app sh

# Jalankan script publisher
go run scripts/mqtt_publisher/main.go
```

Atau jalankan di terminal terpisah (jika Go terinstall di host):

```bash
go run scripts/mqtt_publisher/main.go
```

Script ini akan mengirim data lokasi mock setiap 2 detik ke topik MQTT `/fleet/vehicle/B1234XYZ/location`.

#### 4.3 Test API Endpoints

```bash
# Get semua kendaraan
curl http://localhost:3000/vehicles

# Get lokasi terakhir kendaraan (ganti :id dengan ID kendaraan)
curl http://localhost:3000/vehicles/1/location

# Get riwayat lokasi dengan filter waktu
curl "http://localhost:3000/vehicles/1/history?start=1715000000&end=1715009999"
```

## Testing dengan Postman

### Import Collection

1. Buka Postman
2. Klik "Import" 
3. Pilih file `postman/Fleet_Management_API.postman_collection.json`
4. Collection "Fleet Management API" akan tersedia

### Contoh Testing Flow

1. **Setup**: Pastikan aplikasi berjalan dengan `docker-compose up -d`
2. **Create Vehicle**: Gunakan endpoint "Create Vehicle" untuk membuat kendaraan baru
3. **Start Data Stream**: Jalankan `go run scripts/mqtt_publisher/main.go`
4. **Test Last Location**: Gunakan endpoint "Get Vehicle Last Location"
5. **Test History**: Gunakan endpoint "Get Vehicle Location History" dengan parameter waktu

## API Endpoints

### Vehicles

- `GET /vehicles` - Mendapatkan semua kendaraan
- `POST /vehicles` - Membuat kendaraan baru
- `GET /vehicles/:id` - Mendapatkan detail kendaraan
- `GET /vehicles/:id/location` - Mendapatkan lokasi terakhir kendaraan
- `GET /vehicles/:id/history` - Mendapatkan riwayat lokasi kendaraan

### Geofences

- `GET /geofences` - Mendapatkan semua geofence
- `POST /geofences` - Membuat geofence baru
- `GET /geofences/:id` - Mendapatkan detail geofence
- `PUT /geofences/:id` - Update geofence
- `DELETE /geofences/:id` - Hapus geofence

## Integrasi MQTT

### Format Data Lokasi

Data lokasi dikirim ke topik `/fleet/vehicle/{vehicle_id}/location` dengan format JSON:

```json
{
  "vehicle_id": "B1234XYZ",
  "latitude": -6.2088,
  "longitude": 106.8456,
  "timestamp": 1715000000
}
```

### Validasi Data

- `vehicle_id` harus ada dan tidak kosong
- `latitude` dan `longitude` harus berupa angka valid
- `timestamp` harus berupa Unix timestamp

## Integrasi RabbitMQ

### Geofence Events

Event geofence akan dipublish ke exchange `fleet.events` dengan queue `geofence_alerts` ketika kendaraan memasuki area geofence dengan radius 50 meter.

Format event:

```json
{
  "vehicle_id": "B1234XYZ",
  "geofence_name": "Terminal Pulogadung",
  "event_type": "enter",
  "timestamp": 1715000000
}
```

## Database Schema

### Vehicles
- id (Primary Key)
- name (VARCHAR)
- latitude, longitude (DOUBLE PRECISION)
- last_seen (TIMESTAMP)
- created_at, updated_at, deleted_at

### Vehicle Locations
- id (Primary Key)
- vehicle_id (Foreign Key ke vehicles)
- latitude, longitude (DOUBLE PRECISION)
- timestamp (TIMESTAMP)

### Geofences
- id (Primary Key)
- name (VARCHAR)
- latitude, longitude (DOUBLE PRECISION)
- radius (DOUBLE PRECISION dalam meter)
- created_at, updated_at, deleted_at

## Troubleshooting

### Container Tidak Berjalan

```bash
# Cek status containers
docker-compose ps

# Restart services
docker-compose restart

# Rebuild jika ada perubahan code
docker-compose up --build -d
```

### Database Connection Error

```bash
# Cek logs PostgreSQL
docker-compose logs postgres

# Reset database
docker-compose down -v
docker-compose up -d
```

### MQTT Connection Error

```bash
# Cek logs MQTT broker
docker-compose logs mqtt

# Test koneksi MQTT
docker-compose exec mqtt mosquitto_pub -h localhost -t test -m "hello"
```

## Development

### Menjalankan Tanpa Docker

1. Install dependencies:
```bash
go mod download
```

2. Setup environment variables:
```bash
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=fleet_management
export DB_PORT=5432
export RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

3. Jalankan PostgreSQL dan RabbitMQ secara manual atau dengan Docker:
```bash
docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:15
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
docker run -d --name mqtt -p 1883:1883 eclipse-mosquitto
```

4. Jalankan aplikasi:
```bash
go run main.go
```

### Menambah Migration

1. Buat file migration baru di `database/migrations/`
2. Format: `YYYYMMDDHHMMSS_description.up.sql` dan `YYYYMMDDHHMMSS_description.down.sql`

### Testing

```bash
go test ./...
```

## Persyaratan Teknis yang Dipenuhi

✅ **Menerima data lokasi kendaraan melalui MQTT**
- Topic: `/fleet/vehicle/{vehicle_id}/location`
- Format JSON sesuai spesifikasi
- Validasi data input

✅ **Menyimpan data lokasi ke PostgreSQL**
- Tabel `vehicle_locations` dengan field sesuai spesifikasi
- Service untuk insert data ke database

✅ **API untuk mengakses data lokasi**
- GET `/vehicles/{vehicle_id}/location` - lokasi terakhir
- GET `/vehicles/{vehicle_id}/history` - riwayat dengan filter waktu
- Response format JSON sesuai spesifikasi

✅ **RabbitMQ untuk event geofence**
- Exchange: `fleet.events`
- Queue: `geofence_alerts`
- Event dikirim saat kendaraan masuk radius 50 meter

✅ **Docker untuk deployment**
- docker-compose.yml lengkap dengan semua services
- Environment variables untuk konfigurasi
- Volume untuk persistence data

✅ **Script MQTT Publisher**
- Script Go untuk publish data mock
- Mengirim data setiap 2 detik
- Format data sesuai spesifikasi

✅ **Postman Collection**
- Collection lengkap untuk testing semua endpoint
- Contoh request dan response
- Dokumentasi cara penggunaan

## Contributing

1. Fork repository
2. Buat feature branch
3. Commit changes
4. Push ke branch
5. Buat Pull Request

## License

MIT License
