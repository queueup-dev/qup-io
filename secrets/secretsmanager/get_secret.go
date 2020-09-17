package SecretsManager

import (
	"fmt"
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

func Must(secret *string, err error) string {
	if err != nil {
		panic(err)
	}

	if secret == nil {
		panic(fmt.Errorf("value of a required secret value is nil"))
	}

	return *secret
}
