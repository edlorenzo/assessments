CREATE TABLE weathers (
     id SERIAL PRIMARY KEY,
     timestamp TIMESTAMP NOT NULL,
     city TEXT NOT NULL,
     temperature FLOAT NOT NULL,
     humidity FLOAT NOT NULL,
     wind_speed FLOAT NOT NULL
);

CREATE INDEX idx_id_weathers ON weathers (id);
CREATE INDEX idx_timestamp_weathers ON weathers (timestamp);
