package usecase

import (
	"context"
	"fmt"
	"time"

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

func (fs *FlightService) FetchAndStoreFlightsInRegion(ctx context.Context, bounds domain.GeoBounds) error {
	// Fetch flights in region from OpenSky API
	flights, err := fs.openskyClient.GetStatesInRegion(ctx, bounds)
	if err != nil {
		return fmt.Errorf("failed to fetch flights in region: %w", err)
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
		return fmt.Errorf("failed to store regional flights: %w", err)
	}

	fs.logger.InfoWithFields("Regional flights processed", map[string]interface{}{
		"total": len(flights),
		"valid": len(validFlights),
	})

	return nil
}
func (fs *FlightService) GetActiveFlights(ctx context.Context, limit int) ([]*domain.Flight, error) {
	// Get only recent flights (last 5 minutes) - database-level filtering
	since := time.Now().Add(-5 * time.Minute)
	flights, err := fs.repo.GetInRegion(ctx, domain.GeoBounds{
		NorthLat: 90, SouthLat: -90, WestLon: -180, EastLon: 180, // Global bounds
	}, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get active flights: %w", err)
	}

	// Apply limit
	if len(flights) > limit {
		flights = flights[:limit]
	}

	return flights, nil
}

func (fs *FlightService) GetEmergencyFlights(ctx context.Context) ([]*domain.Flight, error) {
	// Get only live flights (last 2 minutes) for emergency detection
	since := time.Now().Add(-2 * time.Minute)
	flights, err := fs.repo.GetInRegion(ctx, domain.GeoBounds{
		NorthLat: 90, SouthLat: -90, WestLon: -180, EastLon: 180,
	}, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get flights for emergency check: %w", err)
	}

	// Filter for emergency flights only
	emergencyFlights := make([]*domain.Flight, 0)
	for _, flight := range flights {
		if flight.IsEmergency() {
			emergencyFlights = append(emergencyFlights, flight)
		}
	}

	if len(emergencyFlights) > 0 {
		fs.logger.InfoWithFields("Emergency flights detected", map[string]interface{}{
			"count": len(emergencyFlights),
		})
	}

	return emergencyFlights, nil
}

func (fs *FlightService) GetFlightsInRegion(ctx context.Context, bounds domain.GeoBounds) ([]*domain.Flight, error) {
	// Database already filters by time, no need for additional staleness check
	since := time.Now().Add(-2 * time.Minute) // Live data only
	flights, err := fs.repo.GetInRegion(ctx, bounds, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get flights in region: %w", err)
	}

	return flights, nil
}

func (fs *FlightService) CleanUpStaleData(ctx context.Context) (int, error) {
	// Delete flights older than 10 minutes
	cutoffTime := time.Now().Add(-10 * time.Minute)
	
	deletedCount, err := fs.repo.DeleteStale(ctx, cutoffTime)
	if err != nil {
		return 0, fmt.Errorf("failed to clean up stale data: %w", err)
	}

	if deletedCount > 0 {
		fs.logger.InfoWithFields("Stale data cleaned up", map[string]interface{}{
			"deleted_count": deletedCount,
			"cutoff_time":   cutoffTime,
		})
	}

	return deletedCount, nil
}

func (fs *FlightService) GetFlightStatistics(ctx context.Context) (*domain.FlightStats, error) {
	// Get only live flights for accurate statistics
	since := time.Now().Add(-2 * time.Minute)
	liveFlights, err := fs.repo.GetInRegion(ctx, domain.GeoBounds{
		NorthLat: 90, SouthLat: -90, WestLon: -180, EastLon: 180,
	}, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get live flights for statistics: %w", err)
	}

	// Calculate real-time statistics
	stats := &domain.FlightStats{
		TotalFlights: len(liveFlights), // Live count, not historical
		ActiveFlights: len(liveFlights),
		LastUpdated:  time.Now(),
	}

	for _, flight := range liveFlights {
		if flight.OnGround {
			stats.OnGroundFlights++
		} else {
			stats.AirborneFlights++
		}

		if flight.IsEmergency() {
			stats.EmergencyFlights++
		}
	}

	return stats, nil
}