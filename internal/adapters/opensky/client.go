package opensky

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mdeadwiler/EyeSky/internal/domain"
	"github.com/mdeadwiler/EyeSky/internal/platform/config"
)


type Client struct {
	httpClient *http.Client
	baseURL string
	username string
	password string
	timeout time.Duration
	rateLimit int 
}

// New OpenSky API client with configuraton 
func NewClient(cfg config.OpenSkyConfig) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		baseURL: cfg.BaseURL,
		username: cfg.Username,
		password: cfg.Password,
		timeout: cfg.Timeout,
		rateLimit: cfg.RateLimit,
	}
}

// Common HTTP request method - eliminates code duplication
func (c *Client) makeRequest(ctx context.Context, reqURL string) (*StateVectorResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Basic Authentication (Note: OpenSky also supports OAuth2 if preferred)
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenSky API returned status %d", resp.StatusCode)
	}

	var apiResponse StateVectorResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResponse, nil
}

// Convert API response to domain flights - optimized parsing
func (c *Client) convertToFlights(apiResponse *StateVectorResponse) []*domain.Flight {
	if len(apiResponse.States) == 0 {
		return nil
	}

	// Pre-allocate with exact capacity for better performance
	flights := make([]*domain.Flight, 0, len(apiResponse.States))
	
	for _, state := range apiResponse.States {
		if flight := c.parseStateVectorOptimized(state); flight != nil {
			flights = append(flights, flight)
		}
	}
	
	return flights
}

func (c *Client) GetAllStates(ctx context.Context) ([]*domain.Flight, error) {
	reqURL := c.baseURL + "/states/all"
	
	apiResponse, err := c.makeRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	return c.convertToFlights(apiResponse), nil
}

// Optimized US domestic flights - most common use case
func (c *Client) GetDomesticStates(ctx context.Context) ([]*domain.Flight, error) {
	// Pre-built US bounds for maximum performance
	reqURL := c.baseURL + "/states/all?lamin=24.000000&lamax=49.000000&lomin=-125.000000&lomax=-66.000000"
	
	apiResponse, err := c.makeRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	return c.convertToFlights(apiResponse), nil
}

// Optimized parsing - returns nil for invalid flights instead of error
func (c *Client) parseStateVectorOptimized(state []interface{}) *domain.Flight {
	// Fast length check
	if len(state) < 17 {
		return nil
	}
	
	// Required ICAO24 - fail fast if missing
	icao24, ok := state[ICAO24Index].(string)
	if !ok || icao24 == "" {
		return nil
	}

	// Initialize flight with required field
	flight := &domain.Flight{ICAO24: icao24}

	// Direct field extraction - maximum performance
	if callsign, ok := state[CallsignIndex].(string); ok && callsign != "" {
		flight.Callsign = &callsign
	}
	
	if country, ok := state[OriginCountryIndex].(string); ok {
		flight.Country = country
	}
	
	if lastContact, ok := state[LastContactIndex].(float64); ok {
		flight.LastContact = time.Unix(int64(lastContact), 0)
	}
	
	if longitude, ok := state[LongitudeIndex].(float64); ok {
		flight.Longitude = &longitude
	}
	
	if latitude, ok := state[LatitudeIndex].(float64); ok {
		flight.Latitude = &latitude
	}
	
	if baroAlt, ok := state[BaroAltitudeIndex].(float64); ok {
		altFeet := baroAlt * 3.28084 // Convert meters to feet
		flight.BarometricAltitude = &altFeet
	}

	return flight
}



// GetStatesInRegion fetches aircraft states within geographic bounds - optimized
func (c *Client) GetStatesInRegion(ctx context.Context, bounds domain.GeoBounds) ([]*domain.Flight, error) {
	// Build URL with parameters - single allocation
	reqURL := fmt.Sprintf("%s/states/all?lamin=%.6f&lamax=%.6f&lomin=%.6f&lomax=%.6f",
		c.baseURL, bounds.SouthLat, bounds.NorthLat, bounds.WestLon, bounds.EastLon)

	apiResponse, err := c.makeRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	return c.convertToFlights(apiResponse), nil
}

// GetStatesByICAO24 fetches states for specific aircraft - optimized
func (c *Client) GetStatesByICAO24(ctx context.Context, icao24s []string) ([]*domain.Flight, error) {
	if len(icao24s) == 0 {
		return nil, nil
	}

	// Efficient URL building with pre-allocated capacity
	params := make(url.Values, len(icao24s))
	for _, icao24 := range icao24s {
		params.Add("icao24", icao24)
	}

	reqURL := c.baseURL + "/states/all?" + params.Encode()

	apiResponse, err := c.makeRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	return c.convertToFlights(apiResponse), nil
}

// HealthCheck verifies OpenSky API is accessible - optimized
func (c *Client) HealthCheck(ctx context.Context) error {
	reqURL := c.baseURL + "/states/all"
	
	req, err := http.NewRequestWithContext(ctx, "HEAD", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}
	
	// Add auth for health check if available
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("OpenSky API unreachable: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return fmt.Errorf("OpenSky API unhealthy, status: %d", resp.StatusCode)
	}
	
	return nil
}