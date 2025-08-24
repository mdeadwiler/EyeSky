package domain

import (
	"context"
	"time"
)

// Data for flight repository
type FlightRepository interface {
	Store(ctx context.Context, flight *Flight) error
	StoreBatch(ctx context.Context, flights []*Flight) error 
}