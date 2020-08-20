package csv

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
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
	if field.Kind() == reflect.Ptr {
		if omitEmpty && value == "" {
			return nil
		}
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	switch field.Interface().(type) {
	case string:
		field.SetString(value)
	case bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case int, int8, int16, int32, int64:
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(i))
	case uint, uint8, uint16, uint32, uint64:
		ui, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return err
		}
		field.SetUint(ui)
	case float32, float64:
		f, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	}
	return nil
}
