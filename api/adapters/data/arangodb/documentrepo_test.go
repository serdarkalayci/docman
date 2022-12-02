package arangodb

import (
	"context"
	"fmt"
	"testing"

	driver "github.com/arangodb/go-driver"
	"github.com/serdarkalayci/docman/api/adapters/data/arangodb/dao"
	"github.com/serdarkalayci/docman/api/domain"
	"github.com/stretchr/testify/assert"
)

// #region MockArangoHelper

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

// #endregion MockArangoHelper

// #region MockCursor
type MockCursor struct {
}

var (
	readDocumentFunc func(ctx context.Context, result interface{}) (driver.DocumentMeta, error)
	countFunc        func() int64
)

func (mc MockCursor) ReadDocument(ctx context.Context, result interface{}) (driver.DocumentMeta, error) {
	return readDocumentFunc(ctx, result)
}

func (mc MockCursor) Close() error {
	return nil
}

func (mc MockCursor) Count() int64 {
	return countFunc()
}

func (mc MockCursor) HasMore() bool {
	return true
}

func (mc MockCursor) Statistics() driver.QueryStatistics {
	return nil
}

func (mc MockCursor) Extra() driver.QueryExtra {
	return nil
}

//#endregion MockCursor

type MockDocument struct {
	ID      string
	Key     string
	Name    string
	Content string
	Meta    driver.DocumentMeta
	Error   error
}

var (
	mockDocuments = []MockDocument{
		MockDocument{
			ID:    "111",
			Key:   "folders/111",
			Name:  "Folder 111",
			Meta:  driver.DocumentMeta{ID: driver.NewDocumentID("folders", "111"), Key: "folders/111"},
			Error: nil,
		},
		MockDocument{
			ID:      "222",
			Key:     "documents/222",
			Name:    "File 222",
			Content: "File 222 Content",
			Meta:    driver.DocumentMeta{ID: driver.NewDocumentID("documents", "111"), Key: "documents/111"},
			Error:   nil,
		},
		MockDocument{
			Error: driver.NoMoreDocumentsError{},
		},
	}
	mockParent = []MockDocument{
		MockDocument{
			ID:    "111",
			Key:   "folders/111",
			Name:  "Folder 111",
			Meta:  driver.DocumentMeta{ID: driver.NewDocumentID("folders", "111"), Key: "folders/111"},
			Error: nil,
		},
		MockDocument{
			Error: driver.NoMoreDocumentsError{},
		},
	}
	docIndex    = 0
	parentIndex = 0
)

func TestDocumentRepository_List(t *testing.T) {
	docIndex = 0
	dr := DocumentRepository{
		helper: MockArangoHelper{},
	}
	// Lets test what happens when the function cannot find the main item
	findItemFunc = func(ctx context.Context, id string, collection string, item interface{}) error {
		return fmt.Errorf("cannot find item")
	}
	folder, err := dr.List("123")
	assert.Equal(t, domain.Folder{}, folder)
	assert.Equal(t, fmt.Errorf("cannot find item"), err)
	// Lets fix findItem function and test what happens when the function cannot find the children
	findItemFunc = func(ctx context.Context, id string, collection string, item interface{}) error {
		item.(*dao.FolderDAO).ID = "123"
		item.(*dao.FolderDAO).Name = "test"
		item.(*dao.FolderDAO).ParentFolderID = "1"
		return nil
	}
	findChildrenFunc = func(ctx context.Context, id string, graphName string) (driver.Cursor, error) {
		return nil, fmt.Errorf("cannot find children")
	}
	folder, err = dr.List("123")
	assert.Equal(t, domain.Folder{}, folder)
	assert.Equal(t, fmt.Errorf("cannot find children"), err)
	// Lets fix findChildren function and test what happens when the readDocument returns an error
	findChildrenFunc = func(ctx context.Context, id string, graphName string) (driver.Cursor, error) {
		return MockCursor{}, nil
	}
	readDocumentFunc = func(ctx context.Context, result interface{}) (driver.DocumentMeta, error) {
		return driver.DocumentMeta{}, fmt.Errorf("cannot read document")
	}
	folder, err = dr.List("123")
	assert.Equal(t, domain.Folder{}, folder)
	assert.Equal(t, fmt.Errorf("cannot read document"), err)
	// Lets fix readDocument function and test what happens when the findParent returns an error
	countFunc = func() int64 {
		return int64(len(mockDocuments))
	}
	readDocumentFunc = func(ctx context.Context, result interface{}) (driver.DocumentMeta, error) {
		result.(*dao.DocumentDAO).ID = mockDocuments[docIndex].ID
		result.(*dao.DocumentDAO).Key = mockDocuments[docIndex].Key
		result.(*dao.DocumentDAO).Name = mockDocuments[docIndex].Name
		result.(*dao.DocumentDAO).Content = mockDocuments[docIndex].Content
		meta := mockDocuments[docIndex].Meta
		err := mockDocuments[docIndex].Error
		docIndex++
		return meta, err
	}
	findParentFunc = func(ctx context.Context, id string, graphName string) (driver.Cursor, error) {
		return nil, fmt.Errorf("cannot find parent")
	}
	folder, err = dr.List("123")
	assert.Equal(t, domain.Folder{}, folder)
	assert.Equal(t, fmt.Errorf("cannot find parent"), err)
	// Lets fix findParent function and test what happens when the readDocument returns an error
	findParentFunc = func(ctx context.Context, id string, graphName string) (driver.Cursor, error) {
		return MockCursor{}, nil
	}
	readDocumentFunc = func(ctx context.Context, result interface{}) (driver.DocumentMeta, error) {
		return driver.DocumentMeta{}, fmt.Errorf("cannot read parent")
	}
	folder, err = dr.List("123")
	assert.Equal(t, domain.Folder{}, folder)
	assert.Equal(t, fmt.Errorf("cannot read parent"), err)
	// Lets fix readDocument function and test what happens when the function works
	docIndex = 0
	parentIndex = 0
	countFunc = func() int64 {
		return int64(len(mockParent))
	}
	readDocumentFunc = func(ctx context.Context, result interface{}) (driver.DocumentMeta, error) {
		if docIndex < len(mockDocuments) {
			result.(*dao.DocumentDAO).ID = mockDocuments[docIndex].ID
			result.(*dao.DocumentDAO).Key = mockDocuments[docIndex].Key
			result.(*dao.DocumentDAO).Name = mockDocuments[docIndex].Name
			result.(*dao.DocumentDAO).Content = mockDocuments[docIndex].Content
			meta := mockDocuments[docIndex].Meta
			err := mockDocuments[docIndex].Error
			docIndex++
			return meta, err
		} else {
			result.(*dao.DocumentDAO).ID = mockParent[parentIndex].ID
			result.(*dao.DocumentDAO).Key = mockParent[parentIndex].Key
			result.(*dao.DocumentDAO).Name = mockParent[parentIndex].Name
			meta := mockParent[parentIndex].Meta
			err := mockParent[parentIndex].Error
			parentIndex++
			return meta, err
		}
	}
	folder, err = dr.List("123")
	assert.Equal(t, domain.Folder{ID: "123", Name: "test", ParentFolderID: "111", Folders: []domain.Folder{{ID: "111", Name: "Folder 111", ParentFolderID: "", Folders: []domain.Folder(nil), Documents: []domain.Document(nil)}}, Documents: []domain.Document{{ID: "222", Name: "File 222", Content: "File 222 Content", DocumentHistory: []domain.History(nil)}}}, folder)
	assert.Nil(t, err)
}
