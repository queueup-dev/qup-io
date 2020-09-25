package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
)

// The return of the following handlers go straight back to aws.
type KinesisFirehoseEventChannelHandler func(ctx context.Context, qq KinesisFirehoseEvent, channel chan struct {
	KinesisFirehoseEventRecord
	SomeStruct
})
type SNSEventRecordChannelHandler func(ctx context.Context, channel chan struct {
	SNSEventRecord
	SomeStruct
}) (interface{}, error)

// For lambda's of the form 'Firehose' the second argument should implement 'GetCollection',
// of which the inner-type should be casted to the inner-type of the first argument of the channel,
// this first argument of the channel should implement 'GetMessage'.

// Previous one flattened, the result should be aggregated in the handler.
// Added 2nd error return type: last one is breaking, first one is to dlq the current record.
type SNSEventRecordChannelFlattenedHandler func(context.Context, SNSEventRecord, SomeStruct) (interface{}, error, error)

// This event doesnt have collections
type APIGatewayProxyRequestHandler func(ctx context.Context, event APIGatewayProxyRequest, message SomeStruct) (interface{}, error)

type KinesisFirehoseEvent struct {
	events.KinesisFirehoseEvent
}

func (q KinesisFirehoseEvent) GetCollection() []events.KinesisFirehoseEventRecord {
	return q.Records
}

type KinesisFirehoseEventRecord struct {
	events.KinesisFirehoseEventRecord
}

func (q KinesisFirehoseEventRecord) GetMessage() []byte {
	return q.Data
}

type APIGatewayProxyRequest struct {
	events.APIGatewayProxyRequest
}

func (q APIGatewayProxyRequest) GetMessage() []byte {
	return []byte(q.Body)
}

type SNSEventRecord struct {
	events.SNSEventRecord
}

func (q SNSEventRecord) GetMessage() []byte {
	return []byte(q.SNS.Message)
}

type KinesisEventRecord struct {
	events.KinesisEventRecord
}

func (q KinesisEventRecord) GetMessage() []byte {
	return q.Kinesis.Data
}

func test(e events.KinesisFirehoseEventRecord) KinesisFirehoseEventRecord {
	return struct {
		events.KinesisFirehoseEventRecord
	}{e}
}
