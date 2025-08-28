package db
// Postgres connection and migration runner 
import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/mdeadwiler/EyeSky/internal/platform/config"
)


const (
	defaultConnectionTimeout = 5 * time.Second
	migrationPath = "internal/platform/db/migrations"
)

type DB struct {
	conn *sql.DB
	config config.DatabaseConfig
}

func New(cfg config.DatabaseConfig) (*DB, error) {
	// DB connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
		}

	// Configure connection pool for flight data
	conn.SetMaxOpenConns(cfg.MaxOpenConns)
	conn.SetMaxIdleConns(cfg.MaxIdleConns)
	conn.SetConnMaxLifetime(cfg.MaxLifetime)

	// Test Connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{
		conn: conn,
		config: cfg,
	}
	return db, nil	
}

func (db *DB) Migrate() error {
	// Migration path
	migrationDir, err := filepath.Abs(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to get migration directory: %w", err)
	}

	// File source for migrations
	sourceURL := fmt.Sprintf("file://%s", migrationDir)
	source, err := (&file.File{}).Open(sourceURL)
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}
	defer source.Close()

	// DB driver
	driver, err := postgres.WithInstance(db.conn, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// Pass down migration instance
	m, err := migrate.NewWithInstance("file", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	// Run latest migration
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			// No migrations
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

// For DB connection repositories
func (db *DB) GetConnection() *sql.DB {
	return db.conn
}

// Close DB connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Pings DB for health check
func (db *DB) Ping(ctx context.Context) error {
	return db.conn.PingContext(ctx)
}
