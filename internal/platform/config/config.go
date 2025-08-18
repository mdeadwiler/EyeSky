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