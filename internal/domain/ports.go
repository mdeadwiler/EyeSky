package domain

import (
	"context"
	"time"
)


// Geographical bounding box for flight queries( Get flights by region, fetch flights, query flights to specific area, and OpenSky bounding box query efficiency)
type GeoBounds struct {
	NorthLat float64 `json:"north_lat"`
	SouthLat float64 `json:"south_lat"`
	WestLon float64 `json:"west_lon"`
	EastLon float64 `json:"east_lon"`
}

type FlightStats struct {
	TotalFlights int `json:"total_flights"`
	ActiveFlights int `json:"active_flights"`
	OnGroundFlights int `json:"on_ground_flights"`
	AirborneFlights int `json:"airborne_flights"`
	EmergencyFlights int `json:"emergency_flights"`
	LastUpdated time.Time `json:"last_updated"`
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
// Business logic for flight service...flight tracking 
type FlightService interface {
// Fetch and store flight data
FetchAndStoreFlights(ctx context.Context) error
FetchAndStoreFlightsInRegion(ctx context.Context, bounds GeoBounds) error

//Query operations with business logic
GetActiveFlights(ctx context.Context, limit int) ([]*Flight, error)
GetEmergencyFlights(ctx context.Context) ([]*Flight, error)
GetFlightsInRegion(ctx context.Context, bounds GeoBounds) ([]*Flight, error)

//Data management 
CleanUpStaleData(ctx context.Context) (int, error)
GetFlightStatistics(ctx context.Context) (*FlightStats, error)
}