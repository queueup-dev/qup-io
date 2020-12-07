package reflection

import (
	"fmt"
	"reflect"
	"strings"
)

type StructOrPointerToStruct struct {
	OriginalType reflect.Type
	IsPointer    bool
}

func IsStructOrPointerToStruct(typ reflect.Type, nameOfVariable string) (int, error) {
	depth := 0
	for {
		switch typ.Kind() {
		case reflect.Struct:
			return depth, nil
		case reflect.Ptr:
			depth++
			typ = typ.Elem()
		default:
			return depth, fmt.Errorf("expected %v to be of type struct or pointer(s) to a struct, got: %v %v", nameOfVariable, strings.Repeat("*", depth), typ.String())
		}
	}
}

func WalkPointer(typ reflect.Value) reflect.Value {
	for {
		if typ.Kind() != reflect.Ptr {
			break
		}

		if typ.IsNil() {
			if typ.Elem().Kind() == reflect.Ptr {
				typ.Set(reflect.Zero(typ.Type()))
			}
			typ.Set(reflect.New(typ.Type().Elem()))
		}
		typ = typ.Elem()
	}

	return typ
}

func GetAddressOfStruct(val reflect.Value) reflect.Value {
	// To get the address of a value of kind struct we can't just take reflect.PtrTo( type ),
	// as those values are not addressable.
	ptrMainEvent := reflect.New(val.Type())
	ptrMainEvent.Elem().Set(val)
	return ptrMainEvent
}

func GetMethodAddressSafe(typ reflect.Type, interfaceValue reflect.Type, methodName string) (reflect.Method, bool, error) {
	shouldTakePointer := false
	if !typ.Implements(interfaceValue) {
		typ = reflect.PtrTo(typ)
		if !typ.Implements(interfaceValue) {
			return reflect.Method{}, false, fmt.Errorf("expected type %v to implement %v", typ.String(), interfaceValue.Name())
		}
		shouldTakePointer = true
	}
	setMessageMethod, ok := typ.MethodByName(methodName)
	if !ok {
		return reflect.Method{}, false, fmt.Errorf("method %s of interface %v can't be found", methodName, interfaceValue.Name())
	}

	return setMessageMethod, shouldTakePointer, nil
}
