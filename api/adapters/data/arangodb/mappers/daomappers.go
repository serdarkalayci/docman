package mappers

import (
	"github.com/serdarkalayci/docman/api/adapters/data/arangodb/dao"
	"github.com/serdarkalayci/docman/api/domain"
)

// MapDocumentDAO2Document maps dao document to domain document
func MapDocumentDAO2Document(pd dao.DocumentDAO) domain.Document {
	return domain.Document{
		ID:      pd.ID,
		Name:    pd.Name,
		Content: pd.Content,
	}
}

// MapDocument2DocumentDAO maps domain document to dao document
func MapDocument2DocumentDAO(p domain.Document) dao.DocumentDAO {
	return dao.DocumentDAO{
		Name:    p.Name,
		Content: p.Content,
	}
}

func MapFolderDAO2Folder(f dao.FolderDAO) domain.Folder {
	return domain.Folder{
		ID:   f.ID,
		Name: f.Name,
	}
}

func MapFolderTreeDAO2Folder(ft dao.FolderTreeDAO) domain.Folder {
	folder := domain.Folder{
		ID:             ft.CurrentFolder.ID,
		Name:           ft.CurrentFolder.Name,
		ParentFolderID: ft.CurrentFolder.ParentFolderID,
	}
	for _, d := range ft.Documents {
		folder.Documents = append(folder.Documents, MapDocumentDAO2Document(d))
	}
	for _, f := range ft.SubFolders {
		folder.Folders = append(folder.Folders, MapFolderDAO2Folder(f))
	}
	return folder
}

func MapFolder2FolderDAO(f domain.Folder) dao.FolderDAO {
	return dao.FolderDAO{
		ID:   f.ID,
		Name: f.Name,
	}
}
