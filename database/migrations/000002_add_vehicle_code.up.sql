ALTER TABLE vehicles ADD vehicle_code VARCHAR(20);
CREATE UNIQUE INDEX idx_vehicles_vehicle_code ON vehicles(vehicle_code);
