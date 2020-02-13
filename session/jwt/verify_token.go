package Jwt

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
)

func VerifyAndGetClaims(token string, publicKey *rsa.PublicKey) (*jwt.Token, *IdentityClaims, error) {

	claims := IdentityClaims{}

	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	return jwtToken, &claims, nil
}

func Verify(token string, key *rsa.PublicKey) (*jwt.Token, error) {

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
}
