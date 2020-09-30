package envvar

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	reflectTag         = "env"
	reflectSeparator   = ","
	reflectRequiredTag = "required"
)

func ToStruct(input interface{}) error {
	typeReflection := reflect.TypeOf(input)
	valueReflection := reflect.ValueOf(input)

	if typeReflection.Kind() == reflect.Ptr {
		typeReflection = typeReflection.Elem()
	}

	if typeReflection.Kind() != reflect.Struct {
		return errors.New("supplied argument is not a structure")
	}

	for i := 0; i < typeReflection.NumField(); i++ {
		fieldType := typeReflection.Field(i)
		rawTag, ok := fieldType.Tag.Lookup(reflectTag)

		// No tag found, continue
		if !ok {
			continue
		}

		parsedTag := strings.Split(rawTag, reflectSeparator)
		tag := parsedTag[0]
		required := len(parsedTag) > 1 && parsedTag[1] == reflectRequiredTag

		value, err := LookupEnv(tag)

		if err != nil {
			if required {
				return fmt.Errorf("required tag %s is not set in your environment variables", tag)
			}
			continue
		}

		fieldValue := valueReflection.Elem().FieldByName(fieldType.Name)

		if !fieldValue.CanSet() {
			if required {
				return fmt.Errorf("required tag %s can not be set in your structure", tag)
			}
			continue
		}

		if fieldValue.Kind() != reflect.String {
			return errors.New("only structure values of type string can be used for env variable mapping")
		}

		fieldValue.SetString(value)
	}

	return nil
}
