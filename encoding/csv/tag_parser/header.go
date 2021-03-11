package tag_parser

import (
	"errors"
	"fmt"
	"github.com/queueup-dev/qup-io/reflection"
	"github.com/queueup-dev/qup-io/slices"
)

func (b FieldInfo) MapHeader(header []string) (map[int]reflection.StructFieldInfo, error) {
	// todo check for double headers
	// todo implement header count

	for requiredColumn, occurrences := range b.requiredFields {
		for _, occ := range occurrences {
			if _, found := slices.CountedHasString(requiredColumn, header, occ); !found {
				return nil, fmt.Errorf("unable to match the required '%v'-th csv-tag '%s' to a column in the csv header", occ, requiredColumn)
			}
		}
	}

	mappedHeaders := make(map[int]reflection.StructFieldInfo)

	headerCount := make(map[string]int)
	for i, column := range header {
		infoCollection, ok := b.aggregatedFieldInfos[column]
		if !ok {
			continue
		}

		selectedInfo := infoCollection[headerCount[column]]
		headerCount[column]++
		mappedHeaders[i] = *selectedInfo
	}

	if len(mappedHeaders) == 0 {
		return nil, errors.New("no csv struct tags found")
	}

	return mappedHeaders, nil
}
