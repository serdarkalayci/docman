package mappers

import (
	"github.com/serdarkalayci/docman/api/adapters/comm/rest/dto"
	"github.com/serdarkalayci/docman/api/domain"
)

func MapDocumentRequestDTO2Document(doc dto.DocumentRequestDTO) domain.Document {
	return domain.Document{
		ID:      doc.ID,
		Name:    doc.Name,
		Content: doc.Content,
	}
}

func MapDocument2DocumentResponseDTO(doc domain.Document) dto.DocumentResponseDTO {
	return dto.DocumentResponseDTO{
		ID:              doc.ID,
		Name:            doc.Name,
		Content:         doc.Content,
		DocumentHistory: MapDocumentHistoryArray2DocumentHistoryDTOArray(doc.DocumentHistory),
	}
}

func MapDocumentHistory2DocumentHistoryDTO(doc domain.History) dto.HistoryDTO {
	return dto.HistoryDTO{
		EditedBy: doc.EditedBy,
		EditedAt: doc.EditedAt,
	}
}

func MapDocumentHistoryArray2DocumentHistoryDTOArray(doc []domain.History) []dto.HistoryDTO {
	var history []dto.HistoryDTO
	for _, v := range doc {
		history = append(history, MapDocumentHistory2DocumentHistoryDTO(v))
	}
	return history
}

func MapFolder2FolderResponseDTO(folder domain.Folder) dto.FolderResponseDTO {
	fr := dto.FolderResponseDTO{
		ID:   folder.ID,
		Name: folder.Name,
	}
	for _, d := range folder.Documents {
		fr.Documents = append(fr.Documents, MapDocument2DocumentResponseDTO(d))
	}
	for _, f := range folder.Folders {
		fr.Folders = append(fr.Folders, MapFolder2FolderResponseDTO(f))
	}
	return fr
}
