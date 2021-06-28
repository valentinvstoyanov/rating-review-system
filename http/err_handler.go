package http

import (
	"encoding/json"
	"net/http"
)

func handleError(w http.ResponseWriter, err error, statusCode int) {
	type ErrorResp struct {
		Message string `json:"message"`
	}

	w.WriteHeader(statusCode)
	w.Header().Add("Content-type", "application/json")
	_ = json.NewEncoder(w).Encode(ErrorResp{Message: err.Error()})
}
