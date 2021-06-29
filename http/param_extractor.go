package http

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func extractParamId(paramName string, w http.ResponseWriter, req *http.Request) (uint, bool) {
	idStr, err := strconv.ParseUint(mux.Vars(req)[paramName], 10, 32)
	if err != nil {
		log.Printf("Failed to extract uint as id from the url, err=%s\n", err)
		handleError(w, err, http.StatusBadRequest)
		return 0, false
	}
	return uint(idStr), true
}
