package lambda

import (
	"encoding"
	"reflect"
)

var (
	getMessageType = reflect.TypeOf((*SingleMessage)(nil)).Elem()
	//getCollectionType  = reflect.TypeOf((*MessageArray)(nil)).Elem()
	unmarshallTextType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

type SingleMessage interface {
	GetMessage() []byte
}

//type MessageArray interface {
//	GetCollection() interface{}
//}

type AwsEvent interface {
	CastToArray() bool
}
