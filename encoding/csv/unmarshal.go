package csv

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"github.com/queueup-dev/qup-io/reflection"
	"reflect"
)

func (m HeaderMapping) Unmarshal(input []string, out interface{}) error {
	ptrValue := reflect.ValueOf(out)

	if ptrValue.Kind() != reflect.Ptr || ptrValue.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(out)}
	}

	return m.UnmarshalValue(input, ptrValue)
}

func (m HeaderMapping) UnmarshalValue(input []string, value reflect.Value) error {
	typ := value.Elem().Type()
	if typ != m.OriginalStruct {
		return errors.New("output struct is not registered")
	}

	for k, field := range m.HeaderMapping {
		content := input[k]
		if err := setFieldFromIndexChain(value.Elem(), false, field.indexChain, content, field.omitempty); err != nil { // Set field of struct
			parseError := csv.ParseError{
				Column: k + 1,
				Err:    err,
			}
			return &parseError
		}
	}

	return nil
}

func setFieldFromIndexChain(fieldPtr reflect.Value, isPointer bool, index []int, value string, omitEmpty bool) error {

	field := fieldPtr

	if isPointer {
		// initialize nil pointer
		if field.IsNil() {
			setField(field, "", omitEmpty)
		}
		field = fieldPtr.Elem()
	}

	// because pointers can be nil need to recurse one index at a time and perform nil check
	if len(index) > 1 {
		nextField := field.Field(index[0])
		return setFieldFromIndexChain(nextField, nextField.Kind() == reflect.Ptr, index[1:], value, omitEmpty)
	}
	return setField(field.FieldByIndex(index), value, omitEmpty)
}

func setField(field reflect.Value, value string, omitEmpty bool) error {
	return reflection.PopulateFromString(field, value, omitEmpty)
}
