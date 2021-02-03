package helper_functions_aws_events

// Helper functions for aws events.
// As we can't extend the aws types with methods, we store the helper functions in a map indexed by type.

import (
	"github.com/aws/aws-lambda-go/events"
	"reflect"
)

var (
	GetSetMessageFunction func(typ reflect.Type) (reflect.Value, bool) = getFunctionFromMap(setMessageFunctions)
	GetGetMessageFunction func(typ reflect.Type) (reflect.Value, bool) = getFunctionFromMap(getMessageFunctions)
)

var getMessageFunctions = map[reflect.Type]interface{}{
	reflect.TypeOf(events.APIGatewayProxyRequest{}): func(request *events.APIGatewayProxyRequest) []byte {
		return []byte(request.Body)
	},
}

var setMessageFunctions = map[reflect.Type]interface{}{
	reflect.TypeOf(events.APIGatewayProxyResponse{}): func(response *events.APIGatewayProxyResponse, message []byte) {
		response.Body = string(message)
	},
}

func getFunctionFromMap(mp map[reflect.Type]interface{}) func(reflect.Type) (reflect.Value, bool) {
	return func(typ reflect.Type) (reflect.Value, bool) {

		for {
			if typ.Kind() != reflect.Ptr {
				break
			}

			typ = typ.Elem()
		}

		function, ok := mp[typ]

		if !ok {
			return reflect.Value{}, false
		}

		return reflect.ValueOf(function), true
	}
}
