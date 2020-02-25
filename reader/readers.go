package reader

import "io"

func NewJsonReader(stream io.Reader, source string) *JsonReader {
	return &JsonReader{output: stream, source: source}
}

func NewXmlReader(stream io.Reader, source string) *XmlReader {
	return &XmlReader{output: stream, source: source}
}

func NewRawReader(stream io.Reader, source string) *RawReader {
	return &RawReader{output: stream, source: source}
}
