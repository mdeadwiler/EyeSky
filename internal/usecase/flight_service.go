package usecase

import (
	"context"

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
