package reflection

import (
	"reflect"
)

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
