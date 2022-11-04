package middleware

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/serdarkalayci/docman/api/adapters/comm/rest/dto"
	"github.com/serdarkalayci/docman/api/application"
)

// ToJSON serializes the given interface into a string based JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

func readPayload(r *http.Request) (payload []byte, e error) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		e = &application.ErrorReadPayload{}
		log.Error().Err(err)
		return
	}
	if len(payload) == 0 {
		e = &application.ErrorPayloadMissing{}
		log.Error().Err(err)
		return
	}
	return
}

// ExtractAdddocumentPayload extracts document data from the request body
// Returns documentRequestDTO model if found, error otherwise
func ExtractAdddocumentPayload(r *http.Request) (document *dto.DocumentRequestDTO, e error) {
	payload, e := readPayload(r)
	if e != nil {
		return
	}
	err := json.Unmarshal(payload, &document)
	if err != nil {
		e = &application.ErrorParsePayload{}
		log.Error().Err(err)
		return
	}
	return
}
