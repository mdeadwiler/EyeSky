-- Down migration for initial schema
-- Drop tables in reverse order to handle dependencies

DROP INDEX IF EXISTS idx_flights_domestic_active;
DROP INDEX IF EXISTS idx_flights_active;
DROP INDEX IF EXISTS idx_flights_squawk_emergency;
DROP INDEX IF EXISTS idx_flights_coordinates;
DROP INDEX IF EXISTS idx_flights_icao24;
DROP TABLE IF EXISTS flights;

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;
