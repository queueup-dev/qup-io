package main

import (
	"context"
	"fmt"
)

type SomeStruct struct {
	Foo int
}

var (
	gatewaypayload  = []byte("{\"Body\":\"{\\\"Foo\\\":123}\"}")
	snspayload      = []byte("[{\"EventVersion\":\"zxc\",\"SNS\":{\"Message\":\"{\\\"Foo\\\":123}\"}},{\"EventVersion\":\"zxc\",\"SNS\":{\"Message\":\"{\\\"Foo\\\":321}\"}}]")
	firehosepayload = []byte("{\"Region\":\"blaat\",\"records\":[{\"RecordID\":\"zxc\",\"SNS\":{\"Message\":\"{\\\"Foo\\\":123}\"}},{\"RecordID\":\"zxc\",\"SNS\":{\"Message\":\"{\\\"Foo\\\":321}\"}}]}")
)

func main2() {
	obj, err := CompileHandler(handler, false)(context.Background(), gatewaypayload)

	ret, err2 := obj()

	fmt.Println(ret, err)
	fmt.Println(err2)
}

func main3() {
	obj, err := CompileChannelHandler(SNSEventRecordChannelHandlerImpl)(context.Background(), snspayload)

	if err != nil {
		panic(err)
	}

	ret, err2 := obj()

	fmt.Println(ret, err)
	fmt.Println(err2)
}

func main() {
	fmt.Println(string(firehosepayload))
	obj, err := CompileChannelHandler(KinesisFirehoseEventChannelHandlerImpl)(context.Background(), firehosepayload)

	if err != nil {
		panic(err)
	}

	ret, err2 := obj()

	fmt.Println(ret, err)
	fmt.Println(err2)
}

func handler(ctx context.Context, event APIGatewayProxyRequest, message SomeStruct) (interface{}, error) {
	fmt.Printf("hi:")
	fmt.Println((message).Foo)

	return message, nil
}

type SnsChanPayload struct {
	SNSEventRecord SNSEventRecord
	Body           SomeStruct
}

func SNSEventRecordChannelHandlerImpl(ctx context.Context, channel chan struct {
	SNSEventRecord SNSEventRecord
	Body           SomeStruct
}) (interface{}, error) {

	fmt.Println("before reading from channel")

	for mes := range channel {
		fmt.Println("reading a record from the channel")

		fmt.Println(mes.Body)
		fmt.Println(mes.SNSEventRecord.SNSEventRecord.EventVersion)
	}

	fmt.Println("after reading from the channel")

	return nil, nil
}

func KinesisFirehoseEventChannelHandlerImpl(ctx context.Context, qq KinesisFirehoseEvent, channel chan struct {
	KinesisFirehoseEventRecord KinesisFirehoseEventRecord
	Body                       SomeStruct
}) {

	fmt.Println(qq.Region)

	for mes := range channel {
		fmt.Println(mes.KinesisFirehoseEventRecord.RecordID)
		fmt.Println(mes.Body)
	}
}
