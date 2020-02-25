package writer

func NewJsonWriter(input interface{}, target string) *JsonWriter {
	return &JsonWriter{input: input, target: target}
}

func NewXmlWriter(input interface{}, target string) *XmlWriter {
	return &XmlWriter{input: input, target: target}
}

func NewRawWriter(input interface{}, target string) *RawWriter {
	return &RawWriter{input: input, target: target}
}
