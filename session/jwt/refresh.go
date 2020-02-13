package Jwt

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func Refresh(
	token string,
	publicKey *rsa.PublicKey,
	privateKey *rsa.PrivateKey,
	ttl int64,
) (*string, *IdentityClaims, error) {
	claims := IdentityClaims{}

	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	claims.ExpiresAt = time.Now().Unix() + ttl
	refreshed, err := Sign(&claims, privateKey)

	return &refreshed, &claims, nil
}
