package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DuplicateEntryException struct {
	Message string
}

const (
	awsErrConditionalCheckFailed = "ConditionalCheckFailed"
)

func (d DuplicateEntryException) Error() string {
	return d.Message
}

func isConditionalCheckFailedError(err error) bool {
	castError, ok := err.(*dynamodb.TransactionCanceledException)

	if !ok {
		return false
	}

	for _, reason := range castError.CancellationReasons {
		if *reason.Code == awsErrConditionalCheckFailed {
			return true
		}
	}

	return false
}
