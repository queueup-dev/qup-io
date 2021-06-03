package reflection

import (
	"context"
	"fmt"
	"reflect"
)

var (
	ErrorType      = reflect.TypeOf((*error)(nil)).Elem()
	ErrorZeroValue = reflect.Zero(ErrorType)

	ContextType = reflect.TypeOf((*context.Context)(nil)).Elem()
)

func IsOfErrorType(typ reflect.Type, nameOfVariable string) error {
	if !typ.Implements(ErrorType) {
		return fmt.Errorf("expected %v to be the error interface, got: %v", nameOfVariable, typ.Name())
	}
	return nil
}
