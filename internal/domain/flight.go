package domain

import (
	"time"
)

// Aircract Identification
type Flight struct {
	ICAO24 string `json:"icao24"` // Transponder address
	Callsign *string `json:"callsign"` // This can be null for flight 
	Country string `json:"country"` // Registration location

	// Position Data
	Longitude *float64 `json:"longitude"` // this will be in decimal degrees
	Latitude *float64 `json:"latitude"` // this will be in decimal degrees

	// Altitude Data and this will be converted in feet in the adapter layer
	BarometricAltitude *float64 `json:"barometric_altitude"` 
	GeometricAltitude *float64 `json:"geometric_altitude"`

	// Movement Data
Velocity *float64 `json:"velocity"`
Heading *float64 `json:"heading"`
VerticalRate *float64 `json:"vertical_rate"`
}

