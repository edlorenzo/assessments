CREATE TABLE bike_availabilities (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP,
    name VARCHAR(255),
    total_docks INT,
    docks_available INT,
    bikes_available INT,
    classic_bikes_available INT,
    smart_bikes_available INT,
    electric_bikes_available INT,
    reward_bikes_available INT,
    reward_docks_available INT,
    kiosk_status VARCHAR(255),
    kiosk_public_status VARCHAR(255),
    kiosk_connection_status VARCHAR(255),
    kiosk_type INT,
    address_street VARCHAR(255),
    address_city VARCHAR(255),
    address_state VARCHAR(255),
    address_zip_code VARCHAR(255),
    is_event_based BOOLEAN,
    is_virtual BOOLEAN,
    kiosk_id INT,
    trikes_available INT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8)
);

CREATE INDEX idx_id_bike_availabilities ON bike_availabilities (id);
CREATE INDEX idx_timestamp_bike_availabilities ON bike_availabilities (timestamp);
CREATE INDEX idx_kiosk_id_bike_availabilities ON bike_availabilities (kiosk_id);