package envvar

import (
	"fmt"
	"os"
	"testing"
)

type ExampleStruct struct {
	Foo string `env:"foo,required"`
	Bar string
}

type InvalidTypeExampleStruct struct {
	Foo int `env:"foo,required"`
}

func TestToStructWithMissingRequiredEnvVar(t *testing.T) {
	os.Unsetenv("foo")
	testStruct := ExampleStruct{}

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
