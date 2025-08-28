-- DB for EyeSkye flight tracking 
-- This is for domestic US Flight data

CREATE TABLE flights {
    id BIGSERIAL PRIMARY KEY,

    -- Aircraft Id
    icao24 VARCHAR(6) NOT NULL CHECK (length(icao24) = 6),
    callsign VARCHAR(8),
    country VARCHAR(100) NOT NULL,
}