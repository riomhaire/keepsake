package oauth2

type AuthorizeRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
	ClientID     string `json:"client_id,omitempty"`     // Client/Secret grant
	ClientSecret string `json:"client_secret,omitempty"` // Client/Secret grant
	Username     string `json:"username,omitempty"`      // Username/Password grant
	Password     string `json:"password,omitempty"`      // Username/Password grant
	Scope        string `json:"scope,omitempty"`
}

type AuthorizeResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int32  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

type ErrorResponse struct {
	Error       string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
	URI         string `json:"uri,omitempty"`
}
