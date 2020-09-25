package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	types "github.com/queueup-dev/qup-types"
	"io"
	"reflect"
)

type csvReader struct {
	numLine int
	input   *csv.Reader
	HeaderMapping
}

func (rs RegisteredStruct) NewReader(reader *csv.Reader, mapHeaderFromReader bool) (*csvReader, error) {
	m := HeaderMapping{OriginalStruct: rs.originalStruct}

	if mapHeaderFromReader {
		header, err := reader.Read()
		if err != nil {
			fmt.Print("no lines in csv file when trying to read header")
			return nil, err
		}
		m.HeaderMapping, err = MapHeadersToFields(header, rs.fieldInfos)
		if err != nil {
			return nil, err
		}
	} else {
		m.Header, m.HeaderMapping = GenerateHeaderFromStruct(rs.fieldInfos)
	}

	return &csvReader{
		input:         reader,
		HeaderMapping: m,
	}, nil
}

func (m HeaderMapping) NewReader(reader *csv.Reader) types.PayloadReader {
	return &csvReader{
		input:         reader,
		HeaderMapping: m,
	}
}

func (rs RegisteredStruct) NewReaderWithHeader(reader *csv.Reader, i []string) (types.PayloadReader, error) {
	panic("implement me")
}

func (c csvReader) Unmarshal(out interface{}) error {
	val := reflect.ValueOf(out)
	//value := val.Elem()
	originalType := val.Type()

	if val.Kind() == reflect.Chan {
		elementType := originalType.Elem()
		fmt.Println(elementType)

		if c.HeaderMapping.OriginalStruct != elementType {
			return errors.New("struct not registered")
		}

		for {
			row, err := c.input.Read()

			if err != nil {
				if err == io.EOF {
					fmt.Printf("exiting after reading %v lines\n", c.numLine)

					break
				}
				panic(err)
			}

			c.numLine++

			x := reflect.New(elementType)
			err = c.HeaderMapping.UnmarshalValue(row, x)
			if err != nil {
				panic(err)
			}

			val.Send(x.Elem())
		}

		val.Close()
	}

	return nil
}

func (c csvReader) ContentType() string {
	return "csv"
}

func (c csvReader) Valid() bool {
	panic("implement me")
}

func (c csvReader) Bytes() ([]byte, error) {
	panic("implement me")
}

func (c csvReader) Reader() (io.Reader, error) {
	panic("implement me")
}

func (c csvReader) ToString() (*string, error) {
	panic("implement me")
}
