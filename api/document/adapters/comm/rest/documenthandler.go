package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/serdarkalayci/docman/api/document/adapters/comm/rest/mappers"
	"github.com/serdarkalayci/docman/api/document/adapters/comm/rest/middleware"
	"github.com/serdarkalayci/docman/api/document/application"
)

type validateddocument struct{}

// swagger:route GET /folder folder GetFolder
// Return all the documents
// responses:
//	200: OK
//	500: errorResponse

// GetSpace gets the tree of all the documents inside a space.
func (apiContext *APIContext) GetSpace(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	ctx, span := createSpan(ctx, "Rest:DocumentHandler:GetSpace", r)
	defer span.End()
	// parse the document id from the url
	vars := mux.Vars(r)
	id := vars["id"]
	DocumentService := application.NewDocumentService(apiContext.documentRepo)
	folder, err := DocumentService.List(ctx, id)
	if err != nil {
		respondWithError(rw, r, 500, "Cannot get folder contents from database")
	} else {
		respondWithJSON(rw, r, 200, mappers.MapDocumentTreeItemArray2DocumentTreeDTOArray(folder))
	}
}

// swagger:route POST /document document Adddocument
// Adds a new document
// responses:
//	201: Created
//	500: errorResponse

// // Adddocument adds a new documents to the Titanic
// func (apiContext *APIContext) AddDocument(rw http.ResponseWriter, r *http.Request) {
// 	span := createSpan("docman.Add", r)
// 	defer span.Finish()
// 	// Get document data from payload
// 	documentDTO := r.Context().Value(validateddocument{}).(dto.DocumentRequestDTO)
// 	document := mappers.MapDocumentRequestDTO2Document(documentDTO)
// 	// parse the document id from the url
// 	vars := mux.Vars(r)
// 	parentID := vars["id"]
// 	DocumentService := application.NewDocumentService(apiContext.documentRepo)
// 	document, err := DocumentService.Add(document, parentID)
// 	if err != nil {
// 		respondWithError(rw, r, 500, err.Error())
// 	} else {
// 		pDTO := mappers.MapDocument2DocumentResponseDTO(document)
// 		respondWithJSON(rw, r, 201, pDTO)
// 	}
// }

// swagger:route GET /document/{id} document GetDocument
// Return the document with the given id
// responses:
//	200: OK
//  400: Bad Request
//	500: errorResponse

// GetDocument gets the documents of the Titanic with the given id
func (apiContext *APIContext) GetDocument(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	ctx, span := createSpan(ctx, "docman.GetOne", r)
	defer span.End()

	// parse the document id from the url
	vars := mux.Vars(r)
	id := vars["id"]
	DocumentService := application.NewDocumentService(apiContext.documentRepo)
	document, err := DocumentService.Get(ctx, id)
	if err != nil {
		switch err.(type) {
		case *application.ErrorIDFormat:
			respondWithError(rw, r, 400, "Cannot process with the given id")
		case *application.ErrorCannotFinddocument:
			respondWithError(rw, r, 404, "Cannot get document from database")
		default:
			respondWithError(rw, r, 500, "Internal server error")
		}
	} else {
		pDTO := mappers.MapDocument2SingleDocumentResponseDTO(document)
		respondWithJSON(rw, r, 200, pDTO)
	}
}

// swagger:route PUT /document{id} document UpdateDocument
// Updates an existing document
// responses:
//	201: Created
//  400: Bad Request
//	500: errorResponse

// UpdateDocument updates an existing documents on the Titanic
// func (apiContext *APIContext) UpdateDocument(rw http.ResponseWriter, r *http.Request) {
// 	span := createSpan("docman.Update", r)
// 	defer span.Finish()

// 	// parse the document id from the url
// 	vars := mux.Vars(r)
// 	id := vars["id"]
// 	// Get document data from payload
// 	documentDTO := r.Context().Value(validateddocument{}).(dto.DocumentRequestDTO)
// 	document := mappers.MapDocumentRequestDTO2Document(documentDTO)
// 	DocumentService := application.NewDocumentService(apiContext.documentRepo)
// 	err := DocumentService.Update(id, document)
// 	if err != nil {
// 		switch err.(type) {
// 		case *application.ErrorIDFormat:
// 			respondWithError(rw, r, 400, "Cannot process with the given id")
// 		case *application.ErrorCannotFinddocument:
// 			respondWithError(rw, r, 404, "Cannot get document from database")
// 		default:
// 			respondWithError(rw, r, 500, "Internal server error")
// 		}
// 	} else {
// 		respondEmpty(rw, r, 201)
// 	}
// }

// swagger:route DELETE /document/{id} document DeleteDocument
// Deletes the document with the given id
// responses:
//	200: OK
//  400: Bad Request
//	500: errorResponse

// DeleteDocument deletes the documents of the Titanic with the given id
func (apiContext *APIContext) DeleteDocument(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	ctx, span := createSpan(ctx, "docman.Delete", r)
	defer span.End()

	// parse the document id from the url
	vars := mux.Vars(r)
	id := vars["id"]
	DocumentService := application.NewDocumentService(apiContext.documentRepo)
	err := DocumentService.Delete(id)
	if err != nil {
		switch err.(type) {
		case *application.ErrorIDFormat:
			respondWithError(rw, r, 400, "Cannot process with the given id")
		case *application.ErrorCannotFinddocument:
			respondWithError(rw, r, 404, "Cannot get document from database")
		default:
			respondWithError(rw, r, 500, "Internal server error")
		}
	} else {
		respondEmpty(rw, r, 200)
	}
}

// MiddlewareValidateNewDocument Checks the integrity of new document in the request and calls next if ok
func (apiContext *APIContext) MiddlewareValidateNewDocument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user, err := middleware.ExtractAdddocumentPayload(r)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		// validate the user
		errs := apiContext.validation.Validate(user)
		if errs != nil && len(errs) != 0 {
			log.Error().Err(errs[0]).Msg("Error validating the document")

			// return the validation messages as an array
			respondWithJSON(rw, r, http.StatusUnprocessableEntity, errs.Errors())
			return
		}

		// add the rating to the context
		apiContext := context.WithValue(r.Context(), validateddocument{}, *user)
		r = r.WithContext(apiContext)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
