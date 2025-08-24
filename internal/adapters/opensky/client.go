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
