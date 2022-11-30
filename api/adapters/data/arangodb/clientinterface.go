package arangodb

import (
	"context"

	driver "github.com/arangodb/go-driver"
)

type dbHelper interface {
	beginTransaction(ctx context.Context, cols driver.TransactionCollections) (driver.TransactionID, error)
	commitTransaction(ctx context.Context, id driver.TransactionID) error
	abortTransaction(ctx context.Context, id driver.TransactionID) error
	findItem(ctx context.Context, id string, collection string, item interface{}) error
	findParent(ctx context.Context, id string, graphName string) (driver.Cursor, error)
	findChildren(ctx context.Context, id string, graphName string) (driver.Cursor, error)
	insertItem(ctx context.Context, document interface{}, collection string) (string, error)
	insertEdge(ctx context.Context, from string, to string, graphName string) (string, error)
	updateItem(ctx context.Context, id string, collection string, update interface{}) (int, error)
	deleteItem(ctx context.Context, collection string, id string) (int, error)
}
