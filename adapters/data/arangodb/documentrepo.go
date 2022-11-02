package arangodb

import (
	"context"
	"errors"
	driver "github.com/arangodb/go-driver"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/serdarkalayci/docman/adapters/data/arangodb/mappers"
	"github.com/serdarkalayci/docman/application"
	"github.com/serdarkalayci/docman/domain"
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
func (pr DocumentRepository) List() ([]domain.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	//var documentDAO dao.DocumentDAO
	documentDAOs, err := pr.helper.Find(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting documents")
		return nil, errors.New("Error getting documents")
	}
	documents := make([]domain.Document, 0)
	for _, documentDAO := range documentDAOs {
		document := mappers.MapDocumentDAO2Document(documentDAO)
		documents = append(documents, document)
	}
	return documents, nil
}

// Add adds a new document to the underlying database.
// It returns the document inserted on success or error
func (pr DocumentRepository) Add(p domain.Document) (domain.Document, error) {
	pass := mappers.MapDocument2DocumentDAO(p)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := pr.helper.InsertOne(ctx, pass)
	if err != nil {
		log.Error().Err(err).Msg("Error while writing user")
		return domain.Document{}, errors.New("Cannot insert the document")
	}
	log.Info().Msgf("User written: %s", result)
	p.ID = pass.ID
	return p, nil
}

// Get selects a single document from the database with the given unique identifier
// Returns an error if database fails to provide service
func (pr DocumentRepository) Get(id string) (domain.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	documentDAO, err := pr.helper.FindOne(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting document")
		return domain.Document{}, &application.ErrorCannotFinddocument{ID: id}
	}
	return mappers.MapDocumentDAO2Document(documentDAO), nil
}

// Update updates fields of a single document from the database with the given unique identifier
// Returns an error if database fails to provide service
func (pr DocumentRepository) Update(id string, p domain.Document) error {
	p.ID = id
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pDAO := mappers.MapDocument2DocumentDAO(p)
	upDoc := bson.D{{Key: "$set", Value: pDAO}}
	result, err := pr.helper.UpdateOne(ctx, id, upDoc)
	if err != nil {
		log.Error().Err(err).Msgf("Error updating the document with ID: %s", id)
		return errors.New("Error updating the document")
	}
	if result != 1 {
		log.Error().Err(err).Msgf("Could not found the document with ID: %s", id)
		return &application.ErrorCannotFinddocument{ID: id}
	}
	return nil
}

// Delete selects a single document from the database with the given unique identifier
// Returns an error if database fails to provide service
func (pr DocumentRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := pr.helper.DeleteOne(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("Error deleting document with ID: %s", id)
		return errors.New("Error deleting the document")
	}
	if result != 1 {
		log.Error().Err(err).Msgf("Could not found the document with ID: %s", id)
		return &application.ErrorCannotFinddocument{ID: id}
	}
	return nil
}
