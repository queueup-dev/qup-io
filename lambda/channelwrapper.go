package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

//var channel reflect.Value

func CompileChannelHandler(handlerFunc interface{}) func(context.Context, []byte) (func() (interface{}, error), error) {

	handler := reflect.ValueOf(handlerFunc)
	handlerType := reflect.TypeOf(handlerFunc)

	if handlerType.NumIn() < 2 {
		panic(fmt.Errorf("something went wrong with casting handler arguments"))
		//return interface{}(lambda.NewHandler(handlerFunc)).(func(context.Context, []byte) (interface{}, error))
	}

	userChannelType := handlerType.In(handlerType.NumIn() - 1)

	if userChannelType.Kind() != reflect.Chan {
		panic(fmt.Errorf("expected the last argument of the handler to be a channel, got: %s", userChannelType.Name()))
	}

	if userChannelType.Elem().Kind() != reflect.Struct && userChannelType.Elem().NumField() != 2 {
		panic(fmt.Errorf("expected the inner type of the channel to be a struct of 2 elements, got: %s", userChannelType.Elem().Name()))
	}

	chanEventType := userChannelType.Elem().Field(0).Type
	chanEventInnerType := chanEventType
	chanEventWasPtr := false

	if chanEventInnerType.Kind() == reflect.Ptr {
		chanEventInnerType = chanEventInnerType.Elem()
		chanEventWasPtr = true
	}

	if !chanEventInnerType.Implements(getMessageType) {
		panic(fmt.Errorf("the first field of the channel struct doesn't implement %s, got: %s", getMessageType.Name(), chanEventInnerType.Name()))
	}

	getMessageMethod, _ := chanEventInnerType.MethodByName("GetMessage")

	messageType := userChannelType.Elem().Field(1).Type
	innerMessageType := messageType
	//messageWasPtr := false

	if messageType.Kind() == reflect.Ptr {
		innerMessageType = innerMessageType.Elem()
		//messageWasPtr = true
	}

	var unmarshallText *reflect.Method

	if messageType.Implements(unmarshallTextType) {
		val, _ := messageType.MethodByName("UnmarshalText")
		unmarshallText = &val
	}

	var (
		chanEventOuterType  reflect.Type
		getCollectionMethod *reflect.Method
	)

	if handlerType.NumIn() == 3 {
		chanEventOuterType = handlerType.In(1)
		if !chanEventOuterType.Implements(getCollectionType) {
			panic(fmt.Errorf("handler has 3 arguments but the second argument doesn't implement %s, got: %s", getCollectionType.Name(), chanEventOuterType.Name()))
		}

		val, _ := chanEventOuterType.MethodByName(getCollectionType.Method(0).Name)
		getCollectionMethod = &val

		//collectionType := getCollectionMethod.Type.Out(0)

		// @todo have to call the function in order to check the underlying type.
		//if collectionType.Kind() != reflect.Slice || collectionType.Kind() != reflect.Array {
		//	panic(fmt.Errorf("expected an array or slice in the return type of %s, got: %s", getCollectionType.Method(0).Name, collectionType.Name()))
		//}

		//if !chanEventType.AssignableTo(getCollectionMethod.Type.Out(0).Elem()) {
		//	panic(fmt.Errorf("getcollection and getmessage not comparable"))
		//}
	} else if handlerType.NumIn() == 2 {
		chanEventOuterType = reflect.SliceOf(chanEventType)
	} else {
		panic(fmt.Errorf("expected a handler with 2 or 3 arguments"))
	}

	//channel = makeChannel(userChannelType.Elem(), reflect.BothDir, 1)
	channel := reflect.MakeChan(userChannelType, 0)

	newHandler := func(ctx context.Context, payload []byte) (func() (interface{}, error), error) {

		// reflect.MakeChan(ctype, buffer)
		event := reflect.New(chanEventOuterType)

		err := json.Unmarshal(payload, event.Interface())

		if err != nil {
			return nil, err
		}

		var collection reflect.Value

		clientHandlerArguments := []reflect.Value{
			reflect.ValueOf(ctx),
		}

		if getCollectionMethod != nil {
			//if eventWasPtr {
			//	collection = getCollectionMethod.Func.Call([]reflect.Value{event.Elem().Elem()})[0]
			//} else {
			collection = getCollectionMethod.Func.Call([]reflect.Value{event.Elem()})[0]
			clientHandlerArguments = append(clientHandlerArguments, event.Elem())
			//}
		} else {
			collection = event
		}

		go UnmarshallAndSend(ctx, collection, chanEventWasPtr, getMessageMethod, messageType, unmarshallText, userChannelType.Elem(), channel)

		clientHandlerArguments = append(clientHandlerArguments, channel)

		x, y := InvokeClientLambda(handler, clientHandlerArguments)

		return func() (interface{}, error) {
			return x, y
		}, nil
	}
	return newHandler
}

func UnmarshallAndSend(ctx context.Context, collection reflect.Value, chanEventWasPtr bool, getMessageMethod reflect.Method, messageType reflect.Type, unmarshallText *reflect.Method, channelType reflect.Type, channel reflect.Value) {
	for i := 0; i < collection.Elem().Len(); i++ {

		fmt.Println("unmarshalling a message and sending it on the channel")

		innerEvent := collection.Elem().Index(i)
		channelMessage := reflect.New(channelType).Elem()

		val, dlqErr := UnmarshalInnerMessage(ctx, innerEvent, false, getMessageMethod, messageType, unmarshallText)

		// @todo add dlq error handling
		if dlqErr.error != nil {
			fmt.Println(dlqErr.error)
			continue
		}

		channelMessage.Field(0).Set(innerEvent)
		channelMessage.Field(1).Set(val.Elem())

		channel.Send(channelMessage)
	}

	fmt.Println("closing the channel")
	channel.Close()
}

func makeChannel(t reflect.Type, chanDir reflect.ChanDir, buffer int) reflect.Value {
	ctype := reflect.ChanOf(chanDir, t)
	return reflect.MakeChan(ctype, buffer)
}
