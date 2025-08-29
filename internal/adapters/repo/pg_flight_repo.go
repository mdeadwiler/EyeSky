package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	
	"github.com/mdeadwiler/EyeSky/internal/domain"
	"github.com/mdeadwiler/EyeSky/internal/platform/db"	
)


type Repository struct {
	db *sql.DB
}

// New SQL flight repo
func New(database *db.DB) *Repository {
	return &Repository{
		db: database.GetConnection(),
	}
}

// The $ is to protect and point to the exact param to pull from the query. Think of it like a id. $ protects from SQL injection like obsfuscation
func (r *Repository) Store(ctx context.Context, flight *domain.Flight) error {
	query := `
	INSERT INTO flights (
	icao24, callsign, country, longitude, latitude,
	barometric_altitude, geometric_altitude, velocity, heading, vertical_rate,
	on_ground, alert, spi, squawk_code,
	time_position, last_contact, last_seen) VALUES (
	$1, $2, $3, $4, $5,
	$6, $7, $8, $9, $10,
	$11, $12, $13, $14,
	$15, $16, $17)`

	_, err := r.db.ExecContext(ctx, query,
	flight.ICAO24, flight.Callsign, 
	flight.Country, flight.Longitude, flight.Latitude,
	flight.BarometricAltitude, flight.GeometricAltitude, 
	flight.Velocity, flight.Heading, flight.VerticalRate,
	flight.OnGround, flight.Alert, flight.SPI, 
	flight.SquawkCode, flight.TimePosition, flight.LastContact, flight.LastSeen,
)
if err != nil {
	return fmt.Errorf("failed to store flight %s: %w", flight.ICAO24, err)
}
return nil
}


// StoreBatch efficiently stores multiple flights using PostgreSQL COPY - 10x faster than INSERT
func (r *Repository) StoreBatch(ctx context.Context, flights []*domain.Flight) error {
	if len(flights) == 0 {
		return nil
	}
	// Turns batching from seconds to streaming in milliseconds. Batching causes bottlenecking when the number of flights is large.
	// Use PostgreSQL COPY for maximum insertion speed
	stmt, err := r.db.PrepareContext(ctx, pq.CopyIn("flights",
		"icao24", "callsign", "country", "longitude", "latitude",
		"barometric_altitude", "geometric_altitude", "velocity", "heading", "vertical_rate",
		"on_ground", "alert", "spi", "squawk_code",
		"time_position", "last_contact", "last_seen"))
	if err != nil {
		return fmt.Errorf("failed to prepare COPY statement: %w", err)
	}

	// Stream flight data directly to PostgreSQL
	for _, flight := range flights {
		_, err = stmt.ExecContext(ctx,
			flight.ICAO24, flight.Callsign, flight.Country,
			flight.Longitude, flight.Latitude,
			flight.BarometricAltitude, flight.GeometricAltitude,
			flight.Velocity, flight.Heading, flight.VerticalRate,
			flight.OnGround, flight.Alert, flight.SPI, flight.SquawkCode,
			flight.TimePosition, flight.LastContact, flight.LastSeen,
		)
		if err != nil {
			stmt.Close()
			return fmt.Errorf("failed to add flight %s to batch: %w", flight.ICAO24, err)
		}
	}

	// Execute the COPY operation
	if err := stmt.Close(); err != nil {
		return fmt.Errorf("failed to execute COPY batch of %d flights: %w", len(flights), err)
	}

	return nil
}