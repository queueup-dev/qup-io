package oauth2

import (
	"context"
	"golang.org/x/oauth2/clientcredentials"
)

type ClientCredentials struct {
	ClientId     string
	ClientSecret string
	TokenUrl     string
	Scopes       []string
	Context      context.Context
}

type Token struct {
	AccessToken  string
	RefreshToken *string
}

type Scopes []string

func (c ClientCredentials) AccessToken() (*Token, error) {

	config := &clientcredentials.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		TokenURL:     c.TokenUrl,
		Scopes:       c.Scopes,
	}

	token, err := config.TokenSource(c.Context).Token()

	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  token.AccessToken,
		RefreshToken: &token.RefreshToken,
	}, nil
}

func NewClientCredentials(tokenUrl string, clientId string, clientSecret string, scopes Scopes) *ClientCredentials {
	return &ClientCredentials{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		TokenUrl:     tokenUrl,
		Scopes:       scopes,
		Context:      context.Background(),
	}
}
