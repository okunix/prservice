package endpoints

import (
	"encoding/json"
	"io"
	"net/http"
)

// incomplete interface
type ApiEndpoint func(w http.ResponseWriter, r *http.Request) error

func (a ApiEndpoint) Unwrap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			WriteJson(w, 400, err)
			return
		}
	}
}

func WriteJson[T any](w http.ResponseWriter, statusCode int, data T) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func ReadJson[T any](r io.Reader, dest *T) error {
	return json.NewDecoder(r).Decode(dest)
}
