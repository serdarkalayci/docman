package arangodb

import (
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/nicholasjackson/env"
)

var databaseName = env.String("DatabaseName", false, "titanic", "The database name for arangodb")
var connectionString = env.String("ConnectionString", false, "https://127.0.0.1:8529", "Database connection string")
var username = env.String("DbUserName", false, "root", "Database username")
var password = env.String("DbPassword", false, "password", "Database password")

// DataContext represents a struct that holds concrete repositories
type DataContext struct {
	DocumentRepository DocumentRepository
	HealthRepository   HealthRepository
}

// NewDataContext returns a new mongoDB backed DataContext
func NewDataContext() (DataContext, error) {

	env.Parse()
	// We try to get connectionstring value from the environment variables, if not found it falls back to local database

	// Open a client connection
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{*connectionString},
	})
	if err != nil {
		// Handle error
	}

	// Client object
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(*username, *password),
	})
	if err != nil {
		// Handle error
	}

	// Open "examples_books" database
	db, err := client.Database(nil, *databaseName)
	if err != nil {
		// Handle error
	}
	dataContext := DataContext{}
	dataContext.DocumentRepository = newDocumentRepository(db)
	dataContext.HealthRepository = newHealthRepository(db)
	return dataContext, nil
}
