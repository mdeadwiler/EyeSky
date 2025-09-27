# EyeSky

A real-time flight tracking and monitoring system built in Go that fetches live aircraft data from the OpenSky Network API.

## Features

- **Real-time Flight Tracking**: Fetches live aircraft data from OpenSky Network
- **Emergency Detection**: Monitors special transponder codes (7700=Emergency, 7600=Radio Failure, 7500=Hijack)  
- **Geographic Filtering**: Track flights within specific regions using bounding box queries
- **Data Validation**: Validates ICAO24 transponder addresses and flight coordinates
- **PostgreSQL Storage**: Persistent storage for flight data and analytics
- **Clean Architecture**: Domain-driven design with separated business logic

## Tech Stack

- **Language**: Go 1.25
- **Database**: PostgreSQL
- **External API**: OpenSky Network REST API
- **Logging**: Zerolog
- **Environment**: Docker Compose

## Project Structure

```
├── cmd/gateway/          # Application entry point
├── internal/
│   ├── domain/           # Business logic and entities
│   ├── usecase/          # Application services
│   ├── adapters/         # External API clients
│   └── platform/         # Infrastructure (DB, config, logging)
├── configs/              # Configuration files
├── deploy/               # Deployment configurations
└── scripts/              # Development scripts
```

## Prerequisites

- Go 1.25+
- Docker & Docker Compose
- PostgreSQL

## Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/mdeadwiler/EyeSky.git
   cd EyeSky
   ```

2. **Create environment file**
   ```bash
   cp .env.example .env
   ```
   
   Configure your `.env` file with:
   ```env
   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=eyesky_user
   DB_PASSWORD=your_password
   DB_NAME=eyesky
   
   # OpenSky Network API (optional - higher rate limits with account)
   OPEN_SKY_USERNAME=your_username
   OPEN_SKY_PASSWORD=your_password
   
   # Server
   PORT=8080
   ENVIRONMENT=development
   ```

3. **Start PostgreSQL with Docker**
   ```bash
   docker-compose -f docker.yaml up -d
   ```

4. **Install dependencies**
   ```bash
   go mod download
   ```

5. **Run the application**
   ```bash
   go run cmd/gateway/main.go
   ```

## Development Status

**This project is currently in development**

The core domain logic and database integration are complete, but the HTTP API endpoints are still being implemented.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is open source. See LICENSE file for details.
