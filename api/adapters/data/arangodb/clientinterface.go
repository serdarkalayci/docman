package arangodb

import (
	"context"

	"github.com/serdarkalayci/docman/api/adapters/data/arangodb/dao"
)

type dbHelper interface {
	Find(ctx context.Context, id string) (dao.FolderTreeDAO, error)
	InsertOne(ctx context.Context, document interface{}) (string, error)
	FindOne(ctx context.Context, id string) (dao.DocumentDAO, error)
	UpdateOne(ctx context.Context, id string, update interface{}) (int, error)
	DeleteOne(ctx context.Context, id string) (int, error)
}
