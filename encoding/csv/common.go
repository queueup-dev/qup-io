package csv

import (
	"encoding/csv"
	types "github.com/queueup-dev/qup-types"
)

var (
	TagName                 = "csv"
	TagSeparator            = ","
	FailIfDoubleHeaderNames = false
	SkipEmbeddedFields      = false
)

type CsvMarshaller interface {
	Unmarshal(interface{}) error
	Marshal(interface{}) error
	GetHeader() []string
}

type CsvRowMarshaller interface {
	Unmarshal([]string, interface{}) error
	Marshal(interface{}) ([]string, error)
	GetHeader() []string
	NewReader(*csv.Reader) types.PayloadReader
}
