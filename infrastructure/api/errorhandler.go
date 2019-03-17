package api

import (
	"encoding/json"
	"net/http"

	"github.com/riomhaire/keepsake/models"
	"github.com/riomhaire/keepsake/models/oauth2"
)

func handleError(w http.ResponseWriter, errorCode int, content oauth2.ErrorResponse) {
	w.WriteHeader(errorCode) // unprocessable entity
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json.NewEncoder(w).Encode(content)
}

func handleJWTError(w http.ResponseWriter, errorCode int, content models.JWTErrorResponse) {
	w.WriteHeader(errorCode) // unprocessable entity
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json.NewEncoder(w).Encode(content)
}

func handleSimpleError(w http.ResponseWriter, errorCode int, content string) {
	w.WriteHeader(errorCode) // unprocessable entity
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json.NewEncoder(w).Encode(content)
}
