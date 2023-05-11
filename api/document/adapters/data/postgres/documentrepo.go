package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/serdarkalayci/docman/api/document/adapters/data/postgres/dao"
	"github.com/serdarkalayci/docman/api/document/adapters/data/postgres/mappers"
	"github.com/serdarkalayci/docman/api/document/application"
	"github.com/serdarkalayci/docman/api/document/domain"
	"go.opentelemetry.io/otel"
)

// DocumentRepository holds the arangodb client and database name for methods to use
type DocumentRepository struct {
	db *pgx.Conn
}

func newDocumentRepository(database *pgx.Conn) DocumentRepository {
	return DocumentRepository{
		db: database,
	}
}

// List loads all the document records from tha database and returns it
// Returns an error if database fails to provide service
func (dr DocumentRepository) List(ctx context.Context, spaceID string) ([]application.DocumentListItem, error) {
	ctx, childSpan := otel.Tracer("Docman").Start(ctx, "Data:DocumentRepository:List")
	defer childSpan.End()
	documentList := []application.DocumentListItem{}
	rows, err := dr.db.Query(ctx, `WITH RECURSIVE document_hierarchy AS ( 
		SELECT id, name, parent_id, 1 AS depth
		FROM documents
		WHERE parent_id  IS NULL AND space_id = $1 -- root node(s)
		UNION ALL
		SELECT d.id, d.name, d.parent_id , dh.depth + 1 AS depth
		FROM documents d
		JOIN document_hierarchy dh ON d.parent_id  = dh.id
		where d.space_id = $1
	  )
	  SELECT id, name, parent_id , depth
	  FROM document_hierarchy
	  ORDER BY depth, id`, spaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, parentID uuid.UUID
		var name string
		var depth int
		err := rows.Scan(&id, &name, &parentID, &depth)
		if err != nil {
			return nil, err
		}

		documentList = append(documentList, application.DocumentListItem{ID: id.String(), Name: name, ParentID: parentID.String(), Depth: depth })
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return documentList, nil
}

// Get selects a single document from the database with the given unique identifier
// Returns an error if database fails to provide service
func (dr DocumentRepository) Get(documentID string) (domain.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	documentDAO := dao.DocumentDAO{}
	dr.db.QueryRow(ctx, "SELECT id, name, content, parent_id, space_id, created_at FROM documents WHERE id = $1", documentID).Scan(&documentDAO.ID, &documentDAO.Name, &documentDAO.Content, &documentDAO.ParentID, &documentDAO.SpaceID, &documentDAO.CreatedAt)
	return mappers.MapDocumentDAO2Document(documentDAO), nil
}

// AddItem adds a new document or a new folder to the underlying database.
// It returns the document inserted on success or error
func (dr DocumentRepository) Add(document domain.Document, parentID string, spaceID string) (string, error) {
	return "", errors.New("not implemented")
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
