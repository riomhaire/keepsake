package storage

import (
	"encoding/base64"
	"errors"

	"github.com/riomhaire/keepsake/models"
)

type ConfigurationStorageInteractor struct {
	configuration *models.Configuration
}

func NewConfigurationStorageIntegrator(configuration *models.Configuration) (configurationStorageInteractor *ConfigurationStorageInteractor) {

	configurationStorageInteractor = &ConfigurationStorageInteractor{}
	configurationStorageInteractor.configuration = configuration

	return
}

func (s *ConfigurationStorageInteractor) FindClientCredential(clientID string) (clientCredential models.ClientCredential, err error) {

	for _, credential := range s.configuration.ClientCredentials {
		if credential.ClientID == clientID {
			clientCredential = credential
			return
		}
	}
	err = errors.New("unknown client id")
	return
}

func (s *ConfigurationStorageInteractor) FindPublicPrivateKey(clientID string) (certificateCredential models.CertificateCredential, err error) {
	for _, credential := range s.configuration.CertificateCredentials {
		if credential.ClientID == clientID {
			if len(credential.PublicKey) > 0 {
				decoded, _ := base64.StdEncoding.DecodeString(credential.PublicKey)
				certificateCredential.PublicKey = string(decoded)
			}
			if len(credential.PrivateKey) > 0 {
				decoded, _ := base64.StdEncoding.DecodeString(credential.PrivateKey)
				certificateCredential.PrivateKey = string(decoded)
			}
			// log.Println(certificateCredential.PrivateKey)
			return
		}
	}
	err = errors.New("unknown client id")
	return
}

func (s *ConfigurationStorageInteractor) JWKPublicKeyName() (key string, err error) {
	key = s.configuration.JWKKey
	return
}
