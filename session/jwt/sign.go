package Jwt

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
)

func Sign(identityClaims *IdentityClaims, privateKey *rsa.PrivateKey) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, identityClaims)
	return token.SignedString(privateKey)
}
