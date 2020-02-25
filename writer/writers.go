package writer

func NewJsonWriter(input interface{}) *JsonWriter {
	return &JsonWriter{input: input}
}

func NewXmlWriter(input interface{}) *XmlWriter {
	return &XmlWriter{input: input}
}

func NewRawWriter(input interface{}) *RawWriter {
	return &RawWriter{input: input}
}
