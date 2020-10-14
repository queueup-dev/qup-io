package reflection

import (
	"testing"
)

type TestStruct struct {
	ReflectTest string `test:"Test tag value"`
}

func TestGetFieldNamesWithExistingTag(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	fields, err := GetFieldNamesWithTag("test", test)

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

	fields, err := GetFieldNamesWithTag("foo", test)

	if err != nil {
		t.Fail()
	}

	if len(fields) != 0 {
		t.Fail()
	}
}

func TestGetFieldNamesWithTagNonStructFails(t *testing.T) {
	_, err := GetFieldNamesWithTag("foo", "bla")

	if err == nil {
		t.Fail()
	}
}

func TestGetTagValueWithExistingTag(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	val, err := GetTagValue("test", "ReflectTest", test)

	if err != nil {
		t.Fail()
	}

	if val != "Test tag value" {
		t.Fail()
	}
}

func TestGetTagValueWithNonStruct(t *testing.T) {
	_, err := GetTagValue("foo", "ReflectFoo", "test")

	if err == nil || err.Error() != "supplied argument is not a structure" {
		t.Fail()
	}
}

func TestGetTagValueWithNonExistingField(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	_, err := GetTagValue("foo", "ReflectFoo", test)

	if err == nil || err.Error() != "supplied field is not defined in the structure" {
		t.Fail()
	}
}

func TestGetTagValueWithNonExistingTag(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	_, err := GetTagValue("foo", "ReflectTest", test)

	if err == nil || err.Error() != "supplied tag is not present on the field" {
		t.Fail()
	}
}

func TestGetTagValuesWithNonStruct(t *testing.T) {
	_, err := GetTagValues("foo", "test")

	if err == nil || err.Error() != "supplied argument is not a structure" {
		t.Fail()
	}
}

func TestGetTagValuesWithNonExistingTag(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	values, err := GetTagValues("foo", test)

	if err != nil {
		t.Fail()
	}

	if len(values) != 0 {
		t.Fail()
	}
}

func TestGetTagValuesWithExistingTag(t *testing.T) {
	test := &TestStruct{
		ReflectTest: "Hello World",
	}

	values, err := GetTagValues("test", test)

	if err != nil {
		t.Fail()
	}

	if len(values) != 1 {
		t.Fail()
	}
}
