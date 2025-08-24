package domain

import (
	"context"
	"time"
)


// Geographical bounding box for flight queries( Get flights by region, fetch flights, query flights to specific area, and OpenSky bounding box query efficiency)
type GeoBounds struct {
	NorthLat float64 `json:"north_lat"`
	SouthLat float64 `json:"south_lat"`
	WestLong float64 `json:"west_lon"`
	EastLong float64 `json:"east_lon"`
}

// Data for flight repository
type FlightRepository interface {
	Store(ctx context.Context, flight *Flight) error
	StoreBatch(ctx context.Context, flights []*Flight) error 

	// Queries 
	GetByICAO24(ctx context.Context, icao24 string) (*Flight, error)
	GetAll(ctx context.Context, limit int) ([]*Flight, error)
	GetInRegion(ctx context.Context, bounds GeoBounds, since time.Time) ([]*Flight, error)

	// Maintenance Ops
	DeleteStale(ctx context.Context, olderThan time.Time) (int, error)
	Count(ctx context.Context) (int, error)

}

type OpenSkyClient interface {
	// Get all aircraft states globally
	GetAllStates(ctx context.Context) ([]*Flight, error)

	//Get states in region
	GetStatesInRegion(ctx context.Context, bounds GeoBounds) ([]*Flight, error)

	//Get states by ICAO24..Specific aircraft
	GetStatesByICAO24(ctx context.Context, icao24s []string) ([]*Flight, error)

	//Health check for API
	HealthCheck(ctx context.Context) error 
}