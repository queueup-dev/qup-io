package csv

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func (m HeaderMapping) Marshal(input interface{}) ([]string, error) {
	ptrValue := reflect.ValueOf(input)

	if ptrValue.Kind() != reflect.Ptr || ptrValue.IsNil() {
		return nil, &json.InvalidUnmarshalError{Type: reflect.TypeOf(input)}
	}

	return m.MarshalValue(ptrValue.Elem())
}

func (m HeaderMapping) MarshalValue(value reflect.Value) ([]string, error) {
	typ := value.Type()

	if typ != m.OriginalStruct {
		return nil, errors.New("output struct is not registered")
	}

	output := make([]string, len(m.Header))

	for k, field := range m.HeaderMapping {
		stringValue, err := getValueFromIndexChain(value, false, field.indexChain)
		if err != nil {
			return nil, err
		}
		output[k] = stringValue
	}

	return output, nil
}

func getValueFromIndexChain(f reflect.Value, isPointer bool, index []int) (string, error) {
	endField := f
	if isPointer {
		if endField.IsNil() {
			return "", nil
		}
		endField = f.Elem()
	}
	// because pointers can be nil need to recurse one index at a time and perform nil check
	if len(index) > 1 {
		nextField := endField.Field(index[0])
		return getValueFromIndexChain(nextField, nextField.Kind() == reflect.Ptr, index[1:])
	}
	return getFieldAsString(endField.FieldByIndex(index))
}

func getFieldAsString(field reflect.Value) (str string, err error) {
	switch field.Kind() {
	case reflect.Interface, reflect.Ptr:
		if field.IsNil() {
			return "", nil
		}
		return getFieldAsString(field.Elem())
	case reflect.String:
		return field.String(), nil
	default:
		// Check if field is go native type
		switch field.Interface().(type) {
		case string:
			return field.String(), nil
		case bool:
			return strconv.FormatBool(field.Bool()), nil
		case int, int8, int16, int32, int64:
			return strconv.FormatInt(field.Int(), 10), nil
		case uint, uint8, uint16, uint32, uint64:
			return strconv.FormatUint(field.Uint(), 10), nil
		case float32, float64:
			return FormatFloat(field.Float()), nil
		}
	}
	return str, nil
}
