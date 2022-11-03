package arangodb

import (
	"context"
	"errors"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/serdarkalayci/docman/adapters/data/arangodb/dao"
)

type arangoHelper struct {
	db driver.Database
}

func (ah arangoHelper) Find(ctx context.Context, id string) (dao.FolderTreeDAO, error) {
	var folderTree dao.FolderTreeDAO
	var currentFolder dao.FolderDAO
	// Open "folders" collection
	col, err := ah.db.Collection(nil, "folders")
	if err != nil {
		return dao.FolderTreeDAO{}, err
	}
	_, err = col.ReadDocument(nil, id, &currentFolder)
	if err != nil {
		return dao.FolderTreeDAO{}, err
	}
	folderTree.CurrentFolder = currentFolder
	querystring := "FOR v, e, p IN 1..1 OUTBOUND @currentFolder GRAPH 'filesystem' RETURN v"
	bindVars := map[string]interface{}{
		"currentFolder": fmt.Sprintf("folders/%s", id),
	}
	cursor, err := ah.db.Query(ctx, querystring, bindVars)
	if err != nil {
		return dao.FolderTreeDAO{}, err
	}
	defer cursor.Close()
	for {
		var doc dao.DocumentDAO
		meta, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return dao.FolderTreeDAO{}, err
		}
		if meta.ID.Collection() == "folders" {
			folder := dao.FolderDAO{
				ID:   doc.ID,
				Key:  doc.Key,
				Name: doc.Name,
			}
			folderTree.SubFolders = append(folderTree.SubFolders, folder)
		} else {
			folderTree.Documents = append(folderTree.Documents, doc)
		}
	}
	return folderTree, nil
}

func (ah arangoHelper) InsertOne(ctx context.Context, document interface{}) (string, error) {
	return "", errors.New("not implemented")
}
func (ah arangoHelper) FindOne(ctx context.Context, id string) (dao.DocumentDAO, error) {
	var documentDAO dao.DocumentDAO
	// Open "documents" collection
	col, err := ah.db.Collection(nil, "documents")
	if err != nil {
		return documentDAO, err
	}
	_, err = col.ReadDocument(nil, id, &documentDAO)
	if err != nil {
		return documentDAO, err
	}
	return documentDAO, nil
}

func (ah arangoHelper) UpdateOne(ctx context.Context, id string, update interface{}) (int, error) {
	return 0, errors.New("not implemented")
}

func (ah arangoHelper) DeleteOne(ctx context.Context, id string) (int, error) {
	return 0, errors.New("not implemented")
}
