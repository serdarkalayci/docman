package arangodb

import (
	"context"
	"fmt"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/nicholasjackson/env"
	"log"
	"time"
)

var databaseName = env.String("DatabaseName", false, "docman", "The database name for arangodb")
var connectionString = env.String("ConnectionString", false, "http://localhost:8529", "Database connection string")
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
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
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
	users, err := client.Users(ctx)
	fmt.Println(users)
	db_exists, err := client.DatabaseExists(ctx, *databaseName)

	if db_exists {
		fmt.Println("That db exists already")
	} else {
		_, err = client.CreateDatabase(ctx, *databaseName, nil)

		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
	}
	// Open a database
	db, err := client.Database(ctx, *databaseName)
	if err != nil {
		// Handle error
	}
	dataContext := DataContext{}
	dataContext.DocumentRepository = newDocumentRepository(db)
	dataContext.HealthRepository = newHealthRepository(db)
	return dataContext, nil
}
