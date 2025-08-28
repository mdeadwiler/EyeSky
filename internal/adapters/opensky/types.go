package opensky


type StateVectorResponse struct {
	Time int64 `json:"time"`
	States [][]interface{} `json:"states"`
}

// OpenSky state vector field positions in the array response
const (
	ICAO24Index         = 0  // ICAO24 address (string)
	CallsignIndex       = 1  // Callsign (string, can be null)
	OriginCountryIndex  = 2  // Country of registration (string)
	TimePositionIndex   = 3  // Unix timestamp of position (int, can be null)
	LastContactIndex    = 4  // Unix timestamp of last contact (int)
	LongitudeIndex      = 5  // Longitude in decimal degrees (float, can be null)
	LatitudeIndex       = 6  // Latitude in decimal degrees (float, can be null)
	BaroAltitudeIndex   = 7  // Barometric altitude in meters (float, can be null)
	OnGroundIndex       = 8  // On ground status (bool)
	VelocityIndex       = 9  // Ground speed in m/s (float, can be null)
	TrueTrackIndex      = 10 // Track angle in decimal degrees (float, can be null)
	VerticalRateIndex   = 11 // Vertical rate in m/s (float, can be null)
	SensorsIndex        = 12 // IDs of sensors (array, can be null)
	GeoAltitudeIndex    = 13 // Geometric altitude in meters (float, can be null)
	SquawkIndex         = 14 // Squawk code (string, can be null)
	SPIIndex            = 15 // Special purpose indicator (bool)
	PositionSourceIndex = 16 // Position source (int)
)
