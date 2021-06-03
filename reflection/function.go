package reflection

import "reflect"

func PopulateFunction(readerFuncType reflect.Type, readerFuncPtrValue reflect.Value, reflectedFunc func([]reflect.Value) []reflect.Value) {
	newFuncValue := reflect.MakeFunc(readerFuncType, reflectedFunc)
	readerFuncPtrValue.Elem().Set(newFuncValue)
}
