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
}
