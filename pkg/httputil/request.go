package httputil

import (
	"encoding/json"
	"net/http"
)

func DecodeRequest[T any](w http.ResponseWriter, r *http.Request, req *T) bool {
	if json.NewDecoder(r.Body).Decode(req) != nil {
		WriteBadRequest(w, "invalid request body")
		return false
	}
	return true
}
