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

type db struct {
	conn *sql.DB
	config config.DatabaseConfig
}

func New(cfg config.DatabaseConfig) (*DB, error) {
	// DB connection

}

func (db *DB) Migrate() error {
	// Migration runner
}

// For DB connection repositories
func (db *db) GetConnection() *sql.DB {
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