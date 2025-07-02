# Troubleshooting Guide

## Masalah Go Version

### Error: "go.mod requires go >= 1.23.0"

**Penyebab**: Versi Go di go.mod tidak kompatibel dengan versi Go di Docker image.

**Solusi**: 
1. Pastikan versi Go di `go.mod` sudah diubah ke `go 1.22`
2. Rebuild Docker image:

```bash
docker-compose down
docker-compose up --build -d
```

**Alternatif**: Update Dockerfile untuk menggunakan Go versi yang lebih baru:

```dockerfile
FROM golang:1.23-alpine
```

## Masalah Docker Permission

### Error: "permission denied" saat build Docker

**Penyebab**: Permission issue dengan Docker daemon atau buildx.

**Solusi**:

1. **Restart Docker Desktop** (jika menggunakan macOS/Windows)
2. **Reset Docker buildx**:
```bash
docker buildx rm default
docker buildx create --use --name default
```

3. **Jalankan dengan sudo** (Linux):
```bash
sudo docker-compose up --build -d
```

4. **Alternatif tanpa Docker**:
```bash
# Install dependencies lokal
go mod download

# Set environment variables
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=fleet_management
export DB_PORT=5432
export RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Jalankan PostgreSQL dan RabbitMQ dengan Docker
docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:15
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
docker run -d --name mqtt -p 1883:1883 eclipse-mosquitto

# Jalankan aplikasi
go run main.go
```

## Masalah Database Connection

### Error: "connection refused" ke PostgreSQL

**Solusi**:
1. Pastikan PostgreSQL container berjalan:
```bash
docker-compose ps
```

2. Cek logs PostgreSQL:
```bash
docker-compose logs postgres
```

3. Reset database:
```bash
docker-compose down -v
docker-compose up -d
```

## Masalah MQTT Connection

### Error: "connection refused" ke MQTT broker

**Solusi**:
1. Cek status MQTT container:
```bash
docker-compose logs mqtt
```

2. Test koneksi MQTT:
```bash
# Install mosquitto clients
brew install mosquitto  # macOS
apt-get install mosquitto-clients  # Ubuntu

# Test publish
mosquitto_pub -h localhost -p 1883 -t test -m "hello"

# Test subscribe
mosquitto_sub -h localhost -p 1883 -t test
```

## Masalah RabbitMQ

### Error: "connection refused" ke RabbitMQ

**Solusi**:
1. Cek status RabbitMQ:
```bash
docker-compose logs rabbitmq
```

2. Akses RabbitMQ Management UI:
http://localhost:15672 (guest/guest)

3. Reset RabbitMQ:
```bash
docker-compose restart rabbitmq
```

## Masalah API

### Error: "404 Not Found" pada endpoint

**Solusi**:
1. Pastikan aplikasi berjalan di port 3000:
```bash
curl http://localhost:3000/vehicles
```

2. Cek logs aplikasi:
```bash
docker-compose logs app
```

3. Verifikasi routes di `routes/api.go`

## Tips Debugging

1. **Cek semua containers**:
```bash
docker-compose ps
```

2. **Lihat logs semua services**:
```bash
docker-compose logs -f
```

3. **Masuk ke container untuk debugging**:
```bash
docker-compose exec app sh
```

4. **Reset semua (nuclear option)**:
```bash
docker-compose down -v
docker system prune -a
docker-compose up --build -d
