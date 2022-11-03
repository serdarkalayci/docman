package arangodb

import (
	"context"
	"errors"
	"testing"

	"github.com/serdarkalayci/docman/adapters/data/arangodb/dao"
	"github.com/serdarkalayci/docman/domain"
	"github.com/stretchr/testify/assert"
)

// MockArangoHelper is the struct that mimics original arangodb.Client
type MockArangoHelper struct {
}

var (
	// GetDeleteFunc will be used to get different Delete functions for testing purposes
	GetDeleteFunc func(ctx context.Context, id string) (int, error)
	// GetDeleteFunc will be used to get different Update functions for testing purposes
	GetUpdateFunc func(ctx context.Context, id string, update interface{}) (int, error)
	// GetFindOneFunc will be used to get different FindOne functions for testing purposes
	GetFindOneFunc func(ctx context.Context, id string) (dao.DocumentDAO, error)
	// GetInsertOneFunc will be used to get different InsertOne functions for testing purposes
	GetInsertOneFunc func(ctx context.Context, document interface{}) (string, error)
	// GetListFunc will be used to get different List functions for testing purposes
	GetListFunc func(ctx context.Context, id string) (dao.FolderTreeDAO, error)
)

// func (client MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
// 	return GetDoFunc(req)
// }

func (ah MockArangoHelper) Find(ctx context.Context, id string) (dao.FolderTreeDAO, error) {
	return GetListFunc(ctx, id)
}
func (ah MockArangoHelper) InsertOne(ctx context.Context, document interface{}) (string, error) {
	return GetInsertOneFunc(ctx, document)
}
func (ah MockArangoHelper) FindOne(ctx context.Context, id string) (dao.DocumentDAO, error) {
	return GetFindOneFunc(ctx, id)
}
func (ah MockArangoHelper) UpdateOne(ctx context.Context, id string, update interface{}) (int, error) {
	return GetUpdateFunc(ctx, id, update)
}
func (ah MockArangoHelper) DeleteOne(ctx context.Context, id string) (int, error) {
	return GetDeleteFunc(ctx, id)
}

func TestDocumentRepository_Delete_Error(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	GetDeleteFunc = func(ctx context.Context, id string) (int, error) {
		return 0, errors.New("Whatever error")
	}
	err := dr.Delete("id")
	assert.EqualError(t, err, "Error deleting the document")
}

func TestDocumentRepository_Delete_ResultNotOne(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	GetDeleteFunc = func(ctx context.Context, id string) (int, error) {
		return 0, nil
	}
	err := dr.Delete("this_id")
	assert.EqualError(t, err, "Cannot find the document with the ID this_id")
}

func TestDocumentRepository_Delete_ResultSuccess(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	GetDeleteFunc = func(ctx context.Context, id string) (int, error) {
		return 1, nil
	}
	err := dr.Delete("this_id")
	assert.Nil(t, err)
}

func TestDocumentRepository_Update_Error(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	GetUpdateFunc = func(ctx context.Context, id string, update interface{}) (int, error) {
		return 0, errors.New("Whatever error")
	}
	err := dr.Update("id", domain.Document{})
	assert.EqualError(t, err, "Error updating the document")
}

func TestDocumentRepository_Update_ResultNotOne(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	GetUpdateFunc = func(ctx context.Context, id string, update interface{}) (int, error) {
		return 0, nil
	}
	err := dr.Update("this_id", domain.Document{})
	assert.EqualError(t, err, "Cannot find the document with the ID this_id")
}

func TestDocumentRepository_Update_ResultSuccess(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	GetUpdateFunc = func(ctx context.Context, id string, update interface{}) (int, error) {
		return 1, nil
	}
	err := dr.Update("id", domain.Document{})
	assert.Nil(t, err)
}

func TestDocumentRepository_FindOne_Error(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	GetFindOneFunc = func(ctx context.Context, id string) (dao.DocumentDAO, error) {
		return dao.DocumentDAO{}, errors.New("Cannot find the document with the ID this_id")
	}
	pDAO, err := dr.Get("this_id")
	assert.Equal(t, pDAO, domain.Document{})
	assert.EqualError(t, err, "Cannot find the document with the ID this_id")
}

func TestDocumentRepository_FindOne_Success(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	GetFindOneFunc = func(ctx context.Context, id string) (dao.DocumentDAO, error) {
		return dao.DocumentDAO{
			ID:   "id",
			Name: "name",
		}, nil
	}
	pDAO, err := dr.Get("this_id")
	assert.Equal(t, pDAO, domain.Document{
		ID:   "id",
		Name: "name",
	})
	assert.Nil(t, err)
}

func TestDocumentRepository_InsertOne_Error(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	GetInsertOneFunc = func(ctx context.Context, document interface{}) (string, error) {
		return "", errors.New("Whatever error")
	}
	document, err := dr.Add(domain.Document{
		ID: "this_id",
	})
	assert.Equal(t, document, domain.Document{})
	assert.EqualError(t, err, "Error adding document")
}

func TestDocumentRepository_InsertOne_Success(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	GetInsertOneFunc = func(ctx context.Context, document interface{}) (string, error) {
		return "new_id", nil
	}
	document, err := dr.Add(domain.Document{
		Name: "name",
	})
	assert.Equal(t, document, domain.Document{
		ID:   "new_id",
		Name: "name",
	})
	assert.Nil(t, err)
}

func TestDocumentRepository_List_Error(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	GetListFunc = func(ctx context.Context, id string) (dao.FolderTreeDAO, error) {
		return dao.FolderTreeDAO{}, errors.New("Whatever error")
	}
	result, err := dr.List("")
	assert.Equal(t, domain.Folder{}, result)
	assert.EqualError(t, err, "Error getting documents")
}

func TestDocumentRepository_List_Success(t *testing.T) {
	dr := DocumentRepository{MockArangoHelper{}}
	ftDAO := dao.FolderTreeDAO{
		CurrentFolder: dao.FolderDAO{
			ID:   "folders/0",
			Name: "current",
		},
		SubFolders: []dao.FolderDAO{
			{
				ID:   "folders/1",
				Name: "sub1",
			},
			{
				ID:   "folders/2",
				Name: "sub2",
			},
		},
		Documents: []dao.DocumentDAO{
			{
				ID:   "documents/1",
				Name: "name1",
			},
			{
				ID:   "documents/2",
				Name: "name2",
			},
		},
	}
	folder := domain.Folder{
		ID:   "folders/0",
		Name: "current",
		Folders: []domain.Folder{
			{
				ID:   "folders/1",
				Name: "sub1",
			},
			{
				ID:   "folders/2",
				Name: "sub2",
			},
		},
		Documents: []domain.Document{
			{
				ID:   "documents/1",
				Name: "name1",
			},
			{
				ID:   "documents/2",
				Name: "name2",
			},
		},
	}
	GetListFunc = func(ctx context.Context, id string) (dao.FolderTreeDAO, error) {
		return ftDAO, nil
	}
	result, err := dr.List("")
	assert.Nil(t, err)
	assert.Equal(t, result, folder)
}
