// Package application contains the application logic of the document management system.
package application

import "github.com/serdarkalayci/docman/api/document/domain"

// DocumentRepository is the interface that we expect to be fulfilled to be used as a backend for Document Service
type DocumentRepository interface {
	List(spaceID string) ([]DocumentListItem, error)
	Add(document domain.Document, parentID string, spaceID string) (string, error)
	Get(documentID string) (domain.Document, error)
	Update(documentID string, document domain.Document) error
	Delete(documentID string) error
}

// DocumentListItem represents the list of hierarchical documents by rows
type DocumentListItem struct {
	ID        string
	Name      string
	ParentID  string
	Depth	 int
}

// DocumentService represents the struct which contains a DocumentRepository and exports methods to access the data
type DocumentService struct {
	documentRepo DocumentRepository
}

// NewDocumentService creates a new DocumentService instance and sets its repository
func NewDocumentService(dr DocumentRepository) DocumentService {
	if dr == nil {
		panic("missing documentRepository")
	}
	return DocumentService{
		documentRepo: dr,
	}
}

// List loads all the data from the included repository from the given space and returns them
// Returns an error if the repository returns one
func (ps DocumentService) List(spaceID string) ([]domain.DocumentTreeItem, error) {
	doclist, err := ps.documentRepo.List(spaceID)
	if err != nil {
		return nil, err
	}
	documentTree := []domain.DocumentTreeItem{}
	for index, item := range doclist {
		if item.ParentID == "00000000-0000-0000-0000-000000000000" { // just the root documents should start the recursive call
			documentTree = append(documentTree, buildDocumentTree(doclist, index))
		}
	}
	return documentTree, nil
}

func buildDocumentTree(doclist []DocumentListItem, index int) domain.DocumentTreeItem {
	currentItem := domain.DocumentTreeItem{
		ID: (doclist)[index].ID,
		Name: (doclist)[index].Name,
		Children: []domain.DocumentTreeItem{},
	}
	for index, item := range doclist {
		if currentItem.ID == item.ParentID {
			currentItem.Children = append(currentItem.Children, buildDocumentTree(doclist, index))
		}
	}
	return currentItem
}

// Add adds a new document to the included repository, and returns it
// Returns an error if the repository returns one
func (ps DocumentService) Add(p domain.Document, parentID string, spaceID string) (domain.Document, error) {
	documentID, err := ps.documentRepo.Add(p, parentID, spaceID)
	p.ID = documentID
	return p, err
}

// Get selects the document from the included repository with the given unique identifier, and returns it
// Returns an error if the repository returns one
func (ps DocumentService) Get(id string) (domain.Document, error) {
	document, err := ps.documentRepo.Get(id)
	return document, err
}

// Update updates the document on the included repository with the given unique identifier, and returns it
// Returns an error if the repository returns one
func (ps DocumentService) Update(id string, p domain.Document) error {
	err := ps.documentRepo.Update(id, p)
	return err
}

// Delete deletes the document from the included repository with the given unique identifier
// Returns an error if the repository returns one
func (ps DocumentService) Delete(id string) error {
	err := ps.documentRepo.Delete(id)
	return err
}
