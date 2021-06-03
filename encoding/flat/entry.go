package flat

import (
	"github.com/queueup-dev/qup-io/reflection"
	"reflect"
)

var (
	TagName = "flat"
)

var fieldInfoCache = make(map[string]reflection.DefaultTagParser)

func ToMap(structVal interface{}, structTag *string) map[string]interface{} {

	if structTag == nil || *structTag == "" {
		structTag = &TagName
	}

	val := reflect.ValueOf(structVal)
	typ := val.Type()

	fieldInfos, exists := fieldInfoCache[typ.String()]

	if !exists {
		reflection.GetFieldInfo(typ, *structTag, &fieldInfos)
		fieldInfoCache[typ.String()] = fieldInfos
	}

	ret := map[string]interface{}{}

	for _, field := range fieldInfos.PlainFieldInfos {
		value := reflection.GetFieldValueFromIndexChain(val, field.IndexChain)

		if field.OmitEmpty && value.IsZero() {
			continue
		}

		ret[field.Name] = value.Interface()
	}

	return ret
}
