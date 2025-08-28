-- DB for EyeSkye flight tracking 
-- This is for domestic US Flight data


-- Users Table - username login with email registration
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,     -- Login
    email VARCHAR(255) UNIQUE NOT NULL,       -- Registration and notifications
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    active BOOLEAN DEFAULT TRUE
);

-- Primary index for login queries (most frequent)
CREATE INDEX idx_users_username ON users(username);
-- Secondary index for registration validation
CREATE INDEX idx_users_email ON users(email);



CREATE TABLE flights (
    id BIGSERIAL PRIMARY KEY,

    -- Aircraft Id
    icao24 VARCHAR(6) NOT NULL CHECK (length(icao24) = 6),
    callsign VARCHAR(8),
    country VARCHAR(100) NOT NULL,

    -- Position Data tracking
    longitude DECIMAL(10,6) CHECK (longitude >= -180 AND longitude <= 180),
    latitude DECIMAL(10,6) CHECK (latitude >= -90 AND latitude <= 90),

    barometric_altitude REAL CHECK (barometric_altitude >= -1000 AND barometric_altitude <= 60000),
    geometric_altitude REAL CHECK (geometric_altitude >= -1000 AND geometric_altitude <= 60000),
    
    -- Movement Data
    velocity REAL CHECK (velocity >= 0 AND velocity <= 1500), -- Max ~Mach 4.5
    heading REAL CHECK (heading >= 0 AND heading < 360),
    vertical_rate REAL CHECK (vertical_rate >= -10000 AND vertical_rate <= 10000), -- ft/min
    
    -- Aircraft State
    on_ground BOOLEAN NOT NULL DEFAULT FALSE,
    alert BOOLEAN NOT NULL DEFAULT FALSE,
    spi BOOLEAN NOT NULL DEFAULT FALSE,
    squawk_code VARCHAR(4) CHECK (squawk_code ~ '^[0-7]{4}$'), -- Only octal digits
    
    -- Temporal Data
    time_position TIMESTAMP WITH TIME ZONE,
    last_contact TIMESTAMP WITH TIME ZONE NOT NULL,
    last_seen TIMESTAMP WITH TIME ZONE,
    
    -- Metadata for tracking
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for flight domestic flights only 
CREATE INDEX idx_flights_icao24 ON flights(icao24);
CREATE INDEX idx_flights_coordinates ON flights(latitude, longitude) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;
CREATE INDEX idx_flights_squawk_emergency ON flights(squawk_code) WHERE squawk_code IN ('7700', '7600', '7500');
CREATE INDEX idx_flights_active ON flights(last_contact) WHERE last_contact > NOW() - INTERVAL '5 minutes';

-- Composite index for common domestic queries
CREATE INDEX idx_flights_domestic_active ON flights(last_contact, latitude, longitude) 
    WHERE latitude BETWEEN 24 AND 49 AND longitude BETWEEN -125 AND -66 AND last_contact > NOW() - INTERVAL '5 minutes';
