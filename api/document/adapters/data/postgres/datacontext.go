package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nicholasjackson/env"
)

var databaseName = env.String("DatabaseName", false, "docman", "The database name for arangodb")
var connectionString = env.String("ConnectionString", false, "localhost:26257", "Database connection string")
var username = env.String("DbUserName", false, "docmanuser", "Database username")
var password = env.String("DbPassword", false, "docmanpassword", "Database password")

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
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", *username, *password, *connectionString, *databaseName)
	fmt.Printf("Connecting to %s\n", dsn)
	// Client object
	db, err := pgx.Connect(ctx, dsn)
    if err != nil {
		log.Fatalln(err)
        return DataContext{}, fmt.Errorf("cannot create the database client on %s", *connectionString)
    }
	if err != nil {
		return DataContext{}, fmt.Errorf("cannot create the database client on %s", *connectionString)
	}
	// Open a database. In case the database is not ready yet, we retry a few times
	count := 0
	for count < 5 {
		if err := db.Ping(ctx); err == nil {
			break
		}
		count++
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return DataContext{}, fmt.Errorf("cannot connect to the database on %s", *connectionString)
	}
	dataContext := DataContext{}
	dataContext.DocumentRepository = newDocumentRepository(db)
	dataContext.HealthRepository = newHealthRepository(db)
	return dataContext, nil
}
