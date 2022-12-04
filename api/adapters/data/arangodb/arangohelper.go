package arangodb

import (
	"context"
	"errors"
	"fmt"

	driver "github.com/arangodb/go-driver"
)

type arangoHelper struct {
	db driver.Database
}

func (ah arangoHelper) beginTransaction(ctx context.Context, cols []string) (driver.TransactionID, context.Context, error) {
	tranID, err := ah.db.BeginTransaction(ctx, driver.TransactionCollections{Write: cols}, &driver.BeginTransactionOptions{})
	if err != nil {
		return "0", nil, err
	}
	tranctx := driver.WithTransactionID(ctx, tranID)
	return tranID, tranctx, nil
}

func (ah arangoHelper) commitTransaction(ctx context.Context, id driver.TransactionID) error {
	return ah.db.CommitTransaction(ctx, id, &driver.CommitTransactionOptions{})
}

func (ah arangoHelper) abortTransaction(ctx context.Context, id driver.TransactionID) error {
	return ah.db.AbortTransaction(ctx, id, &driver.AbortTransactionOptions{})
}

func (ah arangoHelper) findItem(ctx context.Context, id string, collection string, item interface{}) error {
	// Open specified collection
	col, err := ah.db.Collection(nil, collection)
	if err != nil {
		return err
	}
	// Try to find the item in the corresponding collection with the given id
	_, err = col.ReadDocument(nil, id, item)
	if err != nil {
		return err
	}
	return nil
}

func (ah arangoHelper) findChildren(ctx context.Context, id string, graphName string) (driver.Cursor, error) {
	querystring := fmt.Sprintf("FOR v IN 1..1 OUTBOUND @currentItem GRAPH '%s' RETURN v", graphName)
	bindVars := map[string]interface{}{
		"currentItem": id,
	}
	cursor, err := ah.db.Query(ctx, querystring, bindVars)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (ah arangoHelper) findParent(ctx context.Context, id string, graphName string) (driver.Cursor, error) {
	// Open "folders" collection
	querystring := fmt.Sprintf("FOR v IN 1..1 INBOUND @currentItem GRAPH '%s' RETURN v", graphName)
	bindVars := map[string]interface{}{
		"currentItem": id,
	}
	cursor, err := ah.db.Query(ctx, querystring, bindVars)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (ah arangoHelper) insertItem(ctx context.Context, document interface{}, collection string) (string, error) {
	// Open specified collection
	col, err := ah.db.Collection(ctx, collection)
	if err != nil {
		return "", err
	}
	meta, err := col.CreateDocument(ctx, document)
	if err != nil {
		return "", err
	}
	return meta.Key, nil
}

func (ah arangoHelper) insertEdge(ctx context.Context, fromID string, toID string, collection string) (string, error) {
	// Open specified collection
	col, err := ah.db.Collection(nil, collection)
	if err != nil {
		return "", err
	}
	edge := driver.EdgeDocument{
		From: driver.DocumentID(fromID),
		To:   driver.DocumentID(toID),
	}
	meta, err := col.CreateDocument(ctx, edge)
	if err != nil {
		return "", err
	}
	return meta.Key, nil
}

func (ah arangoHelper) updateItem(ctx context.Context, id string, collection string, update interface{}) (int, error) {
	return 0, errors.New("not implemented")
}

func (ah arangoHelper) deleteItem(ctx context.Context, collection string, id string) (int, error) {
	return 0, errors.New("not implemented")
}
