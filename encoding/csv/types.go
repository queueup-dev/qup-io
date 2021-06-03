package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/queueup-dev/qup-io/encoding/csv/tag_parser"
	"github.com/queueup-dev/qup-io/reflection"
	"reflect"
)

var (
	TagName = "csv"
)

type Flusher interface {
	Flush()
	Error() error
}

type baseRecord struct {
	origStruct reflect.Type
	fieldInfo  tag_parser.FieldInfo
	Comma      *rune
	//FailIfDoubleHeaderNames bool
}

type headerAndBaseRecord struct {
	*baseRecord
	header  []string
	mapping map[int]reflection.StructFieldInfo
}

// todo add reuse record pointer option
type loadedListReader struct {
	typeContext *headerAndBaseRecord
	numLine     int
	payload     *csv.Reader
}

type ParseError struct {
	StartLine  int // Line where the record starts
	Line       int // Line where the error occurred
	Column     int // Column (rune index) where the error occurred
	ColumnName string
	Err        error // The actual error
}

func (e *ParseError) Error() string {
	if e.Err == csv.ErrFieldCount {
		return fmt.Sprintf("record on line %d: %v", e.Line, e.Err)
	}
	if e.StartLine != e.Line {
		return fmt.Sprintf("record on line %d; parse error on line %d, column %d %s: %v", e.StartLine, e.Line, e.Column, e.ColumnName, e.Err)
	}
	return fmt.Sprintf("parse error on line %d, column %d %s: %v", e.Line, e.Column, e.ColumnName, e.Err)
}
