package reflection

import (
	"fmt"
	"reflect"
	"strings"
)

type StructFieldInfo struct {
	Name              string
	Type              reflect.Type
	IndexChain        []int
	OmitEmpty         bool
	Required          bool
	IsArray           bool
	CustomMarshalling bool
	JsonMarshalling   bool
	HasEmptyTag       bool
}

type DefaultTagParser struct {
	PlainFieldInfos []*StructFieldInfo
}

func (b *DefaultTagParser) WriteStruct(structField reflect.StructField, fieldInfo StructFieldInfo, index int) TagInfoWriter {
	return b
}

func (b *DefaultTagParser) WriteField(field reflect.StructField, fieldInfo StructFieldInfo, index int) {
	if !fieldInfo.HasEmptyTag {
		b.PlainFieldInfos = append(b.PlainFieldInfos, &fieldInfo)
	}
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
	for currentIndex := 0; currentIndex < fieldsCount; currentIndex++ {
		field := typ.Field(currentIndex)

		if field.PkgPath != "" || (SkipEmbeddedFields && field.Anonymous) {
			continue fieldLoop
		}

		fieldTag, _ := field.Tag.Lookup(tagName)

		if fieldTag == "-" {
			continue fieldLoop
		}

		fieldInfo := NewFieldInfo(field, fieldTag, parentIndexChain, currentIndex)

		if fieldInfo.JsonMarshalling {
			tagWriter.WriteField(field, fieldInfo, currentIndex)
			continue fieldLoop
		}

		// if the field is a pointer to a struct, follow the pointer then create field info for each field
		for {
			// if the field implements JsonEncoding we consider it as merely a field
			if fieldInfo.Type.Implements(TextEncoderType) {
				fieldInfo.CustomMarshalling = true
				tagWriter.WriteField(field, fieldInfo, currentIndex)
				continue fieldLoop
			}
			if fieldInfo.Type.Kind() == reflect.Array || fieldInfo.Type.Kind() == reflect.Slice {
				fieldInfo.IsArray = true
				tagWriter.WriteField(field, fieldInfo, currentIndex)
				continue fieldLoop
			}

			if fieldInfo.Type.Kind() != reflect.Ptr {
				break
			}
			fieldInfo.Type = fieldInfo.Type.Elem()
		}

		// if the field is a struct, that doesn't implement JsonEncoding, create a fieldInfo for each of its fields
		if fieldInfo.Type.Kind() == reflect.Struct {
			writer := tagWriter.WriteStruct(field, fieldInfo, currentIndex)
			getFieldInfoRecursive(fieldInfo.Type, tagName, fieldInfo.IndexChain, parentChainSet, writer)
			continue fieldLoop
		}

		if !CheckAllowedFields(fieldInfo.Type) {
			fmt.Printf("skipping an unsupported field: %s, kind: %s\n", fieldInfo.Name, fieldInfo.Type.Kind())
			continue fieldLoop
		}
		tagWriter.WriteField(field, fieldInfo, currentIndex)
	}
}

func CheckAllowedFields(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.String:
		return true
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32:
		return true
	case reflect.Float64:
		return true
	default:
		return false
	}
}

func NewFieldInfo(field reflect.StructField, fieldTag string, parentIndexChain []int, currentIndex int) StructFieldInfo {

	fieldInfo := StructFieldInfo{}

	var cpy = make([]int, len(parentIndexChain))
	copy(cpy, parentIndexChain)
	fieldInfo.IndexChain = append(cpy, currentIndex)
	fieldInfo.Type = field.Type

	fieldInfo.HasEmptyTag = len(fieldTag) == 0

	tagArray := strings.Split(fieldTag, TagSeparator)

	for i, tag := range tagArray {
		if i == 0 {
			if tag != "" {
				fieldInfo.Name = tag
			} else {
				fieldInfo.Name = field.Name
			}
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

	return fieldInfo
}
