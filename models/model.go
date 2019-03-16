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

type ConsulConfiguration struct {
	Enabled  bool   `yaml:",enabled,omitempty" json:"enabled,omitempty"`
	Host     string `yaml:",host,omitempty" json:"host,omitempty"`
	ConsulId string `yaml:",id,omitempty" json:"id,omitempty"`
}

type Test struct {
	BaseURL string `yaml:",baseURL,omitempty" json:"baseURL,omitempty"`
}

type Configuration struct {
	ApplicationName           string                  `yaml:",applicationName,omitempty" json:"applicationName,omitempty"`
	Version                   string                  `yaml:",applicationVersion,omitempty" json:"applicationVersion,omitempty"`
	Port                      int                     `yaml:",port,omitempty" json:"port,omitempty"`
	ConfigurationType         string                  `yaml:",configurationType,omitempty" json:"configurationType,omitempty"`
	MasterSecret              string                  `yaml:",masterSecret,omitempty" json:"masterSecret,omitempty"`
	TimeToLiveSeconds         int32                   `yaml:",timeToLiveSeconds,omitempty" json:"timeToLiveSeconds,omitempty"`
	Issuer                    string                  `yaml:",issuer,omitempty" json:"issuer,omitempty"`
	ClientCredentialsType     string                  `yaml:",clientCredentialsType,omitempty" json:"clientCredentialsType,omitempty"`
	ClientCredentials         []ClientCredential      `yaml:",clientCredentials,omitempty" json:"clientCredentials,omitempty"`
	CertificateCredentialType string                  `yaml:",certificateCredentialType,omitempty" json:"CcertificateCredentialType,omitempty"`
	CertificateCredentials    []CertificateCredential `yaml:",certificateCredential,omitempty" json:"certificateCredential,omitempty"`
	Consul                    ConsulConfiguration     `yaml:",consul,omitempty" json:"consul,omitempty"`
	Test                      Test                    `yaml:",test,omitempty" json:"test,omitempty"`
}

type JWTSignResponse struct {
	Token string `json:"jwt,omitempty"`
}

type JWTErrorResponse struct {
	Error       string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
}
