package reflection

import (
	"fmt"
	"log"
	"math"
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
	Map     map[string]int
	Func    func()
	Ptr     *string
	Array   []int
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

	err = PopulateFromString(field, fmt.Sprintf("%f", math.MaxFloat64), false)

	if err == nil {
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

	err = PopulateFromString(field, fmt.Sprintf("1.797693134862315708145274237317043567981e+309"), false)

	if err == nil {
		t.Fail()
	}
}

func TestPopulateFromStringWithBool(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Bool")

	err := PopulateFromString(field, "nonbool-lolkek", false)

	if err == nil {
		t.Fail()
	}

	err = PopulateFromString(field, "true", false)

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

func TestPopulateFromStringWithUnsupportedTypeFails(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Func")

	err := PopulateFromString(field, "Hello World!", false)

	if err == nil {
		t.Fail()
	}

	if err.Error() != "unsupported field type, got: func" {
		t.Fail()
	}
}

// Order matters with this test.
func TestPopulateFromStringWithPointer(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Ptr")

	err := PopulateFromString(field, "", true)

	if err != nil {
		t.Fail()
	}

	if initializedStruct.Ptr != nil {
		t.Fail()
	}

	err = PopulateFromString(field, "", false)

	if err != nil {
		t.Fail()
	}

	if initializedStruct.Ptr == nil {
		t.Fail()
	}

	err = PopulateFromString(field, "Hello World!", false)

	if err != nil {
		t.Fail()
	}

	if *initializedStruct.Ptr != "Hello World!" {
		t.Fail()
	}
}

func TestPopulateFromStringWithArray(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Array")

	err := PopulateFromString(field, "[1,2,3]", false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if len(initializedStruct.Array) != 3 {
		t.Fail()
	}

	if initializedStruct.Array[0] != 1 && initializedStruct.Array[1] != 2 && initializedStruct.Array[2] != 3 {
		t.Fail()
	}
}

func TestPopulateFromStringWithMap(t *testing.T) {
	initializedStruct := &ExampleStruct{}
	reflectedStruct := reflect.ValueOf(initializedStruct)
	field := reflectedStruct.Elem().FieldByName("Map")

	err := PopulateFromString(field, `{"a":1,"b":2,"c":3}`, false)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if len(initializedStruct.Map) != 3 {
		t.Fail()
	}

	if initializedStruct.Map["a"] != 1 && initializedStruct.Map["b"] != 2 && initializedStruct.Map["c"] != 3 {
		t.Fail()
	}
}
