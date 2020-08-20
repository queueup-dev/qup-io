package csv

import (
	"errors"
	"fmt"
	"reflect"
)

type HeaderMapping struct {
	OriginalStruct reflect.Type
	Header         []string
	HeaderMapping  map[int]StructFieldInfo
}

func (rs RegisteredStruct) GenerateHeadersFromStruct() CsvRowMarshaller {
	header, headerMapping := GenerateHeaderFromStruct(rs.fieldInfos)
	return HeaderMapping{
		OriginalStruct: rs.originalStruct,
		Header:         header,
		HeaderMapping:  headerMapping,
	}
}

func (rs RegisteredStruct) MapHeadersToFields(header []string) (CsvRowMarshaller, error) {
	if header == nil {
		return nil, errors.New("no header provided")
	}

	headerMapping, err := MapHeadersToFields(header, rs.fieldInfos)

	if err != nil {
		return nil, err
	}

	return HeaderMapping{
		OriginalStruct: rs.originalStruct,
		Header:         header,
		HeaderMapping:  headerMapping,
	}, nil
}

func (m HeaderMapping) GetHeader() []string {
	return m.Header
}

func GenerateHeaderFromStruct(fieldInfos []StructFieldInfo) ([]string, map[int]StructFieldInfo) {
	header := make([]string, len(fieldInfos))
	headerMapping := make(map[int]StructFieldInfo, len(fieldInfos))

	for i, val := range fieldInfos {
		header[i] = val.name
		headerMapping[i] = val
	}

	return header, headerMapping
}

func MapHeadersToFields(header []string, fieldInfos []StructFieldInfo) (map[int]StructFieldInfo, error) {
	headerMapping := make(map[int]StructFieldInfo, len(fieldInfos))
	headerCount := make(map[string]int)

	if FailIfDoubleHeaderNames {
		for _, name := range header {
			if headerCount[name] > 0 {
				return nil, fmt.Errorf("repeated Header name: %v", name)
			}
			headerCount[name]++
		}
		headerCount = make(map[string]int)
	}

	for _, attribute := range fieldInfos {
		curHeaderCount := headerCount[attribute.name]

		matchedFieldCount := 0
		var selectedPosition *int
		for k, csvColumnHeader := range header {
			if attribute.MatchesKey(csvColumnHeader) {
				if matchedFieldCount >= curHeaderCount {
					selectedPosition = &k
					break
				}
				matchedFieldCount++
			}
		}

		curHeaderCount++

		if selectedPosition != nil {
			headerMapping[*selectedPosition] = attribute
			headerCount[attribute.name] = curHeaderCount
		} else {
			return nil, fmt.Errorf("unable to match the '%v'-th csv-tag '%v' to a column in the csv header", curHeaderCount, attribute.name)
		}
	}

	for i, csvColumnHeader := range header {
		curHeaderCount := headerCount[csvColumnHeader]

		matchedFieldCount := 0
		var selectedField *StructFieldInfo
		for _, field := range fieldInfos {
			if field.MatchesKey(csvColumnHeader) {
				if matchedFieldCount >= curHeaderCount {
					selectedField = &field
					break
				}
				matchedFieldCount++
			}
		}

		if selectedField != nil {
			headerMapping[i] = *selectedField
			curHeaderCount++
			headerCount[csvColumnHeader] = curHeaderCount
		}
	}

	return headerMapping, nil
}
