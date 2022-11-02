package arangodb

import (
	"context"
	"errors"
	"github.com/arangodb/go-driver"
	"github.com/serdarkalayci/docman/adapters/data/arangodb/dao"
)

type arangoHelper struct {
	db driver.Database
}

func (ah arangoHelper) Find(ctx context.Context) ([]dao.DocumentDAO, error) {
	var documentDAOs = make([]dao.DocumentDAO, 0)
	return documentDAOs, errors.New("not implemented")
}

func (ah arangoHelper) InsertOne(ctx context.Context, document interface{}) (string, error) {
	return "", errors.New("not implemented")
}
func (ah arangoHelper) FindOne(ctx context.Context, id string) (dao.DocumentDAO, error) {
	var documentDAO dao.DocumentDAO
	return documentDAO, errors.New("not implemented")
}

func (ah arangoHelper) UpdateOne(ctx context.Context, id string, update interface{}) (int, error) {
	return 0, errors.New("not implemented")
}

func (ah arangoHelper) DeleteOne(ctx context.Context, id string) (int, error) {
	return 0, errors.New("not implemented")
}
