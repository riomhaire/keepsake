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
	jwks = make(map[string]interface{})
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
	jwk, err := gojwk.PublicKey(publicKey)
	if err != nil {
		return
	}

	jwks["kid"] = key

	jwks["kty"] = jwk.Kty
	jwks["n"] = jwk.N
	jwks["e"] = jwk.E

	return
}
