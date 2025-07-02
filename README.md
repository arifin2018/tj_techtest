# Fleet Management System

Sistem manajemen armada kendaraan yang dibangun dengan Go Fiber, PostgreSQL, MQTT, dan RabbitMQ.

## Quick Start

### Prasyarat
- Docker dan Docker Compose terinstall
- Git untuk clone repository
- Go 1.22+ (jika ingin menjalankan tanpa Docker)

### Langkah-langkah

1. Clone repository dan masuk ke direktori proyek
```bash
git clone <repository-url>
cd tj_techtest
```

2. Download dependencies Go (opsional, untuk development)
```bash
go mod tidy
```

3. Jalankan aplikasi dengan Docker Compose
```bash
docker-compose up -d
```

**Catatan**: Jika mengalami masalah Docker permission, lihat [TROUBLESHOOTING.md](TROUBLESHOOTING.md) untuk solusi alternatif.

4. Verifikasi services berjalan
- API: http://localhost:3000
- RabbitMQ Management: http://localhost:15672 (guest/guest)
- PostgreSQL: localhost:5432
- MQTT: localhost:1883

5. Jalankan script publisher untuk simulasi data lokasi
```bash
go run scripts/mqtt_publisher/main.go
```

6. Test API dengan Postman Collection yang tersedia di `postman/Fleet_Management_API.postman_collection.json`

## Dokumentasi Lengkap

Untuk dokumentasi lebih detail termasuk:
- Struktur proyek
- API endpoints
- Format data
- Cara testing
- Troubleshooting
- Dan lainnya

Silakan baca [README_LENGKAP.md](README_LENGKAP.md)

## License

MIT License
