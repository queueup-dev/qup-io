package envvar

import (
	"errors"
	"fmt"
	"github.com/queueup-dev/qup-io/reflection"
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

		value, errString := LookupEnv(tag)

		if errString != nil {
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

		err := reflection.PopulateFromString(fieldValue, value, false)
		if err != nil {
			return err
		}
	}

	return nil
}
