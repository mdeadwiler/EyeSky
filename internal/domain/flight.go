package domain

import (
	"time"
)


// Emergency squawk codes 
const (
	SquawkEmergency = "7700"
	SquawkRadioFailure = "7600"
	SquawkHijack = "7500"
	SquawkVFR = "1200"
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

	// Aircraft State 
	OnGround bool `json:"on_ground"`
	Alert bool `json:"alert"`
	SPI bool `json:"spi"` // Special Purpose Indicator 
	SquawkCode *string `json:"squawk"` // Transponder code

	// Temporal Data
	TimePosition *time.Time `json:"time_position"` // Last position update
	LastContact time.Time `json:"last_contact"`
	LastSeen *time.Time `json:"last_seen"`
}

func (f *Flight) IsValid() bool {
	// Mandatory for tracking flight
	if f.ICAO24 == "" {
		return false
	}
	// Correct format
	if !f.IsValidICAO24() {
		return false
	}
	// Fresh data
	if f.LastContact.IsZero() {
		return false
	}
	// Position Data
	if !f.HasPosition() {
		return false
	}
	return true

}
// Validates ICAO24 transponder address
func (f *Flight) IsValidICAO24() bool {
	if len(f.ICAO24) != 6 {
		return false
	}
	for _, char := range f.ICAO24 {
		if !((char >= '0' && char <= '9') ||
		    (char >= 'A' && char <= 'F') ||
		    (char >= 'a' && char <= 'f')) {
				return false
			}
	}
	return true
}

// Valid coordinate data
func (f *Flight) HasPosition() bool {
	return f.Latitude != nil && f.Longitude != nil &&
		   *f.Latitude >= -90 && *f.Latitude <= 90 &&
		   *f.Longitude >= -180 && *f.Longitude <= 180
}
