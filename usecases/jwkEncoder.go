package usecases

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/mendsley/gojwk"

	"github.com/riomhaire/keepsake/models"
)

type JWKEncoder struct {
	storageInteractor models.StorageInteractor
}

// .well-known/jwks.json

func NewJWKEncoder(storageInteractor models.StorageInteractor) (encoder *JWKEncoder) {
	encoder = &JWKEncoder{storageInteractor}
	return
}

func (s *JWKEncoder) Encode() (jwks map[string]interface{}, err error) {
	jwk := make(map[string]interface{})
	jwks = make(map[string]interface{})

	var values [1]map[string]interface{}

	key, err := s.storageInteractor.JWKPublicKeyName()
	if err != nil {
		return
	}

	// We have a key - find the public key
	certificates, err := s.storageInteractor.FindPublicPrivateKey(key)
	if err != nil {
		return
	}
	// Decode public key
	block, _ := pem.Decode([]byte(certificates.PublicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		err = errors.New("failed to decode pem block containing public key")
		return
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	// convert to jwks
	jwkObj, err := gojwk.PublicKey(publicKey)
	if err != nil {
		return
	}

	jwk["kid"] = key

	jwk["kty"] = jwkObj.Kty
	jwk["n"] = jwkObj.N
	jwk["e"] = jwkObj.E
	jwk["alg"] = "RS256"

	values[0] = jwk
	jwks["keys"] = values

	return
}
