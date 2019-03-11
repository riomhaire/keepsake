package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/riomhaire/keepsake/models/oauth2"
)

var APPLICATION_JSON = "application/json"
var FORM_ENCODED = "application/x-www-form-urlencoded"

func (r *RestAPI) HandleAuthorize(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)
	// Decode request
	var authorizeRequest oauth2.AuthorizeRequest

	// is this JSON or form post?
	content_type := req.Header.Get("Content-Type")
	if strings.Contains(content_type, APPLICATION_JSON) {
		_ = json.NewDecoder(req.Body).Decode(&authorizeRequest)
	} else if strings.Contains(content_type, FORM_ENCODED) {
		req.ParseForm()
		for key, value := range req.Form {
			switch key {
			case "grant_type":
				authorizeRequest.GrantType = formFieldValue(value)
			case "code":
				authorizeRequest.Code = formFieldValue(value)
			case "redirect_uri":
				authorizeRequest.RedirectURI = formFieldValue(value)
			case "client_id":
				authorizeRequest.ClientID = formFieldValue(value)
			case "client_secret":
				authorizeRequest.ClientSecret = formFieldValue(value)
			case "username":
				authorizeRequest.Username = formFieldValue(value)
			case "password":
				authorizeRequest.Password = formFieldValue(value)
			case "scope":
				authorizeRequest.Scope = formFieldValue(value)
			}
		}
	} else {
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{"Invalid Request", "invalid content type", ""})
		return
	}
	// Call interactor - which one is dependent on whether password is present and claims
	var token string
	var err error
	var authorizeResponse oauth2.AuthorizeResponse

	log.Println("Authorize", authorizeRequest.GrantType, authorizeRequest.ClientID, authorizeRequest.ClientSecret)

	// lookup client id
	clientInfo, err := r.ClientStore.FindClientCredential(authorizeRequest.ClientID)
	if err != nil {
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{"Invalid Request", err.Error(), ""})
		return
	}
	// Compare secret
	if clientInfo.ClientSecret != authorizeRequest.ClientSecret {
		// Generate error (if not match)
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{"Invalid Request", "someones forgot something", ""})
		return
	}

	// Build token claims(if secret match)
	claims := make(map[string]interface{})
	claims["iss"] = r.Configuration.Issuer
	claims["sub"] = authorizeRequest.ClientID

	// Generate token
	token, err = r.TokenEncoderDecoder.Sign(claims)

	// Set result
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{"Invalid Request", err.Error(), ""})
	} else {
		authorizeResponse.AccessToken = token
		authorizeResponse.ExpiresIn = r.Configuration.TimeToLiveSeconds
		authorizeResponse.TokenType = "bearer"

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Authorization", fmt.Sprintf("%s %s", bearerPrefix, token))
		json.NewEncoder(w).Encode(authorizeResponse)
	}
}

func formFieldValue(params []string) string {
	if len(params) == 0 {
		return ""
	}
	return params[0]
}
