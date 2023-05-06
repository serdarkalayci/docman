// Package domain describes the domain model of the document management system.
package domain

import "time"

// Document represents a document to be indexed.
type Document struct {
	// ID is the unique identifier of the document.
	ID string
	// Name is the name of the document.
	Name string
	// ParentID is the unique identifier of the parent document.
	ParentID string
	// Content is the content of the document.
	Content string
	// CreatedAt is the date when the document was created.
	CreatedAt time.Time
}

// DocumentTreeItem represents a document to be displayed in the document tree.
type DocumentTreeItem struct {
	// ID is the unique identifier of the document.
	ID string
	// Name is the name of the document.
	Name string
	// Children is the child documents of the document.
	Children []DocumentTreeItem
}




