package arangodb

import (
	"github.com/arangodb/go-driver"
)

// HealthRepository represent a structure that will communicate to MongoDB to accomplish health related transactions
type HealthRepository struct {
	helper dbHelper
}

func newHealthRepository(database driver.Database) HealthRepository {
	return HealthRepository{
		helper: arangoHelper{db: database},
	}
}

// Ready checks the arangodb connection
func (hr HealthRepository) Ready() bool {
	return false
}
