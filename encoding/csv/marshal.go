package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/queueup-dev/qup-io/reflection"
	"io"
	"reflect"
)

func (h headerAndBaseRecord) Marshal(input interface{}) ([]string, error) {
	return h.MarshalValue(reflect.ValueOf(input))
}

func (h headerAndBaseRecord) MarshalValue(value reflect.Value) ([]string, error) {
	output := make([]string, len(h.header))

	if !value.Type().AssignableTo(h.origStruct) {
		// this is needed when the writer functions input type is not the same as the base struct
		// todo make a bug report on this.
		value = reflect.ValueOf(value.Interface())
		if !value.Type().AssignableTo(h.origStruct) {
			return nil, fmt.Errorf("type '%s' provided to the writer function is incompatible with the base type '%s'", value.Type().String(), h.origStruct.String())
		}
	}

	for k, field := range h.mapping {
		stringValue, err := reflection.GetFieldValueFromIndexChainAsString(value, field.IndexChain)
		if err != nil {
			return nil, err
		}
		output[k] = stringValue
	}

	return output, nil
}

func (h headerAndBaseRecord) MarshalList(writer io.Writer, writerFunction interface{}) (Flusher, error) {
	writerVal := reflect.ValueOf(writerFunction)
	writerType := reflection.ValidateWriterFunction(writerVal, h.baseRecord.origStruct)

	csvWriter := csv.NewWriter(writer)

	if h.Comma != nil {
		csvWriter.Comma = *h.Comma
	}

	err := csvWriter.Write(h.header)
	if err != nil {
		return nil, err
	}

	reflection.PopulateFunction(writerType, writerVal, h.makeWriterFunc(csvWriter))

	return csvWriter, nil
}

func (h headerAndBaseRecord) makeWriterFunc(payload *csv.Writer) func([]reflect.Value) []reflect.Value {
	return func(in []reflect.Value) []reflect.Value {
		record := in[0]
		csvRecord, err := h.MarshalValue(record)

		if err != nil {
			return []reflect.Value{reflect.ValueOf(err)}
		}

		err = payload.Write(csvRecord)

		if err != nil {
			return []reflect.Value{reflect.ValueOf(err)}
		}

		return []reflect.Value{reflection.ErrorZeroValue}
	}
}
