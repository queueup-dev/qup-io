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

type UniqueViolationException struct {
	Message string
}

func (d DuplicateEntryException) Error() string {
	return d.Message
}

func (d ItemDoesNotExistException) Error() string {
	return d.Message
}

func (d UniqueViolationException) Error() string {
	return d.Message
}

func isConditionalCheckFailedError(err error) bool {

	_, ok := err.(*dynamodb.ConditionalCheckFailedException)

	if ok {
		return true
	}

	transactionCastError, ok := err.(*dynamodb.TransactionCanceledException)

	if !ok {
		return false
	}

	for _, reason := range transactionCastError.CancellationReasons {
		if *reason.Code == awsErrConditionalCheckFailed {
			return true
		}
	}

	return false
}
