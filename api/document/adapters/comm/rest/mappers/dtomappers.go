package mappers

import (
	"github.com/serdarkalayci/docman/api/document/adapters/comm/rest/dto"
	"github.com/serdarkalayci/docman/api/document/domain"
)

// func MapDocumentRequestDTO2Document(doc dto.DocumentRequestDTO) domain.Document {
// 	return domain.Document{
// 		ID:      doc.ID,
// 		Name:    doc.Name,
// 		Content: doc.Content,
// 	}
// }

func MapDocument2SingleDocumentResponseDTO(doc domain.Document) dto.SingleDocumentResponseDTO {
	return dto.SingleDocumentResponseDTO{
		ID:              doc.ID,
		Name:            doc.Name,
		Content:         doc.Content,
		CreatedAt:     doc.CreatedAt,
	}
}

func MapDocumentTreeItem2DocumentTreeDTO (doc domain.DocumentTreeItem) dto.DocumentTreeDTO {
	return dto.DocumentTreeDTO{
		ID:       doc.ID,
		Name:     doc.Name,
		Children: MapDocumentTreeItemArray2DocumentTreeDTOArray(doc.Children),
	}
}

func MapDocumentTreeItemArray2DocumentTreeDTOArray(doc []domain.DocumentTreeItem) []dto.DocumentTreeDTO {
	var tree []dto.DocumentTreeDTO
	for _, v := range doc {
		tree = append(tree, MapDocumentTreeItem2DocumentTreeDTO(v))
	}
	return tree
}

// func MapDocumentHistory2DocumentHistoryDTO(doc domain.History) dto.HistoryDTO {
// 	return dto.HistoryDTO{
// 		EditedBy: doc.EditedBy,
// 		EditedAt: doc.EditedAt,
// 	}
// }

// func MapDocumentHistoryArray2DocumentHistoryDTOArray(doc []domain.History) []dto.HistoryDTO {
// 	var history []dto.HistoryDTO
// 	for _, v := range doc {
// 		history = append(history, MapDocumentHistory2DocumentHistoryDTO(v))
// 	}
// 	return history
// }


