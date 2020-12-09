package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"testing"
)

func TestIsConditionalCheckFailedError(t *testing.T) {

	testException := &dynamodb.ConditionalCheckFailedException{}

	result := isConditionalCheckFailedError(testException)

	if !result {
		t.Fail()
	}

	testTransactionalExceptionReasonCode := "ConditionalCheckFailed"
	testTransactionalException := &dynamodb.TransactionCanceledException{
		CancellationReasons: []*dynamodb.CancellationReason{
			{
				Code: &testTransactionalExceptionReasonCode,
			},
		},
	}

	result = isConditionalCheckFailedError(testTransactionalException)

	if !result {
		t.Fail()
	}

}
