package reflection

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	ReflectTest string `test:"Test tag value"`
}

type TestMultiple struct {
	ReflectTest  string `test:"Test tag value"`
	ReflectTest2 string `test:"Test tag value2"`
}

func TestGetFieldNamesWithExistingTag(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	fields, err := GetFieldNamesWithTag("test", reflect.TypeOf(test))

	if err != nil {
		t.Fail()
	}

	if len(fields) != 1 {
		t.Fail()
	}
}

func TestGetFieldNamesWithNonExistingTag(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	fields, err := GetFieldNamesWithTag("foo", reflect.TypeOf(test))

	if err != nil {
		t.Fail()
	}

	if len(fields) != 0 {
		t.Fail()
	}
}

func TestGetFieldNamesWithTagNonStructFails(t *testing.T) {
	_, err := GetFieldNamesWithTag("foo", reflect.TypeOf("bla"))

	if err == nil {
		t.Fail()
	}
}

func TestGetTagValueWithExistingTag(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	val, err := GetTagValue("test", "ReflectTest", reflect.TypeOf(test))

	if err != nil {
		t.Fail()
	}

	if val != "Test tag value" {
		t.Fail()
	}
}

func TestGetTagValueWithNonStruct(t *testing.T) {
	_, err := GetTagValue("foo", "ReflectFoo", reflect.TypeOf("test"))

	if err == nil || err.Error() != "supplied argument is not a structure" {
		t.Fail()
	}
}

func TestGetTagValueWithNonExistingField(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	_, err := GetTagValue("foo", "ReflectFoo", reflect.TypeOf(test))

	if err == nil || err.Error() != "supplied field is not defined in the structure" {
		t.Fail()
	}
}

func TestGetTagValueWithNonExistingTag(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	_, err := GetTagValue("foo", "ReflectTest", reflect.TypeOf(test))

	if err == nil || err.Error() != "supplied tag is not present on the field" {
		t.Fail()
	}
}

func TestGetTagValuesWithNonStruct(t *testing.T) {
	_, err := GetTagValues("foo", reflect.TypeOf("test"))

	if err == nil || err.Error() != "supplied argument is not a structure" {
		t.Fail()
	}
}

func TestGetTagValuesWithNonExistingTag(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	values, err := GetTagValues("foo", reflect.TypeOf(test))

	if err != nil {
		t.Fail()
	}

	if len(values) != 0 {
		t.Fail()
	}
}

func TestGetTagValuesWithMultipleTags(t *testing.T) {
	test := &TestMultiple{
		ReflectTest:  "Hello World",
		ReflectTest2: "Foo bar",
	}

	values, err := GetTagValues("test", reflect.TypeOf(test))

	if err != nil {
		t.Fail()
	}

	if len(values) != 2 {
		t.Fail()
	}
}

func TestGetTagValuesWithExistingTag(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	values, err := GetTagValues("test", reflect.TypeOf(test))

	if err != nil {
		t.Fail()
	}

	if len(values) != 1 {
		t.Fail()
	}
}
