package reflection

import (
	"log"
	"reflect"
	"testing"
)

type ExampleStruct struct {
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64
	Bool    bool
	String  string
}

func TestPopulateFromStringWithInt8(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Int8")

	err := PopulateFromString(field, "128", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-129", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-128", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Int8 != -128 {
		t.Fail()
	}

	err = PopulateFromString(field, "28", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Int8 != 28 {
		t.Fail()
	}
}

func TestPopulateFromStringWithInt16(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Int16")

	err := PopulateFromString(field, "32768", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-32769", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-32768", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Int16 != -32768 {
		t.Fail()
	}

	err = PopulateFromString(field, "512", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Int16 != 512 {
		t.Fail()
	}
}

func TestPopulateFromStringWithInt32(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Int32")

	err := PopulateFromString(field, "2147483648", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-2147483649", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-2147483648", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Int32 != -2147483648 {
		t.Fail()
	}

	err = PopulateFromString(field, "214740000", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Int32 != 214740000 {
		t.Fail()
	}
}

func TestPopulateFromStringWithInt64(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Int64")

	err := PopulateFromString(field, "9223372036854775808", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-9223372036854775809", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-9223372036854775808", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Int64 != -9223372036854775808 {
		t.Fail()
	}

	err = PopulateFromString(field, "9223372036854775800", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Int64 != 9223372036854775800 {
		t.Fail()
	}

}

func TestPopulateFromStringWithUint8(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Uint8")

	err := PopulateFromString(field, "256", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-1", false)

	if err == nil {
		log.Print("failed on the negative uint test")
		t.Fail()
	}

	err = PopulateFromString(field, "28", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Uint8 != 28 {
		t.Fail()
	}
}

func TestPopulateFromStringWithUint16(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Uint16")

	err := PopulateFromString(field, "65536", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-1", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "512", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Uint16 != 512 {
		t.Fail()
	}
}

func TestPopulateFromStringWithUint32(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Uint32")

	err := PopulateFromString(field, "4294967296", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-1", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "214740000", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Uint32 != 214740000 {
		t.Fail()
	}
}

func TestPopulateFromStringWithUint64(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Uint64")

	err := PopulateFromString(field, "18446744073709551616", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "-1", false)

	if err == nil {
		log.Print(err)
		t.Fail()
	}

	err = PopulateFromString(field, "9223372036854775800", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Uint64 != 9223372036854775800 {
		t.Fail()
	}
}

func TestPopulateFromStringWithFloat32(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Float32")

	err := PopulateFromString(field, "214740000.23781278321", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Float32 != 214740000.23781278321 {
		t.Fail()
	}
}

func TestPopulateFromStringWithFloat64(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Float64")

	err := PopulateFromString(field, "9223372036854775800.12382136172", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Float64 != 9223372036854775800.12382136172 {
		t.Fail()
	}
}

func TestPopulateFromStringWithBool(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Bool")

	err := PopulateFromString(field, "true", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if !initializedStruct.Bool {
		t.Fail()
	}

	err = PopulateFromString(field, "false", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.Bool {
		t.Fail()
	}
}

func TestPopulateFromStringWithString(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("String")

	err := PopulateFromString(field, "Hello World!", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if initializedStruct.String != "Hello World!" {
		t.Fail()
	}
}
