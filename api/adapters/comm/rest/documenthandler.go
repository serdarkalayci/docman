package rest

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/serdarkalayci/docman/api/adapters/comm/rest/dto"
	"github.com/serdarkalayci/docman/api/adapters/comm/rest/mappers"
	"github.com/serdarkalayci/docman/api/adapters/comm/rest/middleware"
	"github.com/serdarkalayci/docman/api/application"
)

type validateddocument struct{}

// swagger:route GET /folder folder GetFolder
// Return all the documents
// responses:
//	200: OK
//	500: errorResponse

// GetFolder gets all the documents and folders inside the specified folder. If no folder is specified, it returns the root folder.
func (ctx *APIContext) GetFolder(rw http.ResponseWriter, r *http.Request) {
	span := createSpan("docman.GetFolder", r)
	defer span.Finish()
	// parse the document id from the url
	vars := mux.Vars(r)
	id := vars["id"]
	DocumentService := application.NewDocumentService(ctx.documentRepo)
	folder, err := DocumentService.List(id)
	if err != nil {
		respondWithError(rw, r, 500, "Cannot get folder contents from database")
	} else {
		fr := mappers.MapFolder2FolderResponseDTO(folder)
		respondWithJSON(rw, r, 200, fr)
	}

}

// swagger:route POST /document document Adddocument
// Adds a new document
// responses:
//	201: Created
//	500: errorResponse

// Adddocument adds a new documents to the Titanic
func (ctx *APIContext) AddDocument(rw http.ResponseWriter, r *http.Request) {
	span := createSpan("docman.Add", r)
	defer span.Finish()
	// Get document data from payload
	documentDTO := r.Context().Value(validateddocument{}).(dto.DocumentRequestDTO)
	document := mappers.MapDocumentRequestDTO2Document(documentDTO)
	// parse the document id from the url
	vars := mux.Vars(r)
	parentID := vars["id"]
	DocumentService := application.NewDocumentService(ctx.documentRepo)
	document, err := DocumentService.Add(document, parentID)
	if err != nil {
		respondWithError(rw, r, 500, err.Error())
	} else {
		pDTO := mappers.MapDocument2DocumentResponseDTO(document)
		respondWithJSON(rw, r, 201, pDTO)
	}
}

// swagger:route GET /document/{id} document GetDocument
// Return the document with the given id
// responses:
//	200: OK
//  400: Bad Request
//	500: errorResponse

// GetDocument gets the documents of the Titanic with the given id
func (ctx *APIContext) GetDocument(rw http.ResponseWriter, r *http.Request) {
	span := createSpan("docman.GetOne", r)
	defer span.Finish()

	// parse the document id from the url
	vars := mux.Vars(r)
	id := vars["id"]
	DocumentService := application.NewDocumentService(ctx.documentRepo)
	document, err := DocumentService.Get(id)
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
		pDTO := mappers.MapDocument2DocumentResponseDTO(document)
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
func (ctx *APIContext) UpdateDocument(rw http.ResponseWriter, r *http.Request) {
	span := createSpan("docman.Update", r)
	defer span.Finish()

	// parse the document id from the url
	vars := mux.Vars(r)
	id := vars["id"]
	// Get document data from payload
	documentDTO := r.Context().Value(validateddocument{}).(dto.DocumentRequestDTO)
	document := mappers.MapDocumentRequestDTO2Document(documentDTO)
	DocumentService := application.NewDocumentService(ctx.documentRepo)
	err := DocumentService.Update(id, document)
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
		respondEmpty(rw, r, 201)
	}
}

// swagger:route DELETE /document/{id} document DeleteDocument
// Deletes the document with the given id
// responses:
//	200: OK
//  400: Bad Request
//	500: errorResponse

// DeleteDocument deletes the documents of the Titanic with the given id
func (ctx *APIContext) DeleteDocument(rw http.ResponseWriter, r *http.Request) {
	span := createSpan("docman.Delete", r)
	defer span.Finish()

	// parse the document id from the url
	vars := mux.Vars(r)
	id := vars["id"]
	DocumentService := application.NewDocumentService(ctx.documentRepo)
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
func (ctx *APIContext) MiddlewareValidateNewDocument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user, err := middleware.ExtractAdddocumentPayload(r)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		// validate the user
		errs := ctx.validation.Validate(user)
		if errs != nil && len(errs) != 0 {
			log.Error().Err(errs[0]).Msg("Error validating the document")

			// return the validation messages as an array
			respondWithJSON(rw, r, http.StatusUnprocessableEntity, errs.Errors())
			return
		}

		// add the rating to the context
		ctx := context.WithValue(r.Context(), validateddocument{}, *user)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
