package api

import (
	"encoding/json"
	"net/http"

	"github.com/riomhaire/keepsake/models/oauth2"
)

func (r *RestAPI) HandleVerify(w http.ResponseWriter, req *http.Request) {
	tokenString, ok := req.URL.Query()["token"]

	if !ok || len(tokenString) == 0 {
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{"Unauthorized", "someone forgot parameter", ""})
		return
	}
	// Verify
	token, err := r.TokenEncoderDecoder.Decode(tokenString[0])
	if err != nil {
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{"Unauthorized", err.Error(), ""})
		return
	}
	w.WriteHeader(http.StatusOK) // unprocessable entity
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json.NewEncoder(w).Encode(token)
}
