package models

import (
	"github.com/dgrijalva/jwt-go"
)

type StorageInteractor interface {
	FindClientCredential(clientID string) (clientCredential ClientCredential, err error)
	FindPublicPrivateKey(clientID string) (certificateCredential CertificateCredential, err error)
}

type TokenEncoderDecoder interface {
	Sign(claims jwt.MapClaims) (jwt string, err error)
	Decode(tokenString string) (claims jwt.MapClaims, err error)
}

type JWTEncoderDecoder interface {
	Sign(claims jwt.MapClaims) (jwt string, err error)
	Decode(tokenString string) (claims jwt.MapClaims, err error)
}
