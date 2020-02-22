package Jwt

type IdentityClaims struct {
	Audience  string `json:"aud,omitempty"     structs:"aud"`
	ExpiresAt int64  `json:"exp,omitempty"     structs:"exp"`
	Id        string `json:"jti,omitempty"     structs:"jti"`
	IssuedAt  int64  `json:"iat,omitempty"     structs:"iat"`
	Issuer    string `json:"iss,omitempty"     structs:"iss"`
	NotBefore int64  `json:"nbf,omitempty"     structs:"nbf"`
	Subject   string `json:"sub,omitempty"     structs:"sub"`
	Email     string `json:"email,omitempty"   structs:"email"`
	UserId    string `json:"user_id,omitempty" structs:"user_id"`
}

func (i IdentityClaims) Valid() error {
	return nil
}
