package auth

import (
	"time"

	jwtv4 "github.com/golang-jwt/jwt/v4"
)

const (
	issuer = `service-notifications`
)

func CheckJWT(secret string) func(token *jwtv4.Token) (interface{}, error) {
	return func(token *jwtv4.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}

func MakeJWT(secret string) string {
	farFromRefuge := time.Date(2222, 2, 2, 2, 22, 22, 0, time.UTC)
	claims := &jwtv4.RegisteredClaims{
		ExpiresAt: jwtv4.NewNumericDate(farFromRefuge),
		Issuer:    issuer,
	}

	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return signedString
}
