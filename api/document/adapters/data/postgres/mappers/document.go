// Package mappers contains the funtions that maps DAO objects to domain objects and visa versa.
package mappers

import (
	"github.com/serdarkalayci/docman/api/document/adapters/data/postgres/dao"
	"github.com/serdarkalayci/docman/api/document/domain"
)

// MapDocumentDAO2Document maps dao document to domain document
func MapDocumentDAO2Document(pd dao.DocumentDAO) domain.Document {
	return domain.Document{
		ID:      pd.ID.String(),
		Name:    pd.Name,
		Content: pd.Content,
		ParentID: pd.ParentID.String(),
	}
}

// MapDocument2DocumentDAO maps domain document to dao document
func MapDocument2DocumentDAO(p domain.Document) dao.DocumentDAO {
	return dao.DocumentDAO{
		Name:    p.Name,
		Content: p.Content,
	}
}

