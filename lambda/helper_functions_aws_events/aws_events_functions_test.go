package helper_functions_aws_events

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetMessageFunctions(t *testing.T) {
	for typ, function := range getMessageFunctions {
		t.Run(typ.String(), func(t *testing.T) {

			if typ.Kind() == reflect.Ptr {
				t.Fail()
			}

			val := reflect.New(typ)

			funcVal := reflect.ValueOf(function)

			ret := funcVal.Call([]reflect.Value{val})

			if len(ret) != 1 {
				fmt.Println("expected a one return argument")
				t.Fail()
			}

			if _, ok := ret[0].Interface().([]byte); !ok {
				fmt.Println("expected the return type to be []byte")
				t.Fail()
			}
		})
	}
}

func TestSetMessageFunctions(t *testing.T) {
	for typ, function := range setMessageFunctions {
		t.Run(typ.String(), func(t *testing.T) {

			if typ.Kind() == reflect.Ptr {
				t.Fail()
			}

			val := reflect.New(typ)
			payload := reflect.ValueOf([]byte("hello world"))

			funcVal := reflect.ValueOf(function)

			ret := funcVal.Call([]reflect.Value{val, payload})

			if len(ret) != 0 {
				fmt.Println("expected no return arguments")
				t.Fail()
			}
		})
	}
}
