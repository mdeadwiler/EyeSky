package opensky

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mdeadwiler/EyeSky/internal/domain"
)


type Client struct {
	httpClient *http.Client
	baseUrl string
	username string
	password string
	timeout time.Duration
	rateLimit int 
}
