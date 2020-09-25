package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

// 1. eventType is collection   <- getMessage on the inner-event-type, gives []byte
// 2. messageType is collection <- getMessage on the inner-event-type, gives function that for each message inner type gives a []byte
// 3. neither is a collection   <- getMessage on the event-type = inner-event-type, gives []byte

//  - invoke the user-handler function each time <-- best use case for 1
//  - provide a read channel in the message argument <-- best use case for 2.
//  - as is <-- best use case for 3.
//  - forget the event context

// Three kinds of casting:
// A. cast event-type to array and give inner types in both argument, call the user-handler in a for loop
// B. cast body-type  to array and give inner types in both argument, call the user-handler in a for loop
// C. cast as is
// D. cast event-type to array and go-routine the whole function
// E. cast body-type to array and return

// Applicable casting: A. B. C.
// -  event,  body
// - *event,  body
// -  event, *body
// - *event, *body

// - []event,   body    <-- only A. applicable
// - []event,  *body    <-- only A. applicable
// -   event, []body	<-- only C. applicable, but should be discouraged
// -  *event, []body	<-- only C. applicable, but should be discouraged

// This compiles the handler using the straight forward type definition of the second argument.
func CompileHandler(handlerFunc interface{}, eventIsCollection bool) func(context.Context, []byte) (func() (interface{}, error), error) {

	handler := reflect.ValueOf(handlerFunc)
	handlerType := reflect.TypeOf(handlerFunc)

	if handlerType.NumIn() != 3 {
		panic(fmt.Errorf("something went wrong with casting handler arguments"))
		//return interface{}(lambda.NewHandler(handlerFunc)).(func(context.Context, []byte) (interface{}, error))
	}

	eventType := handlerType.In(1)
	eventInnerType := eventType
	outerEventType := eventType
	eventWasPtr := false

	if eventIsCollection {
		outerEventType = reflect.SliceOf(eventType)
	}

	if eventInnerType.Kind() == reflect.Ptr {
		eventInnerType = eventType.Elem()
		eventWasPtr = true
	}

	if !eventInnerType.Implements(getMessageType) {
		panic(fmt.Errorf("handler accepts three arguments, but the second arguments doesn't implement %s, got%s", getMessageType.Name(), eventType.Kind()))
	}

	getMessageMethod, _ := eventInnerType.MethodByName("GetMessage")

	messageType := handlerType.In(2)
	innerMessageType := messageType
	messageWasPtr := false
	messageWasChan := false

	if messageType.Kind() == reflect.Chan {
		innerMessageType = innerMessageType.Elem()
		messageWasChan = true
	}

	if messageType.Kind() == reflect.Ptr {
		innerMessageType = innerMessageType.Elem()
		messageWasPtr = true
	}

	fmt.Println(messageWasChan, messageWasPtr)

	var unmarshallText *reflect.Method

	if messageType.Implements(unmarshallTextType) {
		val, _ := messageType.MethodByName("UnmarshalText")
		unmarshallText = &val
	}

	newHandler := func(ctx context.Context, payload []byte) (func() (interface{}, error), error) {
		event := reflect.New(outerEventType)
		err := json.Unmarshal(payload, event.Interface())

		if err != nil {
			return nil, err
		}

		if eventIsCollection {

			var (
				ret []interface{}
				dlq []struct {
					input interface{}
					error
				}
			)

			for i := 0; i < event.Len(); i++ {

				innerEvent := event.Index(i)

				val, dlqErr := UnmarshalInnerMessage(ctx, innerEvent, eventWasPtr, getMessageMethod, messageType, unmarshallText)

				if dlqErr.error != nil {
					dlq = append(dlq, dlqErr)
				}

				args := []reflect.Value{
					reflect.ValueOf(ctx),
					event.Elem(),
					val.Elem(),
				}

				res, err := InvokeClientLambda(handler, args)

				if err != nil {
					ret = append(ret, res)
				} else {
					dlqErr.error = err
					dlq = append(dlq, dlqErr)
				}
			}

			// @todo dlq handling

			return func() (interface{}, error) {
				return ret, nil
			}, nil
		} else {
			return func() (interface{}, error) {
				mes, dlqError := UnmarshalInnerMessage(ctx, event, eventWasPtr, getMessageMethod, messageType, unmarshallText)
				if dlqError.error != nil {
					return nil, dlqError.error
				}

				args := []reflect.Value{
					reflect.ValueOf(ctx),
					event.Elem(),
					mes.Elem(),
				}

				return InvokeClientLambda(handler, args)
			}, nil
		}
	}
	return newHandler
}

func UnmarshalInnerMessage(ctx context.Context, event reflect.Value, eventWasPtr bool, getMessage reflect.Method, messageType reflect.Type, unmarshallText *reflect.Method) (*reflect.Value, struct {
	input interface{}
	error
}) {
	var (
		err        error
		rawMessage reflect.Value
		message    = reflect.New(messageType)
	)

	if eventWasPtr {
		rawMessage = getMessage.Func.Call([]reflect.Value{event.Elem().Elem()})[0]
	} else {
		rawMessage = getMessage.Func.Call([]reflect.Value{event.Elem()})[0]
	}

	if unmarshallText != nil {
		err = (*unmarshallText).Func.Call([]reflect.Value{message, rawMessage})[0].Interface().(error)
		//message.Interface().(encoding.TextUnmarshaler).UnmarshalText(rawMessage)
	} else {
		err = json.Unmarshal(rawMessage.Interface().([]byte), message.Interface())
	}

	if err != nil {
		return nil, struct {
			input interface{}
			error
		}{
			rawMessage.Interface(),
			err,
		}
	}

	return &message, struct {
		input interface{}
		error
	}{
		rawMessage.Interface(),
		nil,
	}
}

func InvokeClientLambda(handler reflect.Value, args []reflect.Value) (interface{}, error) {
	var err error

	response := handler.Call(args)

	if errVal, ok := response[1].Interface().(error); ok {
		err = errVal
	}

	return response[0].Interface(), err
}
