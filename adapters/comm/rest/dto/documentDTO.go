package dto

import "time"

// DocumentResponseDTO represents the struct that is returned by rest endpoints
type DocumentResponseDTO struct {

	// ID is the unique identifier of the document.
	ID string `json:"id"`
	// Name is the name of the document.
	Name string `json:"name"`
	// Content is the content of the document.
	Content string `json:"content"`
	// Path is the path of the document.
	Path string `json:"path"`
	// DocumentHistory is the history of the document.
	DocumentHistory []HistoryDTO `json:"documentHistory"`
}

type HistoryDTO struct {
	// EditedBy is the user who edited the document.
	EditedBy string `json:"editedBy"`
	// EditedAt is the date when the document was edited.
	EditedAt time.Time `json:"editedAt"`
}

// DocumentRequestDTO represents the struct that is accepted as input for the rest endpoint
type DocumentRequestDTO struct {
	// ID is the unique identifier of the document.
	ID string `json:"id"`
	// Name is the name of the document.
	Name string `json:"name"`
	// Content is the content of the document.
	Path string `json:"path"`
	// Content is the content of the document.
	Content string `json:"content"`
}
