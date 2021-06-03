package tag_parser

import (
	"github.com/queueup-dev/qup-io/reflection"
	"reflect"
)

type FieldInfo struct {
	aggregatedFieldInfos map[string][]*reflection.StructFieldInfo
	PlainFieldInfos      []*reflection.StructFieldInfo
	requiredFields       map[string][]int
	headerCounts         map[string]int
}

func NewFieldInfoWriter() FieldInfo {
	return FieldInfo{
		aggregatedFieldInfos: make(map[string][]*reflection.StructFieldInfo, 0),
		requiredFields:       make(map[string][]int),
		headerCounts:         make(map[string]int),
	}
}

func (b *FieldInfo) WriteStruct(structField reflect.StructField, fieldInfo reflection.StructFieldInfo, index int) reflection.TagInfoWriter {
	return b
}

func (b *FieldInfo) WriteField(field reflect.StructField, fieldInfo reflection.StructFieldInfo, index int) {
	b.headerCounts[fieldInfo.Name]++

	b.aggregatedFieldInfos[fieldInfo.Name] = append(b.aggregatedFieldInfos[fieldInfo.Name], &fieldInfo)
	b.PlainFieldInfos = append(b.PlainFieldInfos, &fieldInfo)
}
