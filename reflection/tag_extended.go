package reflection

import (
	"fmt"
	"reflect"
	"strings"
)

type StructFieldInfo struct {
	Type              reflect.Type
	OmitEmpty         bool
	Required          bool
	IsArray           bool
	CustomMarshalling bool
	JsonMarshalling   bool
	Name              string
	IndexChain        []int
}

type TextEncoder interface {
	MarshalTEXT() ([]byte, error)
	UnmarshalTEXT(data []byte) error
}

var (
	TextEncoderType    = reflect.TypeOf((*TextEncoder)(nil)).Elem()
	SkipEmbeddedFields = false
	TagSeparator       = ","
)

type TagInfoWriter interface {
	WriteField(structField reflect.StructField, fieldInfo StructFieldInfo, index int)
	WriteStruct(structField reflect.StructField, fieldInfo StructFieldInfo, index int) TagInfoWriter
	// WriteArray(structField reflect.StructField, indexChain []int)
}

// Get's the field info's from struct and referenced structs inside of it.
// When a cyclic reference is encountered the function panics.
func GetFieldInfo(structType reflect.Type, tag string, writeFieldInfo TagInfoWriter) {
	getFieldInfoRecursive(structType, tag, []int{}, map[reflect.Type]bool{}, writeFieldInfo)
}

func getFieldInfoRecursive(typ reflect.Type, tagName string, parentIndexChain []int, parentChainSet map[reflect.Type]bool, tagWriter TagInfoWriter) {
	defer func() {
		// When we are finished with the current type, we remove it from the chain set.
		parentChainSet[typ] = false
	}()

	for {
		if typ.Kind() != reflect.Ptr {
			break
		}
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected a struct type, instead got: %v", typ.Kind()))
	}

	if parentChainSet[typ] {
		panic(fmt.Errorf("cyclic dependency found in target struct: '%v'", typ.Name()))
	}

	parentChainSet[typ] = true
	fieldsCount := typ.NumField()

fieldLoop:
	for i := 0; i < fieldsCount; i++ {
		field := typ.Field(i)

		if field.PkgPath != "" || (SkipEmbeddedFields && field.Anonymous) {
			continue fieldLoop
		}

		fieldTag, _ := field.Tag.Lookup(tagName)

		if fieldTag == "-" {
			continue fieldLoop
		}

		fieldInfo := StructFieldInfo{}

		var cpy = make([]int, len(parentIndexChain))
		copy(cpy, parentIndexChain)
		fieldInfo.IndexChain = append(cpy, i)
		fieldInfo.Type = field.Type

		tagArray := strings.Split(fieldTag, TagSeparator)

		if tagArray[0] != "" {
			fieldInfo.Name = tagArray[0]
		} else {
			fieldInfo.Name = field.Name
		}

		for i, tag := range tagArray {
			if i == 0 {
				continue
			}
			switch tag {
			case "json":
				fieldInfo.JsonMarshalling = true
			case "omitempty":
				fieldInfo.OmitEmpty = true
			case "required":
				fieldInfo.Required = true
			}
		}

		if fieldInfo.JsonMarshalling {
			tagWriter.WriteField(field, fieldInfo, i)
			continue fieldLoop
		}

		// if the field is a pointer to a struct, follow the pointer then create field info for each field
		for {
			// if the field implements JsonEncoding we consider it as merely a field
			if fieldInfo.Type.Implements(TextEncoderType) {
				fieldInfo.CustomMarshalling = true
				tagWriter.WriteField(field, fieldInfo, i)
				continue fieldLoop
			}
			if fieldInfo.Type.Kind() == reflect.Array || fieldInfo.Type.Kind() == reflect.Slice {
				fieldInfo.IsArray = true
				tagWriter.WriteField(field, fieldInfo, i)
				continue fieldLoop
			}

			if fieldInfo.Type.Kind() != reflect.Ptr {
				break
			}
			fieldInfo.Type = fieldInfo.Type.Elem()
		}

		// if the field is a struct, that doesn't implement JsonEncoding, create a fieldInfo for each of its fields
		if fieldInfo.Type.Kind() == reflect.Struct {
			writer := tagWriter.WriteStruct(field, fieldInfo, i)
			getFieldInfoRecursive(fieldInfo.Type, tagName, fieldInfo.IndexChain, parentChainSet, writer)
			continue fieldLoop
		}

		tagWriter.WriteField(field, fieldInfo, i)
	}
}
