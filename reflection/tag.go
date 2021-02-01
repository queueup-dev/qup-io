package reflection

import (
	"errors"
	"reflect"
)

func validStruct(typeReflection reflect.Type) (reflect.Type, error) {
	if typeReflection.Kind() == reflect.Ptr {
		return validStruct(typeReflection.Elem())
	}

	if typeReflection.Kind() != reflect.Struct {
		return typeReflection, errors.New("supplied argument is not a structure")
	}

	return typeReflection, nil
}

func GetTagValue(tag string, fieldName string, typeReflection reflect.Type) (string, error) {

	typeReflection, err := validStruct(typeReflection)

	if err != nil {
		return "", err
	}

	field, ok := typeReflection.FieldByName(fieldName)

	if !ok {
		return "", errors.New("supplied field is not defined in the structure")
	}

	// Don't handle unexported fields
	if field.PkgPath != "" {
		return "", errors.New("supplied field is not exported")
	}

	tagValue, ok := field.Tag.Lookup(tag)

	if !ok {
		return "", errors.New("supplied tag is not present on the field")
	}

	return tagValue, nil
}

func GetTagValues(tag string, typeReflection reflect.Type) ([]string, error) {
	values := make([]string, 0)

	typeReflection, err := validStruct(typeReflection)

	if err != nil {
		return values, err
	}

	for i := 0; i < typeReflection.NumField(); i++ {
		fieldType := typeReflection.Field(i)

		// Don't handle unexported fields
		if fieldType.PkgPath != "" {
			continue
		}

		rawTag, ok := fieldType.Tag.Lookup(tag)

		if ok {
			values = append(values, rawTag)
		}
	}

	return values, nil
}

func GetFieldNamesWithTag(tag string, typeReflection reflect.Type) ([]string, error) {
	fields := make([]string, 0)

	typeReflection, err := validStruct(typeReflection)

	if err != nil {
		return fields, err
	}

	for i := 0; i < typeReflection.NumField(); i++ {
		fieldType := typeReflection.Field(i)

		// Don't handle unexported fields
		if fieldType.PkgPath != "" {
			continue
		}

		_, ok := fieldType.Tag.Lookup(tag)

		if ok {
			fields = append(fields, fieldType.Name)
		}
	}

	return fields, nil
}
