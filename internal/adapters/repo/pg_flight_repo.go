package repo

import (
	"context"
	"database/sql"
	"time"
	"fmt"
	"github.com/mdeadwiler/EyeSky/internal/domain"
	"github.com/mdeadwiler/EyeSky/internal/platform/db"	
)


type Repository struct {
	db *sql.DB
}

// New SQL flight repo
func New(database *db.DB) *Repository {
	return &Repository{
		db: database.GetConnection(),
	}
}