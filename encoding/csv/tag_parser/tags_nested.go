package tag_parser

import (
	"github.com/queueup-dev/qup-io/reflection"
	"reflect"
)

var (
	TagName      = "csv"
	TagSeparator = ","
)

// Sum type, only 1 of two can be non-empty
type FieldOrStruct struct {
	StructFieldInfo *reflection.StructFieldInfo
	StructInfo      *NestedStructInfo
	Index           int
}

type NestedStructInfo struct {
	StructFields []FieldOrStruct
	Name         string
	Index        int
}

func (s *NestedStructInfo) WriteStruct(structField reflect.StructField, fieldInfo reflection.StructFieldInfo, index int) reflection.TagInfoWriter {
	structInfo := NestedStructInfo{
		Name:         structField.Name,
		StructFields: make([]FieldOrStruct, 0),
	}

	s.StructFields = append(s.StructFields, FieldOrStruct{StructInfo: &structInfo, Index: index})
	return &structInfo
}

func (s *NestedStructInfo) WriteField(structField reflect.StructField, fieldInfo reflection.StructFieldInfo, index int) {
	s.StructFields = append(s.StructFields, FieldOrStruct{StructFieldInfo: &fieldInfo, Index: index})
}
