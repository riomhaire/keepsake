package usecases

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type SimpleTokenEncoderDecoder struct {
	method jwt.SigningMethod
	secret string
	ttl    int32
}

func NewTokenEncoderDecoder(method jwt.SigningMethod, secret string, timetolive int32) (encoderDecoder *SimpleTokenEncoderDecoder) {
	encoderDecoder = &SimpleTokenEncoderDecoder{method, secret, timetolive}

	return
}

func (this *SimpleTokenEncoderDecoder) Sign(claims jwt.MapClaims) (jwtString string, err error) {
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Second * time.Duration(this.ttl)).Unix()
	claims["jti"] = uuid.NewV4()

	token := jwt.New(this.method)
	token.Claims = claims

	// Sign and get the complete encoded token as a string using the secret
	jwtString, err = token.SignedString([]byte(this.secret))
	return
}

// Decode
func (this *SimpleTokenEncoderDecoder) Decode(tokenString string) (claims jwt.MapClaims, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(this.secret), nil
	})
	if err != nil {
		return
	}
	// return claims
	claims = token.Claims.(jwt.MapClaims)

	return
}
