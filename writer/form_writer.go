package writer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"
)

type FormEncodeWriter struct {
	input interface{}
}

func (f FormEncodeWriter) Bytes() ([]byte, error) {
	marshalOutput, err := f.Marshal()

	if err != nil {
		return nil, err
	}

	return marshalOutput.([]byte), nil
}

func (f FormEncodeWriter) Marshal() (interface{}, error) {
	encoded, err := json.Marshal(f.input)

	if err != nil {
		return nil, err
	}

	var targetMap map[string]interface{}
	err = json.Unmarshal(encoded, &targetMap)

	if err != nil {
		return nil, err
	}

	values := url.Values{}
	for key, val := range targetMap {
		fmt.Print(reflect.TypeOf(val))
		arrayValue, isArray := val.(map[string]interface{})

		if isArray {
			for arrayKey, arrayValue := range arrayValue {
				values.Add(key+"."+arrayKey, arrayValue.(string))
			}

			continue
		}

		stringValue, isString := val.(string)

		if isString {
			values.Add(key, stringValue)

			continue
		}

		return nil, errors.New("invalid data type supplied, can only use map[string]string or string values")
	}

	return values.Encode(), nil
}

func (f FormEncodeWriter) Valid() bool {
	_, err := f.Marshal()

	if err != nil {
		return false
	}

	return true
}

func (f FormEncodeWriter) Reader() (io.Reader, error) {
	data, err := f.Marshal()
	fmt.Print(err)
	content := data.(string)

	if err != nil {
		return nil, err
	}

	return strings.NewReader(content), nil
}

func (f FormEncodeWriter) ToString() (*string, error) {
	content, err := f.Marshal()
	bytes := content.([]byte)

	if err != nil {
		return nil, err
	}

	output := string(bytes)
	return &output, nil
}

func (f FormEncodeWriter) ContentType() string {
	return "application/form-url-encoded"
}
