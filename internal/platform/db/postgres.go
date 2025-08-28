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

type DB struct {
	conn *sql.DB
	config config.DatabaseConfig
}

func New(cfg config.DatabaseConfig) (*DB, error) {
	// DB connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	
		conn, err := sql.Open("postgress", connStr)
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
    // defer ctx.Done()

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
	// Migration runner

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