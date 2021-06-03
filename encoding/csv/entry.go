package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/queueup-dev/qup-io/encoding/csv/tag_parser"
	"github.com/queueup-dev/qup-io/reflection"
	"io"
	"reflect"
)

func RegisterStruct(structVal interface{}) *baseRecord {
	typ := reflect.TypeOf(structVal)

	fieldInfos := tag_parser.NewFieldInfoWriter()
	reflection.GetFieldInfo(typ, TagName, &fieldInfos)

	return &baseRecord{
		origStruct: typ,
		fieldInfo:  fieldInfos,
	}
}

func (b *baseRecord) DeriveHeaderFromStruct() *headerAndBaseRecord {
	header := make([]string, len(b.fieldInfo.PlainFieldInfos))
	mapping := make(map[int]reflection.StructFieldInfo, len(b.fieldInfo.PlainFieldInfos))

	for i, val := range b.fieldInfo.PlainFieldInfos {
		header[i] = val.Name
		mapping[i] = *val
	}

	if len(mapping) == 0 {
		panic(fmt.Errorf("no csv struct tags found for type: %s", b.origStruct.String()))
	}

	return &headerAndBaseRecord{
		baseRecord: b,
		header:     header,
		mapping:    mapping,
	}
}

func (b *baseRecord) RegisterHeader(header []string) (*headerAndBaseRecord, error) {
	mapping, err := b.fieldInfo.MapHeader(header)

	if err != nil {
		return nil, err
	}

	return &headerAndBaseRecord{
		baseRecord: b,
		header:     header,
		mapping:    mapping,
	}, nil
}

func (b *baseRecord) MarshalList(writer io.Writer, writerFunction interface{}) (Flusher, error) {
	return b.DeriveHeaderFromStruct().MarshalList(writer, writerFunction)
}

func (b *baseRecord) UnmarshalList(payload io.Reader, readerFunction interface{}) error {
	csvReader := csv.NewReader(payload)

	if b.Comma != nil {
		csvReader.Comma = *b.Comma
	}

	header, err := csvReader.Read()
	if err != nil {
		fmt.Println("error when reading from csv reader")
		return err
	}

	mapping, err := b.RegisterHeader(header)

	if err != nil {
		return err
	}

	mapping.ProvidePayload(csvReader).UnmarshallList(readerFunction)

	return nil
}
