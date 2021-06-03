package reflection

import (
	"fmt"
	"reflect"
	"strconv"
)

var (
	numeric = []reflect.Kind{
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Float32,
		reflect.Float64,
	}

	iteratable = []reflect.Kind{
		reflect.Array,
		reflect.Map,
	}
)

type StringSerializable interface {
	ToString() string
}

func IsIteratable(input interface{}) bool {
	for _, kind := range iteratable {
		if reflect.ValueOf(input).Kind() == kind {
			return true
		}
	}

	return false
}

func IsNumeric(input interface{}) bool {
	for _, kind := range numeric {
		if reflect.ValueOf(input).Kind() == kind {
			return true
		}
	}

	return false
}

func IntegerOf(input interface{}) int {
	reflectedValue := reflect.ValueOf(input)
	switch reflectedValue.Kind() {
	case reflect.Uint:
		return int(input.(uint))
	case reflect.Uint8:
		return int(input.(uint8))
	case reflect.Uint16:
		return int(input.(uint16))
	case reflect.Uint32:
		return int(input.(uint32))
	case reflect.Uint64:
		return int(input.(uint64))
	case reflect.Int:
		return input.(int)
	case reflect.Int8:
		return int(input.(int8))
	case reflect.Int16:
		return int(input.(int16))
	case reflect.Int32:
		return int(input.(int32))
	case reflect.Int64:
		return int(input.(int64))
	}

	panic("only integer types are supported")
}

func StringValueOf(input interface{}) string {
	reflectedValue := reflect.ValueOf(input)
	switch reflectedValue.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Itoa(IntegerOf(input))
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", input.(float64))
	case reflect.String:
		return input.(string)
	case reflect.Bool:
		if input.(bool) {
			return "true"
		}
		return "false"
	case reflect.Struct:
		s, ok := input.(StringSerializable)

		if ok {
			return s.ToString()
		}
	}

	panic("only primitives and serializable structures are supported")
}
