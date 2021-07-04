package http

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func extractUintPathVar(pathVarName string, w http.ResponseWriter, req *http.Request) (uint, bool) {
	extractedUint, err := strconv.ParseUint(mux.Vars(req)[pathVarName], 10, 32)
	if err != nil {
		log.Printf("Failed to extract uint path variable as %s from the url, err={%s}\n", pathVarName, err)
		handleError(w, err, http.StatusBadRequest)
		return 0, false
	}
	return uint(extractedUint), true
}

func extractUintRequestParam(reqParamName string, w http.ResponseWriter, req *http.Request) (uint, bool) {
	extractedUint, err := strconv.ParseUint(req.URL.Query().Get(reqParamName), 10, 32)
	if err != nil {
		log.Printf("Failed to extract uint request parameter as %s from the url, err={%s}\n", reqParamName, err)
		handleError(w, err, http.StatusBadRequest)
		return 0, false
	}
	return uint(extractedUint), true
}
