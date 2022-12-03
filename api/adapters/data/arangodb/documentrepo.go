package arangodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	driver "github.com/arangodb/go-driver"

	"github.com/rs/zerolog/log"
	"github.com/serdarkalayci/docman/api/adapters/data/arangodb/dao"
	"github.com/serdarkalayci/docman/api/adapters/data/arangodb/mappers"
	"github.com/serdarkalayci/docman/api/application"
	"github.com/serdarkalayci/docman/api/domain"
)

// DocumentRepository holds the arangodb client and database name for methods to use
type DocumentRepository struct {
	helper dbHelper
}

func newDocumentRepository(database driver.Database) DocumentRepository {
	return DocumentRepository{
		helper: arangoHelper{db: database},
	}
}

// List loads all the document records from tha database and returns it
// Returns an error if database fails to provide service
func (dr DocumentRepository) List(id string) (domain.Folder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// A list request should return the item itself, its first level children and its parent (if any)
	var folderTree dao.FolderTreeDAO // This item will be our main carrier of data to be returned
	var currentFolder dao.FolderDAO  // This item will be used to hold the current folder
	// Lets first get the item (in this case the folder) itself
	err := dr.helper.findItem(ctx, id, "folders", &currentFolder)
	if err != nil {
		log.Error().Err(err).Msgf("error finding folder")
		return domain.Folder{}, err
	}
	folderTree.CurrentFolder = currentFolder
	// Lets get the children of the item and
	cursor, err := dr.helper.findChildren(ctx, fmt.Sprintf("folders/%s", id), "filesystem")
	if err != nil {
		log.Error().Err(err).Msgf("error finding children")
		return domain.Folder{}, err
	}
	defer cursor.Close()
	for {
		var doc dao.DocumentDAO
		meta, err := cursor.ReadDocument(ctx, &doc)
		if err != nil && driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return domain.Folder{}, err
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
	var parentFolder dao.FolderDAO
	cursor, err = dr.helper.findParent(ctx, folderTree.CurrentFolder.ID, "filesystem")
	if err != nil {
		return domain.Folder{}, err
	}
	defer cursor.Close()
	for {
		var doc dao.DocumentDAO
		meta, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return domain.Folder{}, err
		}
		if meta.ID.Collection() == "folders" {
			parentFolder = dao.FolderDAO{
				ID:   doc.ID,
				Key:  doc.Key,
				Name: doc.Name,
			}

		}
	}
	folderTree.CurrentFolder.ParentFolderID = parentFolder.ID

	folder := mappers.MapFolderTreeDAO2Folder(folderTree)
	return folder, nil
}

// Get selects a single document from the database with the given unique identifier
// Returns an error if database fails to provide service
func (dr DocumentRepository) Get(id string) (domain.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	documentDAO := dao.DocumentDAO{}
	err := dr.helper.findItem(ctx, id, "documents", &documentDAO)
	if err != nil {
		log.Error().Err(err).Msgf("error getting document")
		return domain.Document{}, &application.ErrorCannotFinddocument{ID: id}
	}
	return mappers.MapDocumentDAO2Document(documentDAO), nil
}

// AddDocument adds a new document to the underlying database.
// It returns the document inserted on success or error
func (dr DocumentRepository) AddDocument(p domain.Document, parentID string) (domain.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	doc := mappers.MapDocument2DocumentDAO(p)
	cols := driver.TransactionCollections{
		Write: []string{"documents", "filesystem"},
	}
	tranID, err := dr.helper.beginTransaction(ctx, cols)
	tranctx := driver.WithTransactionID(ctx, tranID)
	if err != nil {
		log.Error().Err(err).Msgf("Error starting transaction while adding document")
		return domain.Document{}, errors.New("Error adding folder")
	}
	newID, err := dr.helper.insertItem(tranctx, doc, "documents")
	if err != nil {
		dr.helper.abortTransaction(ctx, tranID)
		log.Error().Err(err).Msgf("Error adding document")
		return domain.Document{}, errors.New("Error adding document")
	}
	p.ID = newID
	// Now lets add the document to the filesystem
	_, err = dr.helper.insertEdge(tranctx, fmt.Sprintf("folders/%s", parentID), fmt.Sprintf("documents/%s", newID), "filesystem")
	if err != nil {
		tranerr := dr.helper.abortTransaction(ctx, tranID)
		if tranerr != nil {
			log.Error().Err(tranerr).Msgf("Error aborting transaction while adding document")
		}
		log.Error().Err(err).Msgf("Error adding document to filesystem")
		return domain.Document{}, errors.New("Error adding document to filesystem")
	}
	dr.helper.commitTransaction(ctx, tranID)
	return p, nil
}

// AddFolder adds a new folder to the underlying database.
// It returns the folder inserted on success or error
func (dr DocumentRepository) AddFolder(p domain.Folder, parentID string) (domain.Folder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	folder := mappers.MapFolder2FolderDAO(p)
	cols := driver.TransactionCollections{
		Write: []string{"folders", "filesystem"},
	}
	tranID, err := dr.helper.beginTransaction(ctx, cols)
	if err != nil {
		log.Error().Err(err).Msgf("Error starting transaction while adding folder")
		return domain.Folder{}, errors.New("Error adding folder")
	}
	tranctx := driver.WithTransactionID(ctx, tranID)
	newID, err := dr.helper.insertItem(tranctx, folder, "folders")
	if err != nil {
		dr.helper.abortTransaction(ctx, tranID)
		log.Error().Err(err).Msgf("Error adding folder")
		return domain.Folder{}, errors.New("Error adding folder")
	}
	p.ID = newID
	// Now lets add the document to the filesystem
	_, err = dr.helper.insertEdge(tranctx, fmt.Sprintf("folders/%s", parentID), fmt.Sprintf("documents/%s", newID), "filesystem")
	if err != nil {
		tranerr := dr.helper.abortTransaction(ctx, tranID)
		if tranerr != nil {
			log.Error().Err(tranerr).Msgf("Error aborting transaction while adding folder")
		}
		log.Error().Err(err).Msgf("Error adding folder to filesystem")
		return domain.Folder{}, errors.New("Error adding folder to filesystem")
	}
	dr.helper.commitTransaction(ctx, tranID)
	return p, nil
}

// Update updates fields of a single document from the database with the given unique identifier
// Returns an error if database fails to provide service
func (dr DocumentRepository) Update(id string, p domain.Document) error {
	return errors.New("not implemented")
}

// Delete selects a single document from the database with the given unique identifier
// Returns an error if database fails to provide service
func (dr DocumentRepository) Delete(id string) error {
	return errors.New("not implemented")
}
