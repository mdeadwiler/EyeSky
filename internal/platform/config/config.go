package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment string
	Server  ServerConfig
	Database DatabaseConfig
	OpenSky OpenSkyConfig
	Logging LoggingConfig
}

type ServerConfig struct {
	Port  string
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	IdleTimeout time.Duration
	ShutdownTimeout time.Duration
	TLSEnabled bool // HTTPS
	TLSCertFile string // Cert file path
	TLSKeyFile string // Private key file path
}

type DatabaseConfig struct {
	Host string 
	Port string
	User string
	Password string
	DBName string
	SSLMode string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime time.Duration 
}

type OpenSkyConfig struct {
	BaseURL string
	Username string
	Password string
	Timeout time.Duration
	RateLimit int // This will be request per minute
}

type LoggingConfig struct {
	Level string
	Format string // "json"
}

func New()(*Config, error) {
	env := getEnv("ENVIRONMENT", "development")


	cfg := &Config{
		Environment: env,


		Server: ServerConfig{
			Port: getPortForEnvironment(env),
			ReadTimeout: getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout: getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
			TLSEnabled: env == "production",
			TLSCertFile: getEnv("TLS_CERT_FILE", ""),
			TLSKeyFile: getEnv("TLS_KEY_FILE", ""),
		},

		Database: DatabaseConfig{
			Host: getEnv("DB_HOST", "localhost"),
			Port: getEnv("DB_PORT", "5432"),
			User: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName: getEnv("DB_NAME", "eyesky"),
			SSLMode: getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 20),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 510),
			MaxLifetime: getEnvAsDuration("DB_MAX_LIFETIME", 5*time.Minute),
		},
		OpenSky: OpenSkyConfig{
			BaseURL: "https://opensky-network.org/api",
			Username: getEnv("OPEN_SKY_USERNAME", ""),
			Password: getEnv("OPEN_SKY_PASSWORD", ""),
			Timeout: getEnvAsDuration("OPEN_SKY_TIMEOUT", 30*time.Second),
			RateLimit: getEnvAsInt("OPEN_SKY_RATE_LIMIT", 100),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func getPortForEnvironment(env string) string {
	if env == "production" {
		return getEnv("PORT", "443")
	}
	return getEnv("PORT", "8080")
}


