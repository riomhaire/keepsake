package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/riomhaire/keepsake/models"
	"github.com/riomhaire/keepsake/models/oauth2"
)

func (this *RestAPI) HandleVerifyJWTViaRSA(w http.ResponseWriter, req *http.Request) {
	bearer := "bearer "
	bearer1 := "Bearer "

	// Verify Authorization token (bearer)
	authorizationToken := req.Header.Get("Authorization")
	if len(authorizationToken) == 0 || !(strings.HasPrefix(authorizationToken, bearer) || strings.HasPrefix(authorizationToken, bearer1)) {
		handleJWTError(w, http.StatusUnauthorized, models.JWTErrorResponse{Error: "Unauthorized", Description: "you need a valid authorization token to use this api"})
		return
	}
	// strip off bearer and verify
	authorizationToken = string(authorizationToken[len(bearer):])

	// Should check claims has permissions/roles etc
	_, err := this.TokenEncoderDecoder.Decode(authorizationToken)
	if err != nil {
		handleJWTError(w, http.StatusUnauthorized, models.JWTErrorResponse{"Unauthorized Token", err.Error()})
		return
	}

	// Now Validate jwt parameter and look up issuer so we can verify it
	// IF were a form overwrite value
	content_type := req.Header.Get("Content-Type")
	tokenString := ""
	if strings.Contains(content_type, APPLICATION_JSON) {
		// Convert JSON To Claims
		var claims = make(map[string]interface{})
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			handleJWTError(w, http.StatusBadRequest, models.JWTErrorResponse{Error: "Bad Request", Description: err.Error()})
			return
		}

		err = json.Unmarshal(body, &claims)
		if err != nil {
			handleJWTError(w, http.StatusBadRequest, models.JWTErrorResponse{Error: "Bad Request", Description: err.Error()})
			return
		}
		// Extract JWT
		tokenString = string(claims["jwt"].(string))
	}

	if len(tokenString) == 0 {
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{"Unauthorized", "someone forgot jwt payload parameter", ""})
		return
	}

	claims, err := this.JWTEncoderDecoder.Decode(tokenString)
	if err != nil {
		handleJWTError(w, http.StatusUnauthorized, models.JWTErrorResponse{"Bad Request", err.Error()})
		return

	}

	// We have the json
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(claims)
}
