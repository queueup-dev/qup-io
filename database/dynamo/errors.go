package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	awsErrConditionalCheckFailed = "ConditionalCheckFailed"
)

type DuplicateEntryException struct {
	Message string
}

type ItemDoesNotExistException struct {
	Message string
}

func (d DuplicateEntryException) Error() string {
	return d.Message
}

func (d ItemDoesNotExistException) Error() string {
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
