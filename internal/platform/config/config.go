package config

import (
	"os"
	"strconv"
	"time"
)

export type Config struct {
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

func New(*Config, error) {
	env := getEnv("ENVIRONMENT", "development")

	cfg := &Cnfig{
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
			Host: getEnv("DB_Host", "localhost"),
			Port: getEnv("DB_Port", "5432"),
			User: getEnv("DB_User", "postgres"),
			Password: getEnv("DB_Password", ""),
			DBName: getEnv("DB_Name", "eyesky"),
			SSLMode: getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 20),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 510),
			MaxLifetime: getEnvAsDuration("DB_MAX_LIFETIME", 5*time.Minute),
		},
		

	}
}
