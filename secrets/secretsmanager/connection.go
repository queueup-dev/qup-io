package SecretsManager

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func CreateConnection(initializedSession *session.Session) *secretsmanager.SecretsManager {
	return secretsmanager.New(initializedSession)
}
