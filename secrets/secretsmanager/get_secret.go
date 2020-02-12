package SecretsManager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func GetSecret(manager *secretsmanager.SecretsManager, secretName string) (*string, error) {

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := manager.GetSecretValue(input)

	if err != nil {
		return nil, err
	}

	return result.SecretString, nil
}
