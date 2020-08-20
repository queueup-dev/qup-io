package csv

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type StructFieldInfo struct {
	omitempty  bool
	required   bool
	name       string
	indexChain []int
}

type RegisteredStruct struct {
	originalStruct          reflect.Type
	HeaderDerivedFromStruct []string
	fieldInfos              []StructFieldInfo
}

func RegisterStruct(baseStruct interface{}) *RegisteredStruct {
	ret := RegisteredStruct{}

	ptrValue := reflect.ValueOf(baseStruct)

	if ptrValue.Kind() != reflect.Ptr || ptrValue.IsNil() {
		panic(json.InvalidUnmarshalError{Type: reflect.TypeOf(baseStruct)})
	}

	value := ptrValue.Elem()
	ret.originalStruct = value.Type()

	if !value.CanInterface() {
		panic(errors.New("Can't interface with given type."))
	}

	//if el.Kind() != reflect.Ptr {
	//	return &json.InvalidUnmarshalError{Type: reflect.TypeOf(el)}
	//}

	var err error

	ret.fieldInfos, err = GetFieldInfo(ret.originalStruct, []int{}, map[reflect.Type]bool{})
	if err != nil {
		panic(err)
	}

	if len(ret.fieldInfos) == 0 {
		panic(errors.New("no csv struct tags found"))
	}

	return &ret
}

func RegisterWriter() {

}

func GetFieldInfo(typ reflect.Type, parentIndexChain []int, parentChainSet map[reflect.Type]bool) ([]StructFieldInfo, error) {
	defer func() {
		// When we are finished with the current type, we remove it from the chain set.
		parentChainSet[typ] = false
	}()

	if parentChainSet[typ] {
		return nil, fmt.Errorf("cyclic dependency found in target struct: '%v'", typ.Name())
	}

	parentChainSet[typ] = true

	fieldsCount := typ.NumField()
	fieldsList := make([]StructFieldInfo, 0, fieldsCount)

	for i := 0; i < fieldsCount; i++ {
		field := typ.Field(i)

		if field.PkgPath != "" || (SkipEmbeddedFields && field.Anonymous) {
			continue
		}

		fieldTag := field.Tag.Get(TagName)

		if fieldTag == "-" {
			continue
		}

		var cpy = make([]int, len(parentIndexChain))
		copy(cpy, parentIndexChain)
		indexChain := append(cpy, i)

		// if the field is a pointer to a struct, follow the pointer then create field info for each field
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			fieldInfo, err := GetFieldInfo(field.Type.Elem(), indexChain, parentChainSet)
			if err != nil {
				return nil, err
			}
			fieldsList = append(fieldsList, fieldInfo...)

			continue
		}
		// if the field is a struct, create a fieldInfo for each of its fields
		if field.Type.Kind() == reflect.Struct {
			fieldInfo, err := GetFieldInfo(field.Type, indexChain, parentChainSet)
			if err != nil {
				return nil, err
			}
			fieldsList = append(fieldsList, fieldInfo...)

			continue
		}

		fieldInfo := StructFieldInfo{indexChain: indexChain}

		fieldTags := strings.Split(fieldTag, TagSeparator)

		if fieldTags[0] != "" {
			fieldInfo.name = fieldTags[0]
		} else {
			fieldInfo.name = field.Name
		}

		if len(fieldTags) > 1 {
			switch fieldTags[1] {
			case "omitempty":
				fieldInfo.omitempty = true
			case "required":
				fieldInfo.required = true
			}
		}

		if len(fieldTags) > 2 {
			switch fieldTags[2] {
			case "omitempty":
				fieldInfo.omitempty = true
			case "required":
				fieldInfo.required = true
			}
		}

		fieldsList = append(fieldsList, fieldInfo)
	}

	return fieldsList, nil
}

func (f StructFieldInfo) MatchesKey(key string) bool {
	if strings.TrimSpace(key) == f.name {
		return true
	}

	return false
}
