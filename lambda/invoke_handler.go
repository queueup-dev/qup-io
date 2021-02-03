package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/queueup-dev/qup-io/reflection"
	"reflect"
)

type (
	singleRecordHandler struct {
		input            HandlerInput
		output           HandlerOutput
		originalFunction reflect.Value
	}
	HandlerInput interface {
		prepareArguments(ctx context.Context, payload []byte) ([]reflect.Value, error)
	}
	HandlerOutput interface {
		processReturns([]reflect.Value) ([]byte, error)
	}
)

func (h singleRecordHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {

	args, err := h.input.prepareArguments(ctx, payload)

	if err != nil {
		return nil, err
	}

	response := h.originalFunction.Call(args)

	return h.output.processReturns(response)
}

func (h handlerInputSingleMessage) prepareArguments(ctx context.Context, payload []byte) ([]reflect.Value, error) {
	var args []reflect.Value
	if h.takesContext {
		args = append(args, reflect.ValueOf(ctx))
	}

	if h.mainRequest != nil {
		event := reflect.New(h.mainRequest)
		if err := json.Unmarshal(payload, event.Interface()); err != nil {
			return nil, err
		}
		args = append(args, event.Elem())

		if h.secondaryRequest != nil {
			if h.shouldDereferenceWhenCallingGetMessage {
				event = event.Elem()
			}

			innerPayload := h.getMessageFunction.Call([]reflect.Value{event})[0]
			secondaryEvent := reflect.New(h.secondaryRequest)

			if err := json.Unmarshal(innerPayload.Interface().([]byte), secondaryEvent.Interface()); err != nil {
				return nil, err
			}
			args = append(args, secondaryEvent.Elem())
		}
	}

	return args, nil
}

func (h handlerOutputSingleMessage) processReturns(response []reflect.Value) ([]byte, error) {
	if len(response) > 0 {
		if err, ok := response[len(response)-1].Interface().(error); ok && err != nil {
			return nil, err
		}
	}

	if len(response) == 1 {
		return nil, nil
	}

	mainResponse := response[0]

	if len(response) == 3 {
		secondaryVal := response[1].Interface()

		secondaryPayload, err := json.Marshal(secondaryVal)

		if err != nil {
			fmt.Println("an error occurred when marshalling the secondary payload")
			return nil, err
		}

		if h.shouldTakeAddressWhenCallingSetMessage {
			mainResponse = reflection.GetAddressOfStruct(mainResponse)
		}

		(*h.setMessageFunction).Call([]reflect.Value{mainResponse, reflect.ValueOf(secondaryPayload)})
	}

	return json.Marshal(mainResponse.Interface())
}
