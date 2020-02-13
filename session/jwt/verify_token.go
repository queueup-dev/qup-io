package Jwt

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
)

func VerifyAndGetClaims(token string, key *rsa.PublicKey) (*jwt.Token, map[string]interface{}, error) {

	jwtToken, err := Verify(token, key)

	if err != nil {
		return nil, nil, err
	}

	return jwtToken, jwtToken.Claims.(jwt.MapClaims), nil
}

func Verify(token string, key *rsa.PublicKey) (*jwt.Token, error) {

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
}
