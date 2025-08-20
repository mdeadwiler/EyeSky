package domain

import (
	"time"
)

// Aircract Identification
type Flight struct {
	ICAO24 string `json:"icao24"` // Transponder address
	Callsign *string `json:"callsign"` // This can be null for flight 
	Country string `json:"country"` // Registration location
}