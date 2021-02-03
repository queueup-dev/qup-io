package lambda

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/queueup-dev/qup-io/lambda/helper_functions_aws_events"
	"github.com/queueup-dev/qup-io/reflection"
	"reflect"
)

type (
	handlerInputSingleMessage struct {
		takesContext                           bool
		mainRequest                            reflect.Type
		getMessageFunction                     reflect.Value
		shouldDereferenceWhenCallingGetMessage bool
		secondaryRequest                       reflect.Type
	}
	handlerOutputSingleMessage struct {
		setMessageFunction                     *reflect.Value
		shouldTakeAddressWhenCallingSetMessage bool
	}

	SingleRequest interface {
		GetMessage() []byte
	}
	SingleResponse interface {
		SetMessage([]byte)
	}
)

var (
	singleRequestType  = reflect.TypeOf((*SingleRequest)(nil)).Elem()
	singleResponseType = reflect.TypeOf((*SingleResponse)(nil)).Elem()
)

func NewHandler(handlerFunction interface{}) lambda.Handler {
	if handlerFunction == nil {
		panic(fmt.Errorf("handler is nil"))
	}

	handlerValue := reflect.ValueOf(handlerFunction)
	handlerType := handlerValue.Type()

	if handlerType.Kind() != reflect.Func {
		panic(fmt.Errorf("handler kind %s is not %s", handlerType.Kind(), reflect.Func))
	}

	return singleRecordHandler{
		input:            validateInput(handlerType),
		output:           validateReturn(handlerType),
		originalFunction: handlerValue,
	}
}

func validateInput(handlerType reflect.Type) HandlerInput {
	numIn := handlerType.NumIn()

	if numIn == 0 {
		return handlerInputSingleMessage{}
	}

	var (
		ret                handlerInputSingleMessage
		takesContextOffset int
	)

	if handlerType.In(0).Implements(reflection.ContextType) {
		ret.takesContext = true
		takesContextOffset = 1
	}

	switch n := numIn - takesContextOffset; {
	case n == 0:
		return ret
	case n == 1:
		ret.mainRequest = handlerType.In(takesContextOffset)
		return ret
	case n == 2:
		mainRequest := handlerType.In(takesContextOffset)

		getMessageFunction, ok := helper_functions_aws_events.GetGetMessageFunction(mainRequest)
		shouldDereference := mainRequest.Kind() == reflect.Ptr

		if !ok {
			getMessageMethod, haveTakenAddress, err := reflection.GetMethodAddressSafe(mainRequest, singleRequestType, "GetMessage")
			if err != nil {
				panic(err)
			}
			shouldDereference = !haveTakenAddress
			getMessageFunction = getMessageMethod.Func
		}

		return handlerInputSingleMessage{
			takesContext:                           ret.takesContext,
			mainRequest:                            mainRequest,
			getMessageFunction:                     getMessageFunction,
			shouldDereferenceWhenCallingGetMessage: shouldDereference,
			secondaryRequest:                       handlerType.In(takesContextOffset + 1),
		}
	default:
		panic(fmt.Errorf("to many arguments input the handler function"))
	}
}

func validateReturn(handler reflect.Type) HandlerOutput {
	n := handler.NumOut()

	if n > 0 {
		if err := reflection.IsOfErrorType(handler.Out(n-1), "last variable of the variables returned by the handler function"); err != nil {
			panic(err)
		}
	}

	switch {
	case n > 3:
		panic(fmt.Errorf("handler may not return more than three values"))
	case n == 3:

		mainResponse := handler.Out(0)

		setMessageFunction, ok := helper_functions_aws_events.GetSetMessageFunction(mainResponse)
		shouldTakeAddress := mainResponse.Kind() != reflect.Ptr

		if !ok {
			setMessageMethod, haveTakenAddress, err := reflection.GetMethodAddressSafe(mainResponse, singleResponseType, "SetMessage")
			if err != nil {
				panic(err)
			}
			shouldTakeAddress = haveTakenAddress
			setMessageFunction = setMessageMethod.Func
		}

		return handlerOutputSingleMessage{
			setMessageFunction:                     &setMessageFunction,
			shouldTakeAddressWhenCallingSetMessage: shouldTakeAddress,
		}
	default:
		return handlerOutputSingleMessage{}
	}
}
