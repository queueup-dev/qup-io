package reflection

import (
	"fmt"
	"reflect"
)

func ValidateReaderFunction(readerFunc reflect.Value, baseType reflect.Type) (reflect.Type, reflect.Type) {
	if readerFunc.Kind() != reflect.Ptr || readerFunc.Elem().Kind() != reflect.Func {
		panic(fmt.Errorf("reader function value has to be passed by pointer"))
	}

	readerFuncType := readerFunc.Elem().Type()

	if readerFuncType.NumIn() != 0 {
		panic(fmt.Errorf("reader function can't have input arguments"))
	}

	if readerFuncType.NumOut() != 3 {
		panic(fmt.Errorf("reader function is expected to have 3 return arguments"))
	}
	if readerFuncType.Out(1).Kind() != reflect.Bool {
		panic(fmt.Errorf("expected bool type as the second output argument of the reader function, got: %v", readerFuncType.Out(1).Kind()))
	}

	err := IsOfErrorType(readerFuncType.Out(2), "the third output argument")
	if err != nil {
		panic(err)
	}

	baseTypFromReader := readerFuncType.Out(0)

	isAssignable := baseType.AssignableTo(baseTypFromReader)

	if !isAssignable {
		panic(fmt.Errorf("the base type '%v' is not assignable to the return type '%v' of the reader function", baseType.String(), baseTypFromReader.String()))
	}

	return readerFuncType, baseTypFromReader
}
