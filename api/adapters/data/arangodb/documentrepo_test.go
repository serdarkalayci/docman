package arangodb

import (
	"context"

	driver "github.com/arangodb/go-driver"
)

// MockArangoHelper is the struct that mimics original arangodb.Client
type MockArangoHelper struct {
}

var (
	beginTransactionFunc  func(ctx context.Context, cols driver.TransactionCollections) (driver.TransactionID, error)
	commitTransactionFunc func(ctx context.Context, id driver.TransactionID) error
	abortTransactionFunc  func(ctx context.Context, id driver.TransactionID) error
	findItemFunc          func(ctx context.Context, id string, collection string, item interface{}) error
	findParentFunc        func(ctx context.Context, id string, graphName string) (driver.Cursor, error)
	findChildrenFunc      func(ctx context.Context, id string, graphName string) (driver.Cursor, error)
	insertItemFunc        func(ctx context.Context, document interface{}, collection string) (string, error)
	insertEdgeFunc        func(ctx context.Context, from string, to string, graphName string) (string, error)
	updateItemFunc        func(ctx context.Context, id string, collection string, update interface{}) (int, error)
	deleteItemFunc        func(ctx context.Context, collection string, id string) (int, error)
)

func (ah MockArangoHelper) beginTransaction(ctx context.Context, cols driver.TransactionCollections) (driver.TransactionID, error) {
	return beginTransactionFunc(ctx, cols)
}

func (ah MockArangoHelper) commitTransaction(ctx context.Context, id driver.TransactionID) error {
	return commitTransactionFunc(ctx, id)
}
func (ah MockArangoHelper) abortTransaction(ctx context.Context, id driver.TransactionID) error {
	return abortTransactionFunc(ctx, id)
}
func (ah MockArangoHelper) findItem(ctx context.Context, id string, collection string, item interface{}) error {
	return findItemFunc(ctx, id, collection, item)
}
func (ah MockArangoHelper) findParent(ctx context.Context, id string, graphName string) (driver.Cursor, error) {
	return findParentFunc(ctx, id, graphName)
}
func (ah MockArangoHelper) findChildren(ctx context.Context, id string, graphName string) (driver.Cursor, error) {
	return findChildrenFunc(ctx, id, graphName)
}
func (ah MockArangoHelper) insertItem(ctx context.Context, document interface{}, collection string) (string, error) {
	return insertItemFunc(ctx, document, collection)
}
func (ah MockArangoHelper) insertEdge(ctx context.Context, from string, to string, graphName string) (string, error) {
	return insertEdgeFunc(ctx, from, to, graphName)
}
func (ah MockArangoHelper) updateItem(ctx context.Context, id string, collection string, update interface{}) (int, error) {
	return updateItemFunc(ctx, id, collection, update)
}
func (ah MockArangoHelper) deleteItem(ctx context.Context, collection string, id string) (int, error) {
	return deleteItemFunc(ctx, collection, id)
}
