package models

type ClientCredential struct {
	ClientID     string `yaml:"clientID" json:"clientID"`
	ClientSecret string `yaml:"clientSecret" json:"clientSecret"`
}

type CertificateCredential struct {
	ClientID   string `yaml:"clientID" json:"clientID"`
	PublicKey  string `yaml:"publicKey" json:"publicKey"`
	PrivateKey string `yaml:"privateKey" json:"privateKey"`
}

type Configuration struct {
	Port                      int32                   `yaml:",port" json:"port"`
	ConfigurationType         string                  `yaml:",configurationType" json:"configurationType"`
	MasterSecret              string                  `yaml:",masterSecret" json:"masterSecret"`
	TimeToLiveSeconds         int32                   `yaml:",timeToLiveSeconds" json:"timeToLiveSeconds"`
	Issuer                    string                  `yaml:",issuer" json:"issuer"`
	ClientCredentialsType     string                  `yaml:",clientCredentialsType" json:"clientCredentialsType"`
	ClientCredentials         []ClientCredential      `yaml:",clientCredentials" json:"clientCredentials"`
	CertificateCredentialType string                  `yaml:",certificateCredentialType" json:"CcertificateCredentialType"`
	CertificateCredentials    []CertificateCredential `yaml:",certificateCredential" json:"certificateCredential"`
}

type JWTSignResponse struct {
	Token string `json:"jwt,omitempty"`
}

type JWTErrorResponse struct {
	Error       string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
}
