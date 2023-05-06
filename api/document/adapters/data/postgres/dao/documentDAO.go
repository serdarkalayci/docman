package dao

import (
	"time"

	"github.com/google/uuid"
)

// DocumentDAO represents the struct of document type to be stored in mongoDB
type DocumentDAO struct {
	// ID is the unique identifier of the document.
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	// ParentID is the unique identifier of the parent document.
	ParentID uuid.UUID
	// Name is the name of the document.
	Name string 
	// Content is the content of the document.
	Content string 
	// SpaceID is the unique identifier of the space.
	SpaceID uuid.UUID
	// CreatedAt is the date when the document was created.
	CreatedAt time.Time

}


