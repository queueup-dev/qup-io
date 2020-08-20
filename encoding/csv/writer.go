package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// The exported fields can be changed to customize the details before the first call
type CsvWriter struct {
	HeaderMapping         HeaderMapping
	chanInput             *reflect.Value
	currentLine           int
	currentLineAfterFlush int

	// Only access the exported fields of csv.Writer, not the methods.
	CsvWriter   *csv.Writer
	BufferSize  int
	WriteHeader bool
}

func (rs RegisteredStruct) NewWriterFromChannel(input interface{}, header []string) (*CsvWriter, error) {

	var err error
	mapping := HeaderMapping{
		OriginalStruct: rs.originalStruct,
	}

	if header == nil {
		mapping.Header, mapping.HeaderMapping = GenerateHeaderFromStruct(rs.fieldInfos)
	} else {
		mapping.Header = header
		mapping.HeaderMapping, err = MapHeadersToFields(header, rs.fieldInfos)
		if err != nil {
			return nil, err
		}
	}

	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Chan {
		return nil, fmt.Errorf("expected a channel, instead got: %v", val.Kind())
	}
	originalType := val.Type()
	elementType := originalType.Elem()
	//fmt.Println(elementType)

	if rs.originalStruct != elementType {
		return nil, errors.New("struct not registered")
	}

	return &CsvWriter{
		chanInput:     &val,
		HeaderMapping: mapping,
	}, nil
}

func (writer CsvWriter) Marshal(outWriter io.Writer) {

	writer.CsvWriter = csv.NewWriter(outWriter)

	if writer.WriteHeader {
		writer.CsvWriter.Write(writer.HeaderMapping.Header)
	}

	if writer.chanInput != nil {
		for {
			z, ok := writer.chanInput.Recv()
			if !ok {
				writer.CsvWriter.Flush()
				if err := writer.CsvWriter.Error(); err != nil {
					panic(err)
				}

				fmt.Println("the write channel is now closed and the buffer flushed")
				break
			}

			out, err := writer.HeaderMapping.MarshalValue(z)
			if err != nil {
				panic(err)
			}

			err = writer.CsvWriter.Write(out)
			if err != nil {
				panic(err)
			}

			writer.currentLine++
			writer.currentLineAfterFlush++

			if writer.BufferSize >= 0 && writer.currentLineAfterFlush >= writer.BufferSize {
				if err := writer.CsvWriter.Error(); err != nil {
					panic(err)
				}
				writer.CsvWriter.Flush()
				writer.currentLineAfterFlush = 0
				if err := writer.CsvWriter.Error(); err != nil {
					panic(err)
				}
			}
		}
	}
}

func (writer CsvWriter) Valid() bool {
	panic("implement me")
}

func (writer CsvWriter) Bytes() ([]byte, error) {
	panic("implement me")
}

func (writer CsvWriter) Reader() (io.Reader, error) {
	panic("implement me")
}

func (writer CsvWriter) ToString() (*string, error) {
	panic("implement me")
}

func (writer CsvWriter) ContentType() string {
	return "csv"
}

func (writer CsvWriter) GetHeader() []string {
	return writer.HeaderMapping.Header
}

func (writer CsvWriter) Flush() {
	writer.CsvWriter.Flush()
}

func (writer CsvWriter) GetError() error {
	return writer.CsvWriter.Error()
}
