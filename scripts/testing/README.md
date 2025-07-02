# Testing Scripts

## Test Location Updates

Untuk menguji pengiriman dan penerimaan data lokasi kendaraan:

1. Pastikan aplikasi utama sudah berjalan:
```bash
docker-compose up -d
```

2. Jalankan script test location:
```bash
cd scripts/testing/test_location
go run main.go
```

3. Periksa log aplikasi utama untuk memastikan pesan lokasi diterima dan disimpan ke database.

4. Gunakan API endpoint untuk memverifikasi data lokasi:
- GET /api/vehicles/:id/locations
- GET /api/vehicles/:id/last-location

Pastikan vehicle_code yang dikirim sesuai dengan data kendaraan di database.
