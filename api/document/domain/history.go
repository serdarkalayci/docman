package domain

import "time"

// History represents a historical record of a document.
type History struct {
	// ID is the unique identifier of the historical record of the document.
	ID string
	// EditedAt is the date when the document was edited.
	EditedAt time.Time
	// Content is the content of the document at that specific time.
	Content string
}