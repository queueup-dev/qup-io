package reflection

import (
	"errors"
	"reflect"
)

const (
	reflectSeparator = ","
)

func GetTagValue(tag string, fieldName string, object interface{}) (string, error) {
	typeReflection := reflect.TypeOf(object)

	if typeReflection.Kind() == reflect.Ptr {
		typeReflection = typeReflection.Elem()
	}

	if typeReflection.Kind() != reflect.Struct {
		return "", errors.New("supplied argument is not a structure")
	}

	field, ok := typeReflection.FieldByName(fieldName)

	if !ok {
		return "", errors.New("supplied field is not defined in the structure")
	}

	tagValue, ok := field.Tag.Lookup(tag)

	if !ok {
		return "", errors.New("supplied tag is not present on the field")
	}

	return tagValue, nil
}

func GetTagValues(tag string, object interface{}) ([]string, error) {
	var values = make([]string, 0)
	typeReflection := reflect.TypeOf(object)

	if typeReflection.Kind() == reflect.Ptr {
		typeReflection = typeReflection.Elem()
	}

	if typeReflection.Kind() != reflect.Struct {
		return values, errors.New("supplied argument is not a structure")
	}

	for i := 0; i < typeReflection.NumField(); i++ {
		fieldType := typeReflection.Field(i)
		rawTag, ok := fieldType.Tag.Lookup(tag)

		if ok {
			values = append(values, rawTag)
		}
	}

	return values, nil
}

func GetFieldNamesWithTag(tag string, object interface{}) ([]string, error) {
	var fields []string = make([]string, 0)
	typeReflection := reflect.TypeOf(object)

	if typeReflection.Kind() == reflect.Ptr {
		typeReflection = typeReflection.Elem()
	}

	if typeReflection.Kind() != reflect.Struct {
		return fields, errors.New("supplied argument is not a structure")
	}

	for i := 0; i < typeReflection.NumField(); i++ {
		fieldType := typeReflection.Field(i)
		_, ok := fieldType.Tag.Lookup(tag)

		if ok {
			fields = append(fields, fieldType.Name)
		}
	}

	return fields, nil
}
