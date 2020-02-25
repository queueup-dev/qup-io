package reader

import "io"

func NewJsonReader(stream io.Reader) *JsonReader {
	return &JsonReader{output: stream}
}

func NewXmlReader(stream io.Reader) *XmlReader {
	return &XmlReader{output: stream}
}

func NewRawReader(stream io.Reader) *RawReader {
	return &RawReader{output: stream}
}
