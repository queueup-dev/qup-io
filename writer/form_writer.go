package writer

import (
	"encoding/json"
	"io"
	"net/url"
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

	var targetMap map[string]string
	err = json.Unmarshal(encoded, &targetMap)

	if err != nil {
		return nil, err
	}

	values := url.Values{}
	for key, val := range targetMap {
		values.Add(key, val)
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
