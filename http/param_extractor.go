package http

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func extractParamId(paramName string, w http.ResponseWriter, req *http.Request) (uint, bool) {
	if idStr, err := strconv.ParseUint(mux.Vars(req)[paramName], 10, 32); err != nil {
		log.Printf("Failed to extract uint as id from the url, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return 0, false
	} else {
		return uint(idStr), true
	}
}
