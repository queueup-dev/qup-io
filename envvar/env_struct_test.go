package envvar

import (
	"fmt"
	"os"
	"testing"
)

type NotStruct bool

type ExampleStruct struct {
	Foo   string `env:"foo,required"`
	Bar   string
	Unset string `env:"not-here"`
}

type InvalidTypeExampleStruct struct {
	Foo int `env:"foo,required"`
}

type UnexportedVarExampleStruct struct {
	bar string `env:"bar"`
	foo string `env:"foo,required"`
}

func TestToStructWithNoStructFails(t *testing.T) {
	test := NotStruct(false)

	err := ToStruct(&test)

	if err == nil {
		t.Fail()
	}
}

func TestToStructWithMissingRequiredEnvVar(t *testing.T) {
	os.Unsetenv("foo")
	testStruct := ExampleStruct{}

	err := ToStruct(&testStruct)

	if err == nil {
		t.Fail()
	}
}

func TestToStructWithUnexportedVarFails(t *testing.T) {
	os.Setenv("bar", "hi")
	os.Setenv("foo", "hi")
	testStruct := UnexportedVarExampleStruct{}

	err := ToStruct(&testStruct)

	if err == nil {
		t.Fail()
	}
}

func TestToStructWithInvalidType(t *testing.T) {
	testStruct := InvalidTypeExampleStruct{}
	os.Setenv("foo", "Hello World!")
	err := ToStruct(&testStruct)

	if err == nil {
		t.Fail()
	}
}

func TestToStruct(t *testing.T) {
	testStruct := ExampleStruct{}
	os.Setenv("foo", "Hello World!")
	err := ToStruct(&testStruct)

	if err != nil {
		fmt.Print(err)
	}

	if testStruct.Foo != "Hello World!" {
		t.Fail()
	}
}
