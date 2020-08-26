package main

import (
	"context"
	"fmt"
)

type SomeStruct struct {
	Foo int
}

var (
	gatewaypayload  = []byte("{\"HTTPMethod\":\"awesome-header\",\"Body\":\"{\\\"Foo\\\":123}\"}")
	snspayload      = []byte("[{\"EventVersion\":\"zxc\",\"SNS\":{\"Message\":\"{\\\"Foo\\\":1234}\"}},{\"EventVersion\":\"zxc\",\"SNS\":{\"Message\":\"{\\\"Foo\\\":4321}\"}}]")
	firehosepayload = []byte("{\"Region\":\"blaat\",\"records\":[{\"RecordID\":\"zxc\",\"Data\":\"eyJGb28iOjEyM30=\"},{\"RecordID\":\"asd\",\"Data\":\"eyJGb28iOjMyMX0=\"}]}")
)

func main() {
	obj, err := CompileHandler(handler, false)(context.Background(), gatewaypayload)

	if err != nil {
		panic(err)
	}

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

func main2() {
	//val, _ := json.Marshal(struct{Bar []byte}{Bar: []byte("{\"Foo\":321}")})
	//fmt.Println(string(val))

	obj, err := CompileChannelHandler(KinesisFirehoseEventChannelHandlerImpl)(context.Background(), firehosepayload)

	if err != nil {
		panic(err)
	}

	ret, err2 := obj()

	fmt.Println(ret, err)
	fmt.Println(err2)
}

func handler(ctx context.Context, event *APIGatewayProxyRequest, message *SomeStruct) (interface{}, error) {
	fmt.Println("hi:")
	fmt.Println((message).Foo)
	fmt.Println(event.HTTPMethod)

	return message, nil
}

type SnsChanPayload struct {
	SNSEventRecord SNSEventRecord
	Body           SomeStruct
}

func SNSEventRecordChannelHandlerImpl(ctx context.Context, channel chan struct {
	SNSEventRecord SNSEventRecord `lambda:"event"`
	Body           SomeStruct     `lambda:"message"`
}) (interface{}, error) {

	fmt.Println("before reading from channel")

	for mes := range channel {
		fmt.Println("reading a record from the channel")

		fmt.Println(mes.Body.Foo)
		fmt.Println(mes.SNSEventRecord.SNSEventRecord.EventVersion)
	}

	fmt.Println("after reading from the channel")

	return nil, nil
}

type FirehoseRecord struct {
	KinesisFirehoseEventRecord KinesisFirehoseEventRecord `lambda:"event"`
	Body                       SomeStruct                 `lambda:"message"`
}

func KinesisFirehoseEventChannelHandlerImpl(ctx context.Context, mainEvent *KinesisFirehoseEvent, channel chan FirehoseRecord) (interface{}, error) {

	fmt.Println("inside firehose handler")
	fmt.Printf("region: %s\n", mainEvent.Region)

	for mes := range channel {
		fmt.Println(mes.KinesisFirehoseEventRecord.RecordID)
		fmt.Println(mes.Body.Foo)
	}

	return nil, nil
}
