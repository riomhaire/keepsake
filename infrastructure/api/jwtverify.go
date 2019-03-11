package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/riomhaire/keepsake/models"
	"github.com/riomhaire/keepsake/models/oauth2"
)

func (r *RestAPI) HandleVerifyJWTViaRSA(w http.ResponseWriter, req *http.Request) {
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

	// Now Validate token parameter
	tokenString, ok := req.URL.Query()["token"]

	if !ok || len(tokenString) == 0 {
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{"Unauthorized", "someone forgot parameter", ""})
		return
	}

	// ok find issuer in token and look it up
	parts := strings.Split(tokenString[0], ".")
	if len(parts) != 3 {
		handleJWTError(w, http.StatusUnauthorized, models.JWTErrorResponse{"Unauthorized", err.Error()})
		return
	}
	// Part 1 is base64 json
	j, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		handleJWTError(w, http.StatusUnauthorized, models.JWTErrorResponse{"Unauthorized", err.Error()})
		return

	}
	claims := make(map[string]interface{})
	err = json.Unmarshal(j, &claims)
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

	// Verify Content if we have a public key
	if len(certificates.PublicKey) == 0 {
		handleJWTError(w, http.StatusBadRequest, models.JWTErrorResponse{"Bad Request", "No Public Certificate For That Issuer"})
		return
	}
	_, err = r.JWTEncoderDecoder.Decode(certificates.PublicKey, tokenString[0])
	if err != nil {
		handleJWTError(w, http.StatusBadRequest, models.JWTErrorResponse{"Bad Request", err.Error()})
		return

	}

	// We have the json
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(claims)
}
