package reflection

import (
	"fmt"
	"reflect"
	"strconv"
)

func GetFieldValueFromIndexChain(structField reflect.Value, index []int) reflect.Value {
	for {
		if structField.Kind() != reflect.Ptr {
			break
		}

		if structField.IsNil() {
			return structField
		}

		structField = structField.Elem()
	}

	if len(index) == 0 {
		return structField
	}

	// because pointers can be nil need to recurse one index at a time and perform nil check
	if len(index) > 1 {
		nextField := structField.Field(index[0])
		return GetFieldValueFromIndexChain(nextField, index[1:])
	}

	return GetFieldValueFromIndexChain(structField.FieldByIndex(index), nil)
}

func GetFieldValueFromIndexChainAsString(field reflect.Value, index []int) (string, error) {

	value := GetFieldValueFromIndexChain(field, index)

	if value.IsNil() {
		return "", nil
	}

	return GetFieldValueAsString(field)
}

func GetFieldValueAsString(field reflect.Value) (string, error) {
	switch field.Kind() {
	case reflect.Interface, reflect.Ptr:
		if field.IsNil() {
			return "", nil
		}
		return GetFieldValueAsString(field.Elem())
	case reflect.String:
		return field.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(field.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(field.Float(), 'f', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("unexpected field kind, got: %v", field.Kind())
	}
}
