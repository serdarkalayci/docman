package domain

import (
	"time"
)

// Document represents a document to be indexed.
type Document struct {
	// ID is the unique identifier of the document.
	ID string
	// Name is the name of the document.
	Name string
	// Content is the content of the document.
	Path string
	// Content is the content of the document.
	Content string
	// DocumentHistory is the history of the document.
	DocumentHistory []History
}

type History struct {
	// ID is the unique identifier of the historical record of the document.
	ID string
	// EditedBy is the user who edited the document.
	EditedBy string
	// EditedAt is the date when the document was edited.
	EditedAt time.Time
}

type Folder struct {
	// Name is the name of the folder.
	Name string
	// Path is the path of the folder.
	Path string
	// Folders is the folders in the folder.
	Folders []Folder
	// Documents is the documents in the folder.
	Documents []Document
}
