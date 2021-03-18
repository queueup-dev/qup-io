package reflection

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func SetFieldFromIndexChain(field reflect.Value, index []int, value string, fieldInfo StructFieldInfo) error {
	field = WalkPointer(field, false)

	if field.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected a struct, got: %v", field.String()))
	}

	// because pointers can be nil need to recurse one index at a time and perform nil check
	if len(index) > 1 {
		nextField := field.Field(index[0])
		return SetFieldFromIndexChain(nextField, index[1:], value, fieldInfo)
	}

	if field.Kind() == reflect.Ptr {
		if fieldInfo.OmitEmpty && value == "" {
			return nil
		}

		if field.Type().Implements(TextEncoderType) {
			return PopulateFromCustomUnmarshaller(field.FieldByIndex(index), value)
		}

		field = WalkPointer(field, false)
	}

	if fieldInfo.JsonMarshalling {
		return PopulateFromJsonUnmarshaller(field.FieldByIndex(index), value, fieldInfo.OmitEmpty)
	}

	return PopulateFromString(field.FieldByIndex(index), value, fieldInfo.OmitEmpty)
}

func PopulateFromString(field reflect.Value, value string, omitEmpty bool) error {
	if field.Kind() == reflect.Ptr {
		if omitEmpty && value == "" {
			return nil
		}

		field = WalkPointer(field, omitEmpty)
	}

	if value == "" {
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Int, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(int64(i))
	case reflect.Int8:
		i, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return err
		}
		field.SetInt(int64(i))
	case reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return err
		}
		field.SetInt(int64(i))
	case reflect.Int32:
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(int64(i))
	case reflect.Uint, reflect.Uint64:
		ui, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return err
		}
		field.SetUint(ui)
	case reflect.Uint8:
		ui, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return err
		}
		field.SetUint(ui)
	case reflect.Uint16:
		ui, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return err
		}
		field.SetUint(ui)
	case reflect.Uint32:
		ui, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetUint(ui)
	case reflect.Float32:
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case reflect.Struct:
		if field.Type().Implements(TextEncoderType) {
			errorValue := field.MethodByName("UnmarshalJSON").Call([]reflect.Value{reflect.ValueOf([]byte(value))})[0]
			return errorValue.Interface().(error)
		}
		return fmt.Errorf("unsupported field type, got struct without custom text decoding method provided")
	case reflect.Slice:
		arrayValue := reflect.New(field.Type())

		err := json.Unmarshal([]byte(value), arrayValue.Interface())
		if err != nil {
			return err
		}
		field.Set(arrayValue.Elem())
	default:
		return fmt.Errorf("unsupported field type, got: %v", field.Kind())
	}

	return nil
}

func PopulateFromCustomUnmarshaller(field reflect.Value, value string) error {
	if value == "" {
		return nil
	}

	if field.IsNil() {
		if field.Elem().Kind() == reflect.Ptr {
			field.Set(reflect.Zero(field.Type()))
		}
		field.Set(reflect.New(field.Type().Elem()))
	}

	method := field.MethodByName("UnmarshalTEXT")
	errorValue := method.Call([]reflect.Value{reflect.ValueOf([]byte(value))})[0]

	if errorValue.IsNil() {
		return nil
	}

	return errorValue.Interface().(error)
}

func PopulateFromJsonUnmarshaller(field reflect.Value, value string, omitEmpty bool) error {
	if value == "" && omitEmpty {
		return nil
	}

	if field.Kind() != reflect.Struct && field.IsNil() {
		if field.Elem().Kind() == reflect.Ptr {
			field.Set(reflect.Zero(field.Type()))
		}
		field.Set(reflect.New(field.Type().Elem()))
	}

	if value == "" {
		return nil
	}

	return json.Unmarshal([]byte(value), field.Interface())
}
