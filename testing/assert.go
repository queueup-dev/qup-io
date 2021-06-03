package testing

import (
	"github.com/queueup-dev/qup-io/reflection"
	"reflect"
	"strings"
)

type Assert struct {
	Input interface{}
}

func (a Assert) Numeric(compare interface{}) bool {
	return reflection.IsNumeric(compare)
}

func (a Assert) Eq(compare interface{}) bool {
	return compare == a.Input
}

func (a Assert) Same(compare interface{}) bool {
	return reflection.StringValueOf(a.Input) == reflection.StringValueOf(compare)
}

func (a Assert) Gt(compare interface{}) bool {
	if !reflection.IsNumeric(compare) {
		panic("only numeric values can be compared with Gt")
	}

	return reflection.IntegerOf(a.Input) > reflection.IntegerOf(compare)
}

func (a Assert) Lt(compare interface{}) bool {
	if !reflection.IsNumeric(compare) {
		panic("only numeric values can be compared with Gt")
	}

	return reflection.IntegerOf(a.Input) < reflection.IntegerOf(compare)
}

func (a Assert) Gte(compare interface{}) bool {
	if !reflection.IsNumeric(compare) {
		panic("only numeric values can be compared with Gt")
	}

	return reflection.IntegerOf(a.Input) >= reflection.IntegerOf(compare)
}

func (a Assert) Lte(compare interface{}) bool {
	if !reflection.IsNumeric(compare) {
		panic("only numeric values can be compared with Gt")
	}

	return reflection.IntegerOf(a.Input) <= reflection.IntegerOf(compare)
}

func (a Assert) Contains(compare interface{}) bool {
	return strings.Contains(reflection.StringValueOf(a.Input), reflection.StringValueOf(compare))
}

func (a Assert) Type(compare reflect.Kind) bool {
	return reflect.ValueOf(a.Input).Kind() == compare
}

type AssertInstance struct {
	assert  func(Assert, interface{}) bool
	compare interface{}
	logger  Logger
}

func (b *AssertInstance) Eq(compare interface{}) {
	b.assert = Assert.Eq
	b.compare = compare
}

func (b *AssertInstance) Same(compare interface{}) {
	b.assert = Assert.Same
	b.compare = compare
}

func (b *AssertInstance) Execute(input interface{}) bool {
	return b.assert(Assert{Input: input}, b.compare)
}
