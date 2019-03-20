package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/riomhaire/keepsake/models/oauth2"
)

func (this *RestAPI) HandleVerify(w http.ResponseWriter, req *http.Request) {
	tokenString, ok := req.URL.Query()["token"]
	tokenValue := ""
	if !ok || len(tokenString) == 0 {
		// OK lets usee bearer
		reqToken := req.Header.Get("Authorization")
		if len(reqToken) == 0 {
			handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{Error: "Unauthorized", Description: "someone forgot parameter"})
			return
		}
		splitToken := strings.Split(reqToken, " ")
		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{Error: "Unauthorized", Description: "unsupported token type"})
			return

		}
		tokenValue = splitToken[1]

	} else {
		tokenValue = tokenString[0]
	}
	// Verify
	token, err := this.TokenEncoderDecoder.Decode(tokenValue)
	if err != nil {
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{Error: "Unauthorized", Description: err.Error()})
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(token)
	w.WriteHeader(http.StatusOK) // unprocessable entity

}
