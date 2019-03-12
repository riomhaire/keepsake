package usecases

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTEncoderDecoder struct {
	ttl int32
}

func NewJWTEncoderDecoder(ttl int32) (encoderDecoder *JWTEncoderDecoder) {
	encoderDecoder = &JWTEncoderDecoder{ttl}

	return
}

func (s *JWTEncoderDecoder) Sign(pemString string, claims jwt.MapClaims) (jwtString string, err error) {

	block, _ := pem.Decode([]byte(pemString))
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	// create a signer for rsa 256
	jwtSigner := jwt.New(jwt.GetSigningMethod("RS256"))

	// set the expire time
	// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
	claims["exp"] = time.Now().Add(time.Second * time.Duration(s.ttl)).Unix()
	jwtSigner.Claims = claims

	jwtString, err = jwtSigner.SignedString(key)
	return
}

func (s *JWTEncoderDecoder) Decode(pemString string, tokenString string) (claims jwt.MapClaims, err error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil || block.Type != "PUBLIC KEY" {
		err = errors.New("failed to decode PEM block containing public key")
		return
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	// create a validator for rsa 256
	parts := strings.Split(tokenString, ".")

	method := jwt.GetSigningMethod("RS256")
	err = method.Verify(strings.Join(parts[0:2], "."), parts[2], publicKey)
	if err != nil {
		return
	}
	// decode
	rawtoken, _ := jwt.Parse(tokenString, nil)
	if rawtoken == nil {
		err = errors.New("cannot parse token")
		return
	}

	var ok bool
	claims, ok = rawtoken.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot parse token")

	}
	return
}
