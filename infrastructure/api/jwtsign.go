package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/riomhaire/keepsake/models"
)

func (this *RestAPI) HandleSignJSONViaRSA(w http.ResponseWriter, req *http.Request) {
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
		handleJWTError(w, http.StatusUnauthorized, models.JWTErrorResponse{Error: "Unauthorized", Description: err.Error()})
		return
	}
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

	// OK Sign
	token, err := this.JWTEncoderDecoder.Sign(claims)
	if err != nil {
		handleJWTError(w, http.StatusBadRequest, models.JWTErrorResponse{Error: "Bad Request", Description: err.Error()})
		return
	}

	content := models.JWTSignResponse{Token: token}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json.NewEncoder(w).Encode(content)
}
