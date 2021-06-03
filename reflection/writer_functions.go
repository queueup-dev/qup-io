package reflection

import (
	"fmt"
	"reflect"
)

func ValidateWriterFunction(writerFuncPtrValue reflect.Value, baseType reflect.Type) reflect.Type {
	if writerFuncPtrValue.Kind() != reflect.Ptr || writerFuncPtrValue.Elem().Kind() != reflect.Func {
		panic(fmt.Errorf("reader function value has to be passed by pointer"))
	}

	writerFuncType := writerFuncPtrValue.Elem().Type()

	if writerFuncType.NumIn() != 1 {
		panic(fmt.Errorf("expected 1 input argument for the baseRecord function, got: %v", writerFuncType.NumIn()))
	}

	baseReturnType := writerFuncType.In(0)

	isAssignable := baseType.AssignableTo(baseReturnType)

	if !isAssignable {
		panic(fmt.Errorf("the return type '%v' of the writer function is incompatible with the base type '%v", baseReturnType.String(), baseType.String()))
	}

	return writerFuncType
}
