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

func (c *Client) GetAllStates(ctx context.Context) ([] *domain.Flight, error) {
	// Request URL
	reqURL := fmt.Sprintf("%s/states/all" , c.baseURL)

	// HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	//OAuth
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	// HTTP Request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenSky API return status %d", resp.StatusCode)
	}
	// Parse JSON response
	var apiResponse StateVectorResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse);
	err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to domain flights
	flights := make([]*domain.Flight, 0, len(apiResponse.States))
	for _, state := range apiResponse.States {
		if flight, err := c.parseStateVector(state);
		err == nil {
			flights = append(flights, flight)
		}
	}
	return flights, nil
}

func (c *Client) parseStateVector(state []interface{}) (*domain.Flight, error) {
	//min len of array(data)
	if len(state) < 17 {
		return nil, fmt.Errorf("invalid state vector: expected 17 fields, got %d", len(state))
	}
}