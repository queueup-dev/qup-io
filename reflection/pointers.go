package reflection

import (
	"reflect"
)

func WalkPointer(val reflect.Value, omitEmpty bool) reflect.Value {
	for {
		if val.Kind() != reflect.Ptr {
			break
		}

		if val.IsNil() {
			if omitEmpty {
				return val
			} else {
				if val.Elem().Kind() == reflect.Ptr {
					val.Set(reflect.Zero(val.Type()))
				}
				val.Set(reflect.New(val.Type().Elem()))
			}
		}
		val = val.Elem()
	}

	return val
}

func GetAddressOfStruct(val reflect.Value) reflect.Value {
	// To get the address of a value of kind struct we can't just take reflect.PtrTo( type ),
	// as those values are not addressable.
	ptrMainEvent := reflect.New(val.Type())
	ptrMainEvent.Elem().Set(val)
	return ptrMainEvent
}
