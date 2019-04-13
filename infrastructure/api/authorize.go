package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/riomhaire/keepsake/models/oauth2"
)

var APPLICATION_JSON = "application/json"
var FORM_ENCODED = "application/x-www-form-urlencoded"

func (this *RestAPI) HandleAuthorize(w http.ResponseWriter, req *http.Request) {
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
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{Error: "Invalid Request", Description: "invalid content type"})
		return
	}
	// If we have an Authorize header use it - trumps payload
	clientid, clientsecret, err := parseBasicAuthorizeHeader(req.Header.Get("Authorization"))
	if err == nil && len(clientid) > 0 {
		// IF Client id was in the header then use that
		authorizeRequest.ClientID = clientid
		authorizeRequest.ClientSecret = clientsecret
	}

	// Call interactor - which one is dependent on whether password is present and claims
	var token string
	var authorizeResponse oauth2.AuthorizeResponse

	//	log.Println("Authorize", authorizeRequest.GrantType, authorizeRequest.ClientID, authorizeRequest.ClientSecret)

	// lookup client id
	clientInfo, err := this.ClientStore.FindClientCredential(authorizeRequest.ClientID)
	if err != nil {
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{Error: "Invalid Request", Description: err.Error()})
		return
	}
	// Compare secret
	if clientInfo.ClientSecret != authorizeRequest.ClientSecret {
		// Generate error (if not match)
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{Error: "Invalid Request", Description: "someones forgot something"})
		return
	}

	// Build token claims(if secret match)
	claims := make(map[string]interface{})
	claims["iss"] = this.Configuration.Issuer
	claims["sub"] = authorizeRequest.ClientID
	claims["scope"] = clientInfo.Scope

	// Generate token
	token, err = this.TokenEncoderDecoder.Sign(claims)

	// Set result
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		handleError(w, http.StatusUnauthorized, oauth2.ErrorResponse{Error: "Invalid Request", Description: err.Error()})
	} else {
		authorizeResponse.AccessToken = token
		authorizeResponse.ExpiresIn = this.Configuration.TimeToLiveSeconds
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

/**
* Parses a basic authorization header and splits into client and secret. Returns error
* if not
**/
func parseBasicAuthorizeHeader(authorization string) (clientid, clientsecret string, err error) {
	auth := strings.SplitN(authorization, " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		err = errors.New("Basic encoding expected")
		return
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		err = errors.New("Basic encoding has two pieces")
	} else {
		clientid = pair[0]
		clientsecret = pair[1]
	}
	return
}
