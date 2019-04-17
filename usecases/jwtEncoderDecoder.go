package usecases

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"
	"github.com/riomhaire/keepsake/models"
)

type JWTEncoderDecoder struct {
	storageInteractor models.StorageInteractor
	ttl               int32
}

func NewJWTEncoderDecoder(ttl int32, storageInteractor models.StorageInteractor) (encoderDecoder *JWTEncoderDecoder) {
	encoderDecoder = &JWTEncoderDecoder{storageInteractor, ttl}

	return
}

func (s *JWTEncoderDecoder) Sign(claims jwt.MapClaims) (jwtString string, err error) {

	// Lookup issuer
	issuer := claims["iss"]
	if issuer == nil {
		err = errors.New("no issuer")
		return
	}
	certificates, err := s.storageInteractor.FindPublicPrivateKey(issuer.(string))
	if err != nil {
		return
	}
	// Sign Content if we have a private key
	if len(certificates.PrivateKey) == 0 {
		err = errors.New("no private certificate for that issuer")
		return
	}
	log.Println(certificates.PrivateKey)

	block, _ := pem.Decode([]byte(certificates.PrivateKey))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	// create a signer for rsa 256
	jwtSigner := jwt.New(jwt.GetSigningMethod("RS256"))

	// set the expire time
	// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
	claims["exp"] = time.Now().Add(time.Second * time.Duration(s.ttl)).Unix()
	jwtSigner.Claims = claims

	jwtString, err = jwtSigner.SignedString(key)
	return
}

func (s *JWTEncoderDecoder) Decode(tokenString string) (claims jwt.MapClaims, err error) {

	rawtoken, _ := jwt.Parse(tokenString, nil)
	if rawtoken == nil {
		err = errors.New("cannot parse token")
		return
	}
	// extract claims and lookup issuer
	claims, ok := rawtoken.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot parse token")
		return

	}
	issuer := string(claims["iss"].(string))
	if issuer == "" {
		err = errors.New("no issuer")
		return
	}

	certificates, err := s.storageInteractor.FindPublicPrivateKey(issuer)
	if err != nil {
		return
	}

	// Verify Content if we have a public key
	if len(certificates.PublicKey) == 0 {
		err = errors.New("no public certificate for that issuer")
		return
	}

	block, _ := pem.Decode([]byte(certificates.PublicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		err = errors.New("failed to decode pem block containing public key")
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
	return
}
