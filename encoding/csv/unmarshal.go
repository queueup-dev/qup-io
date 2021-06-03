package csv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/queueup-dev/qup-io/reflection"
	"io"
	"reflect"
)

func (h headerAndBaseRecord) Unmarshal(input []string, out interface{}) error {
	return h.UnmarshalValue(input, reflect.ValueOf(out))
}

func (h headerAndBaseRecord) UnmarshalValue(csvRecord []string, value reflect.Value) error {
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return &json.InvalidUnmarshalError{Type: value.Type()}
	}

	typ := value.Elem().Type()
	if typ != h.baseRecord.origStruct {
		return fmt.Errorf("unknown type to write to. Got: %v expected: %v", value.Type().String(), h.baseRecord.origStruct.String())
	}

	for k, field := range h.mapping {
		content := csvRecord[k]
		if err := reflection.SetFieldFromIndexChain(value.Elem(), field.IndexChain, content, field); err != nil { // Set field of struct
			parseError := ParseError{
				Column:     k + 1,
				ColumnName: field.Name,
				Err:        err,
			}
			return &parseError
		}
	}

	return nil
}

func (h *headerAndBaseRecord) UnmarshalList(payload io.Reader, readerFunction interface{}) {
	csvReader := csv.NewReader(payload)

	if h.Comma != nil {
		csvReader.Comma = *h.Comma
	}

	h.ProvidePayload(csvReader).UnmarshallList(readerFunction)
}

func (h *headerAndBaseRecord) ProvidePayload(payload *csv.Reader) *loadedListReader {
	return &loadedListReader{
		typeContext: h,
		numLine:     0,
		payload:     payload,
	}
}

func (lr *loadedListReader) UnmarshallList(readerFunction interface{}) {
	readerPtrFuncValue := reflect.ValueOf(readerFunction)

	readerFuncTyp, _ := reflection.ValidateReaderFunction(readerPtrFuncValue, lr.typeContext.baseRecord.origStruct)

	reflection.PopulateFunction(readerFuncTyp, readerPtrFuncValue, lr.typeContext.makeReaderFunc(lr.payload))
}

func (h headerAndBaseRecord) makeReaderFunc(payload *csv.Reader) func([]reflect.Value) []reflect.Value {
	return func([]reflect.Value) []reflect.Value {
		ret := make([]reflect.Value, 3)
		csvRecord, err := payload.Read()

		if err != nil {
			ret[0] = reflect.Zero(h.origStruct)

			if err == io.EOF {
				ret[1] = reflect.ValueOf(true)
				ret[2] = reflection.ErrorZeroValue
				return ret
			}

			ret[1] = reflect.ValueOf(false)
			ret[2] = reflect.ValueOf(err)

			return ret
		}

		value := reflect.New(h.origStruct)
		err = h.UnmarshalValue(csvRecord, value)

		ret[0] = value.Elem()
		ret[1] = reflect.ValueOf(false)
		if err != nil {
			ret[2] = reflect.ValueOf(err)
		} else {
			ret[2] = reflection.ErrorZeroValue
		}

		return ret
	}
}
