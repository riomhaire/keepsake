package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/riomhaire/keepsake/models"
)

func (r *RestAPI) HandleSignJSONViaRSA(w http.ResponseWriter, req *http.Request) {
	bearer := "bearer "
	bearer1 := "Bearer "

	// Verify Authorization token (bearer)
	authorizationToken := req.Header.Get("Authorization")
	if len(authorizationToken) == 0 || !(strings.HasPrefix(authorizationToken, bearer) || strings.HasPrefix(authorizationToken, bearer1)) {
		handleJWTError(w, http.StatusUnauthorized, models.JWTErrorResponse{"Unauthorized", "you need a valid authorization token to use this api"})
		return
	}
	// strip off bearer and verify
	authorizationToken = string(authorizationToken[len(bearer):])

	// Should check claims has permissions/roles etc
	_, err := r.TokenEncoderDecoder.Decode(authorizationToken)
	if err != nil {
		handleJWTError(w, http.StatusUnauthorized, models.JWTErrorResponse{"Unauthorized", err.Error()})
		return
	}
	// Convert JSON To Claims
	var claims = make(map[string]interface{})
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleJWTError(w, http.StatusBadRequest, models.JWTErrorResponse{"Bad Request", err.Error()})
		return
	}

	err = json.Unmarshal(body, &claims)
	if err != nil {
		handleJWTError(w, http.StatusBadRequest, models.JWTErrorResponse{"Bad Request", err.Error()})
		return
	}

	// Lookup issuer
	issuer := claims["iss"]
	if issuer == nil {
		handleJWTError(w, http.StatusBadRequest, models.JWTErrorResponse{"Bad Request", "No Issuer"})
		return
	}
	certificates, err := r.ClientStore.FindPublicPrivateKey(issuer.(string))
	if err != nil {
		handleJWTError(w, http.StatusBadRequest, models.JWTErrorResponse{"Bad Request", err.Error()})
		return
	}
	// Sign Content if we have a private key
	if len(certificates.PrivateKey) == 0 {
		handleJWTError(w, http.StatusBadRequest, models.JWTErrorResponse{"Bad Request", "No Private Certificate For That Issuer"})
		return
	}

	// OK Sign
	token, err := r.JWTEncoderDecoder.Sign(certificates.PrivateKey, claims)
	if err != nil {
		handleJWTError(w, http.StatusBadRequest, models.JWTErrorResponse{"Bad Request", err.Error()})
		return
	}

	content := models.JWTSignResponse{token}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json.NewEncoder(w).Encode(content)
}
