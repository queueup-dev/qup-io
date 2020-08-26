package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

type TypeAndPointerDepth struct {
	Type         reflect.Type
	PointerDepth int
}

type ChannelEnvelopeType struct {
	Type               reflect.Type
	AwsEventRecordType struct {
		TypeAndPointerDepth
		GetMessageMethod reflect.Method
	}
	UserRecord struct {
		TypeAndPointerDepth
		UnmarshallMethod *reflect.Method
	}
}

type ChannelType struct {
	Type reflect.Type
	ChannelEnvelopeType
}

type MainEventType struct {
	Type          reflect.Type
	GetCollection *reflect.Method
	CastFunction  func(value reflect.Value) reflect.Value
}

type ChannelHandler struct {
	ChannelType
	MainEventType
	Values struct {
		ActualChannel reflect.Value
		DlqChannel    reflect.Value
	}
}

func ValidateChannel(in reflect.Type) *ChannelType {
	var c ChannelType

	if in.Kind() != reflect.Chan {
		panic(fmt.Errorf("expected the last argument of the handler to be a channel, got: %s", c.Type.Name()))
	}
	c.Type = in

	c.ChannelEnvelopeType = *ValidateEnvelope(in.Elem())

	return &c
}

func ValidateEnvelope(in reflect.Type) *ChannelEnvelopeType {
	if in.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected the inner type of the channel to be a struct of 2 elements, got: %s", in.Name()))
	}
	c := ChannelEnvelopeType{Type: in}

	c.AwsEventRecordType.Type = c.Type.Field(0).Type
	c.UserRecord.Type = c.Type.Field(1).Type

	if !c.AwsEventRecordType.Type.Implements(getMessageType) {
		panic(fmt.Errorf("the first field of the channel struct doesn't implement %s, got: %s", getMessageType.Name(), c.AwsEventRecordType.Type.Name()))
	}
	c.AwsEventRecordType.GetMessageMethod, _ = c.AwsEventRecordType.Type.MethodByName("GetMessage")

	if c.UserRecord.Type.Implements(unmarshallTextType) {
		val, _ := c.UserRecord.Type.MethodByName("UnmarshalText")
		c.UserRecord.UnmarshallMethod = &val
	}

	return &c
}

func ValidateMainEventType(in reflect.Type, AwsEventRecordType reflect.Type) *MainEventType {

	c := MainEventType{Type: in}

	val, ok := c.Type.MethodByName("GetCollection")

	if !ok {
		panic(fmt.Errorf("handler has 3 arguments but the second argument doesn't have a method called 'GetCollection', got: %s", c.Type.Name()))
	}
	if val.Type.NumIn() != 1 || val.Type.NumOut() != 1 {
		panic(fmt.Errorf("the 'GetCollection' method of the second argument should have zero inputs and one return type"))
	}
	collectionType := val.Type.Out(0)

	if !(collectionType.Kind() == reflect.Slice || collectionType.Kind() == reflect.Array) {
		panic(fmt.Errorf("expected an array or slice in the return type of the 'GetCollection' method, got: %s", collectionType.Kind()))
	}
	c.CastFunction = CreateCastWrapper(collectionType.Elem(), AwsEventRecordType)
	c.GetCollection = &val

	return &c
}

func CompileChannelHandler(handlerFunc interface{}) func(context.Context, []byte) (func() (interface{}, error), error) {
	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func {
		panic(fmt.Errorf("expected a function, got: %s", handlerType.Kind()))
	}
	if !(handlerType.NumIn() == 2 || handlerType.NumIn() == 3) {
		panic(fmt.Errorf("expected a handler with 2 or 3 input arguments"))
	}

	// the last input argument is always the channel were the individual records are sent to the handler.
	inputArguments := ChannelHandler{
		ChannelType: *ValidateChannel(handlerType.In(handlerType.NumIn() - 1)),
	}

	if handlerType.NumIn() == 3 {
		inputArguments.MainEventType = *ValidateMainEventType(handlerType.In(1), inputArguments.ChannelEnvelopeType.AwsEventRecordType.Type)
	} else if handlerType.NumIn() == 2 {
		inputArguments.MainEventType.Type = reflect.SliceOf(inputArguments.ChannelType.ChannelEnvelopeType.AwsEventRecordType.Type)
	}

	inputArguments.Values.ActualChannel = reflect.MakeChan(inputArguments.ChannelType.Type, 0)
	inputArguments.Values.DlqChannel = reflect.MakeChan(reflect.ChanOf(reflect.BothDir, inputArguments.ChannelEnvelopeType.AwsEventRecordType.Type), 0)

	newHandler := func(ctx context.Context, payload []byte) (func() (interface{}, error), error) {
		clientHandlerArguments := []reflect.Value{
			reflect.ValueOf(ctx),
		}

		event := reflect.New(inputArguments.MainEventType.Type)
		err := json.Unmarshal(payload, event.Interface())
		if err != nil {
			return nil, err
		}

		var collection reflect.Value

		if inputArguments.MainEventType.GetCollection != nil {
			collection = inputArguments.MainEventType.GetCollection.Func.Call([]reflect.Value{event.Elem()})[0]
			clientHandlerArguments = append(clientHandlerArguments, event.Elem())
			//}
		} else {
			collection = event.Elem()
		}

		go UnmarshallAndSend(ctx, collection, false, inputArguments.ChannelType.ChannelEnvelopeType.AwsEventRecordType.GetMessageMethod, inputArguments.ChannelType.ChannelEnvelopeType.UserRecord.Type, inputArguments.ChannelType.ChannelEnvelopeType.UserRecord.UnmarshallMethod, inputArguments.ChannelType.ChannelEnvelopeType.Type, inputArguments.Values.ActualChannel, inputArguments.MainEventType.CastFunction)

		clientHandlerArguments = append(clientHandlerArguments, inputArguments.Values.ActualChannel)
		x, y := InvokeClientLambda(reflect.ValueOf(handlerFunc), clientHandlerArguments)

		return func() (interface{}, error) {
			return x, y
		}, nil
	}
	return newHandler
}

func CreateCastWrapper(inner reflect.Type, outer reflect.Type) func(reflect.Value) reflect.Value {
	if !outer.Field(0).Type.AssignableTo(inner) {
		panic(fmt.Errorf("trying to cast the event record from the 'GetCollection' method to the event record form the channel: from %s to %s\n", inner.Name(), outer.Name()))
	}

	return func(innerValue reflect.Value) reflect.Value {
		outerEvent := reflect.New(outer)
		outerEvent.Elem().Field(0).Set(innerValue)
		return outerEvent
	}
}

func UnmarshallAndSend(ctx context.Context, collection reflect.Value, chanEventWasPtr bool, getMessageMethod reflect.Method, messageType reflect.Type, unmarshallText *reflect.Method, channelType reflect.Type, channel reflect.Value, castFunction func(reflect.Value) reflect.Value) {
	for i := 0; i < collection.Len(); i++ {

		fmt.Println("unmarshalling a message and sending it on the channel")

		event := collection.Index(i)
		var innerEvent reflect.Value

		if castFunction != nil {
			innerEvent = castFunction(event)
		} else {
			innerEvent = event.Addr()
		}

		channelMessage := reflect.New(channelType).Elem()

		val, dlqErr := UnmarshalInnerMessage(ctx, innerEvent, chanEventWasPtr, getMessageMethod, messageType, unmarshallText)

		// @todo add dlq error handling
		if dlqErr.error != nil {
			fmt.Println(dlqErr.error)
			continue
		}

		channelMessage.Field(0).Set(innerEvent.Elem())
		channelMessage.Field(1).Set(val.Elem())

		channel.Send(channelMessage)
	}

	fmt.Println("closing the channel")
	channel.Close()
}
