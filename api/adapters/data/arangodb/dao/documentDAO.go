package dao

import "time"

// DocumentDAO represents the struct of document type to be stored in mongoDB
type DocumentDAO struct {
	// ID is the unique identifier of the document.
	ID string `json:"_id"`
	// Key is the unique identifier of the document.
	Key string `json:"_key"`
	// Name is the name of the document.
	Name string `json:"name"`
	// Content is the content of the document.
	Path string
	// Content is the content of the document.
	Content string `json:"content"`
	// DocumentHistory is the history of the document.
	DocumentHistory []HistoryDAO
}

type HistoryDAO struct {
	// ID is the unique identifier of the historical record of the document.
	ID string
	// EditedBy is the user who edited the document.
	EditedBy string
	// EditedAt is the date when the document was edited.
	EditedAt time.Time
}

type FolderDAO struct {
	// ID is the unique identifier of the folder.
	ID string `json:"_id"`
	// Key is the unique identifier of the folder.
	Key string `json:"_key"`
	// Name is the name of the document.
	Name string `json:"name"`
	// ParentFolderID is the unique identifier of the parent folder.
	ParentFolderID string `json:"parentFolderID"`
}

type FolderTreeDAO struct {
	// CurrentFolder is the current folder.
	CurrentFolder FolderDAO
	// SubFolders is the subfolders of the current folder.
	SubFolders []FolderDAO
	// Documents is the documents of the current folder.
	Documents []DocumentDAO
}
