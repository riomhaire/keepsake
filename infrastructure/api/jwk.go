package api

import (
	"encoding/json"
	"net/http"

	"github.com/riomhaire/keepsake/models"
)

func (this *RestAPI) HandleJWKPublicGet(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	publickey, err := this.JWKEncoder.Encode()
	if err != nil {
		handleJWTError(w, http.StatusUnauthorized, models.JWTErrorResponse{"No Token", err.Error()})
		return
	}

	json.NewEncoder(w).Encode(publickey)

	return
}
