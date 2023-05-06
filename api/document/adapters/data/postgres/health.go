package postgres

import (
	"github.com/jackc/pgx/v5"
)

// HealthRepository represent a structure that will communicate to MongoDB to accomplish health related transactions
type HealthRepository struct {
	db *pgx.Conn
}

func newHealthRepository(database *pgx.Conn) HealthRepository {
	return HealthRepository{
		db: database,
	}
}

// Ready checks the arangodb connection
func (hr HealthRepository) Ready() bool {
	return false
}
