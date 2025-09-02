package usecase

import (
	"context"
	"fmt"

	"github.com/mdeadwiler/EyeSky/internal/domain"
	"github.com/mdeadwiler/EyeSky/internal/platform/logging"
)

type FlightService struct {
	repo domain.FlightRepository
	openskyClient domain.OpenSkyClient
	logger *logging.Logger
}

func NewFlightService(repo domain.FlightRepository, client domain.OpenSkyClient, logger *logging.Logger) *FlightService {
	return &FlightService{
		repo: repo,
		openskyClient: client,
		logger: logger,
	}
}

func (fs *FlightService) FetchAndStoreFlights(ctx context.Context) error {
	// Fetch flights
	flights, err := fs.openskyClient.GetAllStates(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch flights: %w", err)
	}

	// Filter valid flights
	validFlights := make([]*domain.Flight, 0, len(flights))
	for _, flight := range flights {
		if flight.IsValid() {
			validFlights = append(validFlights, flight)
		}
	}
	// Store flights in DB
	if err := fs.repo.StoreBatch(ctx, validFlights); err != nil {
		return fmt.Errorf("failed to store flights: %w", err)
	}

	fs.logger.InfoWithFields("Flights processed", map[string]interface{}{
		"total": len(flights),
		"valid": len(validFlights),
	})
	return nil
}
