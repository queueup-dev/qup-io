package reader

import "io"

func NewJsonReader(stream io.Reader) *jsonReader {
	return &jsonReader{input: stream}
}

func NewXmlReader(stream io.Reader) *xmlReader {
	return &xmlReader{input: stream}
}

func NewRawReader(stream io.Reader) *rawReader {
	return &rawReader{input: stream}
}
